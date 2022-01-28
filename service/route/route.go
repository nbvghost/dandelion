package route

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/service/serviceobject"
	"net/http"
)

type RegisterRoute func(path string, handler constrain.IHandler, withoutAuth ...bool)
type RegisterView func(path string, handler constrain.IViewHandler, result constrain.IViewResult, withoutAuth ...bool)

type IRoute interface {
	RegisterRoute(path string, handler constrain.IHandler, withoutAuth ...bool)
	RegisterView(path string, handler constrain.IViewHandler, result constrain.IViewResult, withoutAuth ...bool)
	GetInfo(desc *serviceobject.GrpcRequest) (constrain.IRouteInfo, error)
	Handle(parent constrain.IContext, routeInfo constrain.IRouteInfo, desc *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error)

	ExecuteInterceptor(context constrain.IContext, routeInfo constrain.IRouteInfo, writer http.ResponseWriter, request *http.Request) (broken bool, err error)
	RegisterInterceptors(prefixPath string, interceptors ...constrain.IInterceptor)
}
