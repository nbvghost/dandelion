package grpcext

import (
	"context"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/server/serviceobject"
	"google.golang.org/grpc"
)

type Register func(serviceDesc grpc.ServiceDesc, handlers []constrain.IGrpcHandler, withoutAuth ...bool)
type IGrpc interface {
	Register(serviceDesc grpc.ServiceDesc, handlers []constrain.IGrpcHandler, withoutAuth ...bool)
	AddCallback(callbacks ...constrain.IMappingCallback)
	Listen()
}

type IGrpcClient interface {
	Call(ctx context.Context, appName key.MicroServerKey, request *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error)
}
