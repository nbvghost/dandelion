package constrain

import (
	"net/http"
	"reflect"
)

type IRouteInfo interface {
	GetHandlerType() reflect.Type
	//GetWithoutAuth() bool
}
type IRoute interface {
	RegisterRoute(path string, handler IHandler)
	RegisterView(path string, handler IViewHandler)
	GetMappingCallback() IMappingCallback
	CreateHandle(isApi bool, r *http.Request) (IRouteInfo, error)
	ExecuteInterceptors(context IContext, routeHandler any) (bool, error)
	RegisterInterceptors(prefixPath string, excludedPath []string, interceptors ...IInterceptor)
}
