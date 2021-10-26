package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/service/iservice"

	"github.com/nbvghost/dandelion/service/serviceobject"
	"github.com/nbvghost/dandelion/utils"
	"github.com/nbvghost/gweb"
	"google.golang.org/grpc"
)

type service struct {
	Conf  config.Config
	Route iservice.IRoute
	Call  func(desc serviceobject.ServerDesc)
}

func (m *service) Listen() {
	var ip = m.Conf.IP
	var port = m.Conf.Port
	if ip == "" {
		ip = utils.NetworkIP()
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

	s := grpc.NewServer()
	defer s.Stop()
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: m.Conf.ServerName,
		HandlerType: new(gweb.IHandler),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "/",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					log.Println(srv, ctx)
					return nil, nil
				},
			},
		},
		Streams:  nil,
		Metadata: nil,
	}, nil)
	if m.Call != nil {
		m.Call(serviceobject.ServerDesc{
			ServerName: m.Conf.ServerName,
			Port:       port,
			IP:         ip,
		})
	}
	if err = s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
func New(
	conf config.Config,
	route iservice.IRoute,
	call func(desc serviceobject.ServerDesc),
) iservice.IGrpc {
	return &service{
		Conf:  conf,
		Route: route,
		Call:  call,
	}
}
