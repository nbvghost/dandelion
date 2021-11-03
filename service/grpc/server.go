package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/library/result"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strconv"

	"github.com/nbvghost/dandelion/service/iservice"

	"github.com/nbvghost/dandelion/service/serviceobject"
	"github.com/nbvghost/dandelion/utils"
	"google.golang.org/grpc"
)

type Option func(serviceobject.ServerDesc) error
type service struct {
	serviceobject.UnimplementedServerServer

	conf    *config.ServerConfig
	route   iservice.IRoute
	options []Option
}

func (m *service) Call(ctx context.Context, request *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error) {
	response, err := m.route.Handle(ctx, request)
	if err != nil {
		if v, ok := err.(*result.ActionResult); ok {
			return nil, status.Error(codes.DataLoss, v.Message)
		} else {
			return nil, status.Error(codes.DataLoss, err.Error())
		}
	}

	return response, nil
}

func (m *service) Listen() {
	var ip = m.conf.IP
	var port = m.conf.Port
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

	desc := serviceobject.ServerDesc{
		ServerName: m.conf.ServerName,
		Port:       port,
		IP:         ip,
	}
	for i := range m.options {
		if err = m.options[i](desc); err != nil {
			log.Fatalln(err)
		}
	}

	s := grpc.NewServer()
	defer s.Stop()

	serviceobject.RegisterServerServer(s, m)
	if err = s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
func New(
	conf *config.ServerConfig,
	route iservice.IRoute,
	options ...Option,
) iservice.IGrpc {
	return &service{
		conf:    conf,
		route:   route,
		options: options,
	}
}
