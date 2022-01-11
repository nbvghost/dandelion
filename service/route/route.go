package route

import (
	"github.com/gin-gonic/gin"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/service/serviceobject"
)

type RegisterRoute func(path string, handler constrain.IHandler, withoutAuth ...bool)
type RegisterView func(path string, handler constrain.IViewHandler, result constrain.IViewResult, withoutAuth ...bool)

type IRoute interface {
	RegisterRoute(path string, handler constrain.IHandler, withoutAuth ...bool)
	RegisterView(path string, handler constrain.IViewHandler, result constrain.IViewResult, withoutAuth ...bool)
	Handle(parent constrain.IContext, routeInfo constrain.IRouteInfo, desc *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error)
	ExecuteInterceptor(context constrain.IContext, routeInfo constrain.IRouteInfo, ginContext *gin.Context) (broken bool, err error)
	RegisterInterceptors(prefixPath string, interceptors ...constrain.IInterceptor)
	GetInfo(desc *serviceobject.GrpcRequest) (constrain.IRouteInfo, error)
}
