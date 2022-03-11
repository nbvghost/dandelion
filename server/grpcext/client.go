package grpcext

import (
	"context"
	"github.com/nbvghost/dandelion/constrain"
	"log"

	"github.com/nbvghost/dandelion/server/serviceobject"
	"google.golang.org/grpc"
)

type client struct {
	etcd constrain.IEtcd
}

func (m *client) Call(ctx context.Context, appName constrain.MicroServerKey, request *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error) {
	endpoint, err := m.etcd.SelectServer(appName)
	if err != nil {
		return nil, err
	}
	log.Printf("call server addres:%s", endpoint)
	cl, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return serviceobject.NewServerClient(cl).Call(ctx, request)
}
func NewClient(etcd constrain.IEtcd) IGrpcClient {
	return &client{etcd: etcd}
}
