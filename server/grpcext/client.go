package grpcext

import (
	"github.com/nbvghost/dandelion/constrain"
)

type client struct {
	etcd constrain.IEtcd
}

/*func (m *client) Call(ctx context.Context, appName key.MicroServer, request *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error) {
	endpoint, err := m.etcd.SelectInsideServer(appName)
	if err != nil {
		return nil, err
	}
	log.Printf("call server addres:%s", endpoint)
	cl, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return serviceobject.NewServerClient(cl).Call(ctx, request), nil
}
func NewClient(etcd constrain.IEtcd) IGrpcClient {
	return &client{etcd: etcd}
}
*/
