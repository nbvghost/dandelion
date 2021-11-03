package iservice

import (
	"context"
	"github.com/nbvghost/dandelion/service/serviceobject"
)

type IRoute interface {
	RegisterRoute(path string, handler IHandler, withoutAuth ...bool)
	Handle(ctx context.Context, desc *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error)
}
