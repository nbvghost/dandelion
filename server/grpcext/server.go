package grpcext

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/server/route"
	"github.com/nbvghost/tool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/nbvghost/dandelion/server/serviceobject"

	"google.golang.org/grpc"
)

type iCustomizeService interface {
	Call(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)
}
type customizeService struct {
	server      config.MicroServerConfig
	serviceDesc grpc.ServiceDesc
	routes      map[string]*route.Info
	redis       constrain.IRedis
	callbacks   []constrain.IMappingCallback
}

func (m *customizeService) Call(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	serverTransportStream := grpc.ServerTransportStreamFromContext(ctx)

	var ok bool
	var md metadata.MD
	md, ok = metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Newf(codes.Unknown, "无效的上下文:%v", ctx).Err()
	}

	var uid string
	clientIDs := md.Get("clientid")
	if len(clientIDs) > 0 {
		uid = clientIDs[0]
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger = logger.Named("GrpcContext").With(zap.String("TraceID", tool.UUID()))
	defer logger.Sync()

	currentContext := contexext.New(ctx, m.server.MicroServer.Name, uid, serverTransportStream.Method(), m.redis, "", logger, "")

	var r *route.Info

	if r, ok = m.routes[serverTransportStream.Method()]; !ok {
		return nil, status.New(codes.NotFound, "没有找到路由").Err()
	}

	handle := reflect.New(r.GetHandlerType()).Interface().(constrain.IGrpcHandler)
	if err := dec(handle); err != nil {
		return nil, err
	}

	for index := range m.callbacks {
		item := m.callbacks[index]
		if err := item.Before(currentContext, handle); err != nil {
			return nil, err
		}
	}

	if interceptor == nil {
		out, st := handle.Handle(currentContext)
		if st != nil {
			return out, st.Err()
		}
		return out, nil
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: serverTransportStream.Method(),
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		out, st := handle.Handle(currentContext)
		if st != nil {
			return out, st.Err()
		}
		return out, nil
	}
	return interceptor(currentContext, handle, info, handler)

}

type Option func(*serviceobject.ServerDesc, *grpc.Server) error
type service struct {
	//serviceobject.UnimplementedServerServer
	server     config.MicroServerConfig
	routes     map[string]*route.Info
	redis      constrain.IRedis
	grpcServer *grpc.Server
	option     Option
	callbacks  []constrain.IMappingCallback
}

func loggersss(format string, a ...interface{}) {
	fmt.Printf("LOG:\t"+format+"\n", a...)
}
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

func (m *service) Register(serviceDesc grpc.ServiceDesc, handlers []constrain.IGrpcHandler, withoutAuth ...bool) {
	var _withoutAuth bool
	if len(withoutAuth) > 0 {
		_withoutAuth = withoutAuth[0]
	}

	messageType := reflect.TypeOf(new(proto.Message)).Elem()
	t := reflect.TypeOf(serviceDesc.HandlerType).Elem()

	num := t.NumMethod()
	for i := 0; i < num; i++ {
		method := t.Method(i)
		if method.IsExported() {
			var isFound bool
			for index := range handlers {
				hT := reflect.TypeOf(handlers[index]).Elem()
				hFieldNum := hT.NumField()
				for ii := 0; ii < hFieldNum; ii++ {
					if reflect.New(hT.Field(ii).Type).Type().Implements(messageType) {
						aT := method.Type.In(1).Elem()
						bT := hT.Field(ii).Type
						if fmt.Sprintf("%s.%s", aT.PkgPath(), aT.Name()) == fmt.Sprintf("%s.%s", bT.PkgPath(), bT.Name()) {
							fullServiceName := fmt.Sprintf("/%s/%s", serviceDesc.ServiceName, method.Name)
							if _, ok := m.routes[fullServiceName]; ok {
								panic(errors.New(fmt.Sprintf("存在相同的路由:%s", fullServiceName)))
							}
							m.routes[fullServiceName] = &route.Info{
								HandlerType: hT,
								WithoutAuth: _withoutAuth,
							}
							isFound = true
						}
					}
				}
			}
			if !isFound {
				panic(errors.New(fmt.Sprintf("没有处理路由:%s", method.Name)))
			}

		}
	}

	if m.grpcServer != nil {
		customize := &customizeService{
			server:      m.server,
			serviceDesc: serviceDesc,
			routes:      m.routes,
			redis:       m.redis,
			callbacks:   m.callbacks,
		}
		serviceDesc.HandlerType = (*iCustomizeService)(nil)
		for i := 0; i < len(serviceDesc.Methods); i++ {
			serviceDesc.Methods[i].Handler = customize.Call
		}
		m.grpcServer.RegisterService(&serviceDesc, customize)
	}
}

func (m *service) getRouteInfo(serverInfo *grpc.UnaryServerInfo) (constrain.IRouteInfo, error) {
	var routeInfo *route.Info
	var ok bool
	var err error

	if routeInfo, ok = m.routes[serverInfo.FullMethod]; !ok {
		err = action.NewCodeWithError(action.NotFoundRoute, errors.New("没有找到路由"))
	}

	return routeInfo, err
}

func (m *service) AddCallback(callbacks ...constrain.IMappingCallback) {
	m.callbacks = append(m.callbacks, callbacks...)
}
func (m *service) Listen() {
	var ip = m.server.IP
	var port = m.server.Port
	if ip == "" {
		ip = util.NetworkIP()
		if ip == "" {
			log.Fatalln(errors.New("无法获取本机ip"))
		}
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		log.Fatalln(lis.Close())
	}()

	if port == 0 {
		_, _port, err := net.SplitHostPort(lis.Addr().String())
		if err != nil {
			log.Fatalln(err)
		}
		port, _ = strconv.Atoi(_port)
	}

	desc := &serviceobject.ServerDesc{
		MicroServer: m.server.MicroServer,
		Port:        port,
		IP:          ip,
	}

	//s := grpc.NewServer(grpc.UnaryInterceptor(m.unaryInterceptor))
	//s := grpc.NewServer()
	defer m.grpcServer.Stop()

	if err = m.option(desc, m.grpcServer); err != nil {
		log.Fatalln(err)
	}
	reflection.Register(m.grpcServer)
	if err = m.grpcServer.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}

func New(server config.MicroServerConfig, redis constrain.IRedis, option Option, serverOptions ...grpc.ServerOption) IGrpc {
	return &service{
		server:     server,
		routes:     make(map[string]*route.Info),
		option:     option,
		redis:      redis,
		grpcServer: grpc.NewServer(serverOptions...),
	}
}
