package grpcext

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/server/route"
	"github.com/nbvghost/tool"

	"google.golang.org/grpc"
)

type iCustomizeService interface {
	Call(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)
}
type customizeService struct {
	server      config.MicroServerConfig
	serviceDesc grpc.ServiceDesc
	routes      map[string]*route.RouteInfo
	etcd        constrain.IEtcd
	redis       constrain.IRedis
	callback    constrain.IMappingCallback
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
	logger = logger.With(zap.String("Path", tool.UUID()))
	defer logger.Sync()

	currentContext := contexext.New(ctx, m.server.MicroServer.Name, uid, serverTransportStream.Method(), m.callback, m.etcd, m.redis, "", logger, "")

	var r *route.RouteInfo

	if r, ok = m.routes[serverTransportStream.Method()]; !ok {
		return nil, status.New(codes.NotFound, "没有找到路由").Err()
	}

	handle := reflect.New(r.GetHandlerType()).Interface().(constrain.IGrpcHandler)
	if err := dec(handle); err != nil {
		return nil, err
	}

	if m.callback != nil {
		m.callback.Mapping(currentContext, handle)
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

type Option func(*config.MicroServerConfig, *grpc.Server) error
type service struct {
	//serviceobject.UnimplementedServerServer
	server     config.MicroServerConfig
	routes     map[string]*route.RouteInfo
	redis      constrain.IRedis
	etcdServer constrain.IEtcd
	grpcServer *grpc.Server
	option     Option
	callback   constrain.IMappingCallback
}

func (m *service) Server() *grpc.Server {
	return m.grpcServer
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
							m.routes[fullServiceName] = &route.RouteInfo{
								HandlerType: hT,
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
			etcd:        m.etcdServer,
			redis:       m.redis,
			callback:    m.callback,
		}
		serviceDesc.HandlerType = (*iCustomizeService)(nil)
		for i := 0; i < len(serviceDesc.Methods); i++ {
			serviceDesc.Methods[i].Handler = customize.Call
		}
		m.grpcServer.RegisterService(&serviceDesc, customize)
	}
}

func (m *service) getRouteInfo(serverInfo *grpc.UnaryServerInfo) (constrain.IRouteInfo, error) {
	var routeInfo *route.RouteInfo
	var ok bool
	var err error

	if routeInfo, ok = m.routes[serverInfo.FullMethod]; !ok {
		err = result.NewCodeWithMessage(result.NotFound, "没有找到路由")
	}

	return routeInfo, err
}

func (m *service) AddMapping(callback constrain.IMappingCallback) {
	m.callback = callback
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

	desc := &config.MicroServerConfig{
		MicroServer: m.server.MicroServer,
		Port:        port,
		IP:          ip,
	}

	//s := grpc.NewServer(grpc.UnaryInterceptor(m.unaryInterceptor))
	//s := grpc.NewServer()
	defer m.grpcServer.Stop()

	desc, err = m.etcdServer.Register(config.NewMicroServerConfig(desc.MicroServer, desc.Port, desc.IP))
	if err != nil {
		panic(err)
	}

	if err = m.option(desc, m.grpcServer); err != nil {
		log.Fatalln(err)
	}
	//reflection.Register(m.grpcServer)

	if err = m.grpcServer.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}

func New(server config.MicroServerConfig, redis constrain.IRedis, etcdServer constrain.IEtcd, option Option, serverOptions ...grpc.ServerOption) IGrpc {
	return &service{
		server:     server,
		routes:     make(map[string]*route.RouteInfo),
		option:     option,
		redis:      redis,
		etcdServer: etcdServer,
		grpcServer: grpc.NewServer(serverOptions...),
	}
}
