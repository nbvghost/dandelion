package grpcext

import (
	"github.com/nbvghost/dandelion/constrain"
	"google.golang.org/grpc"
)

type Register func(serviceDesc grpc.ServiceDesc, handlers []constrain.IGrpcHandler, withoutAuth ...bool)
type IGrpc interface {
	Register(serviceDesc grpc.ServiceDesc, handlers []constrain.IGrpcHandler, withoutAuth ...bool)
	AddCallback(callbacks ...constrain.IMappingCallback)
	Listen()
}

/*type IGrpcClient interface {
	Call(ctx context.Context, appName key.MicroServer, request *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error)
}
*/
