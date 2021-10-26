package grpc

import (
	"context"

	"github.com/nbvghost/dandelion/service/iservice"
	"github.com/nbvghost/dandelion/service/serviceobject"
	"google.golang.org/grpc"
)

type client struct {
	etcd iservice.IEtcdClient
}

func (m *client) Call(ctx context.Context, appName string, request serviceobject.GrpcRequest, response *serviceobject.GrpcResponse) error {
	endpoint := m.etcd.SelectServer(appName)
	cl, err := grpc.Dial(endpoint)
	if err != nil {
		return err
	}
	return cl.Invoke(ctx, "/", request, response)
}
func NewClient(etcd iservice.IEtcdClient) iservice.IGrpcClient {
	return &client{etcd: etcd}
}
