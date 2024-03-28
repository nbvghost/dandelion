package grpcext

import (
	"google.golang.org/grpc"

	"github.com/nbvghost/dandelion/constrain"
)

type Register func(serviceDesc grpc.ServiceDesc, handlers []constrain.IGrpcHandler, withoutAuth ...bool)

type IGrpc interface {
	Register(serviceDesc grpc.ServiceDesc, handlers []constrain.IGrpcHandler, withoutAuth ...bool)
	AddMapping(callback constrain.IMappingCallback)
	Server() *grpc.Server
	Listen()
}

/*type IGrpcClient interface {
	Call(ctx context.Context, appName key.MicroServer, request *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error)
}
*/
