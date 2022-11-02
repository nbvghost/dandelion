package constrain

import (
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
)

type IHandlerPost interface {
	IHandler
	HandlePost(ctx IContext) (IResult, error)
}
type IHandlerHead interface {
	IHandler
	HandleHead(ctx IContext) (IResult, error)
}
type IHandlerPut interface {
	IHandler
	HandlePut(ctx IContext) (IResult, error)
}
type IHandlerPatch interface {
	IHandler
	HandlePatch(ctx IContext) (IResult, error)
}
type IHandlerDelete interface {
	IHandler
	HandleDelete(ctx IContext) (IResult, error)
}
type IHandlerConnect interface {
	IHandler
	HandleConnect(ctx IContext) (IResult, error)
}
type IHandlerOptions interface {
	IHandler
	HandleOptions(ctx IContext) (IResult, error)
}
type IHandlerTrace interface {
	IHandler
	HandleTrace(ctx IContext) (IResult, error)
}
type IHandler interface {
	Handle(ctx IContext) (IResult, error)
}
type IViewHandler interface {
	Render(ctx IContext) (IViewResult, error)
}
type IGrpcHandler interface {
	Handle(ctx IContext) (proto.Message, *status.Status)
}
