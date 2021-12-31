package route

import (
	"context"

	"github.com/nbvghost/dandelion/library/handler"
	"github.com/nbvghost/dandelion/service/serviceobject"
)

type IRoute interface {
	RegisterRoute(path string, handler handler.IHandler, withoutAuth ...bool)
	RegisterView(path string, handler handler.IViewHandler, withoutAuth ...bool)
	Handle(ctx context.Context, desc *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error)
}
