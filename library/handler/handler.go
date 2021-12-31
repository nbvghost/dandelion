package handler

import (
	"github.com/nbvghost/dandelion/library/context"
	"github.com/nbvghost/dandelion/library/result"
)

type IHandlerGet interface {
	IHandler
	HandleGet(ctx context.IContext) (result.Result, error)
}
type IHandlerPost interface {
	IHandler
	HandlePost(ctx context.IContext) (result.Result, error)
}
type IHandlerHead interface {
	IHandler
	HandleHead(ctx context.IContext) (result.Result, error)
}
type IHandlerPut interface {
	IHandler
	HandlePut(ctx context.IContext) (result.Result, error)
}
type IHandlerPatch interface {
	IHandler
	HandlePatch(ctx context.IContext) (result.Result, error)
}
type IHandlerDelete interface {
	IHandler
	HandleDelete(ctx context.IContext) (result.Result, error)
}
type IHandlerConnect interface {
	IHandler
	HandleConnect(ctx context.IContext) (result.Result, error)
}
type IHandlerOptions interface {
	IHandler
	HandleOptions(ctx context.IContext) (result.Result, error)
}
type IHandlerTrace interface {
	IHandler
	HandleTrace(ctx context.IContext) (result.Result, error)
}
type IHandler interface {
	Handle(ctx context.IContext) (result.Result, error)
}
type IViewHandler interface {
	Render(ctx context.IContext) (result.ViewResult, error)
}
