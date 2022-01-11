package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/service/redis"
	"github.com/nbvghost/dandelion/service/route"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/url"
	"strconv"

	"github.com/nbvghost/dandelion/service/serviceobject"
	"github.com/nbvghost/dandelion/utils"
	"google.golang.org/grpc"
)

type Option func(serviceobject.ServerDesc) error
type service struct {
	serviceobject.UnimplementedServerServer

	server  config.MicroServerConfig
	route   route.IRoute
	redis   redis.IRedis
	options []Option
}

func (m *service) Call(ctx context.Context, request *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error) {

	query, err := url.ParseQuery(request.Query)
	if err != nil {
		return nil, err
	}

	currentContext := contexext.New(ctx, request.AppName, request.UID, request.Route, query, m.redis)

	info, err := m.route.GetInfo(request)
	if err != nil {
		return nil, err
	}

	response, err := m.route.Handle(currentContext, info, request)
	if err != nil {
		if v, ok := err.(*action.ActionResult); ok {
			return nil, status.Error(codes.DataLoss, v.Message)
		} else {
			return nil, status.Error(codes.DataLoss, err.Error())
		}
	}

	return response, nil
}

func (m *service) Listen() {
	var ip = m.server.IP
	var port = m.server.Port
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
		Name: m.server.Name,
		Port: port,
		IP:   ip,
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
	server config.MicroServerConfig,
	route route.IRoute,
	options ...Option,
) IGrpc {
	return &service{
		server:  server,
		route:   route,
		options: options,
	}
}
