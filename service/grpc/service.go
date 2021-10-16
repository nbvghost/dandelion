package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/service/workobject"
	"github.com/nbvghost/dandelion/utils"
	"github.com/nbvghost/gweb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

type service struct {
	conf config.Config
	call func(desc workobject.ServerDesc)
}

func (m *service) Listen() {
	var ip string = m.conf.IP
	var port int = m.conf.Port
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
		ServiceName: m.conf.ServerName,
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
	if m.call != nil {
		m.call(workobject.ServerDesc{
			ServerName: m.conf.ServerName,
			Port:       port,
			IP:         ip,
		})
	}
	if err = s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
func New(conf config.Config, call func(desc workobject.ServerDesc)) *service {
	return &service{conf: conf, call: call}
}
