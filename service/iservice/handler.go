package iservice

import (
	"github.com/nbvghost/dandelion/library/result"
)

type IHandlerGet interface {
	IHandler
	HandleGet(ctx IContext) (result.Result, error)
}
type IHandlerPost interface {
	IHandler
	HandlePost(ctx IContext) (result.Result, error)
}
type IHandlerHead interface {
	IHandler
	HandleHead(ctx IContext) (result.Result, error)
}
type IHandlerPut interface {
	IHandler
	HandlePut(ctx IContext) (result.Result, error)
}
type IHandlerPatch interface {
	IHandler
	HandlePatch(ctx IContext) (result.Result, error)
}
type IHandlerDelete interface {
	IHandler
	HandleDelete(ctx IContext) (result.Result, error)
}
type IHandlerConnect interface {
	IHandler
	HandleConnect(ctx IContext) (result.Result, error)
}
type IHandlerOptions interface {
	IHandler
	HandleOptions(ctx IContext) (result.Result, error)
}
type IHandlerTrace interface {
	IHandler
	HandleTrace(ctx IContext) (result.Result, error)
}
type IHandler interface {
	Handle(ctx IContext) (result.Result, error)
}
