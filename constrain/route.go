package constrain

import "reflect"

type IRouteInfo interface {
	GetHandlerType() reflect.Type
	GetWithoutAuth() bool
}
type IRoute interface {
	RegisterRoute(path string, handler IHandler, withoutAuth ...bool)
	RegisterView(path string, handler IViewHandler, result IViewResult, withoutAuth ...bool)
	GetMappingCallback() IMappingCallback
	CreateHandle(isApi bool, pathTemplate string) (any, bool, error)
	Handle(context IContext, withoutAuth bool, routeHandler any) (bool, error)
	RegisterInterceptors(prefixPath string, interceptors ...IInterceptor)
}
