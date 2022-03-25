package route

import (
	"github.com/nbvghost/dandelion/constrain"
)

type RegisterRoute func(path string, handler constrain.IHandler, withoutAuth ...bool)
type RegisterView func(path string, handler constrain.IViewHandler, result constrain.IViewResult, withoutAuth ...bool)

type IRoute interface {
	RegisterRoute(path string, handler constrain.IHandler, withoutAuth ...bool)
	RegisterView(path string, handler constrain.IViewHandler, result constrain.IViewResult, withoutAuth ...bool)
	GetMappingCallback() constrain.IMappingCallback
	Handle(context constrain.IContext, isApi bool, route string, binddataFunc func(apiHandler interface{}) error) (bool, bool, interface{}, error)
	RegisterInterceptors(prefixPath string, interceptors ...constrain.IInterceptor)
}
