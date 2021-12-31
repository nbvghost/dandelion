package grpc

import (
	"context"

	"github.com/nbvghost/dandelion/service/serviceobject"
)

type IGrpc interface {
	Listen()
}

type IGrpcClient interface {
	Call(ctx context.Context, appName string, request *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error)
}
