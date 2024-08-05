package route

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/pkg/errors"
)

type RouteInfo struct {
	HandlerType reflect.Type
	//WithoutAuth bool
}

func (m *RouteInfo) GetHandlerType() reflect.Type {
	return m.HandlerType
}

/*func (m *RouteInfo) GetWithoutAuth() bool {
	return m.WithoutAuth
}*/

type emptyViewBase struct {
	extends.ViewBase
	Data any
}

func NewViewResult(name string, data any) constrain.IViewResult {
	return &emptyViewBase{
		ViewBase: extends.ViewBase{Name: name},
		Data:     data,
	}
}

type NoneView struct {
	extends.ViewBase
}

func NewNoneViewResult() constrain.IViewResult {
	return &NoneView{
		ViewBase: extends.ViewBase{},
	}
}

type scopeInterceptor struct {
	Interceptors []constrain.IInterceptor
	ExcludedPath []string
}

type service struct {
	Routes     map[string]*RouteInfo
	ViewRoutes map[string]*RouteInfo

	//redis           constrain.IRedis
	mappingCallback constrain.IMappingCallback
	interceptors    map[string]scopeInterceptor
	router          *mux.Router
}

func (m *service) GetMappingCallback() constrain.IMappingCallback {
	return m.mappingCallback
}
func (m *service) encodingViewData(ctx constrain.IContext, r constrain.IViewResult) ([]byte, string, error) {
	buffer := bytes.NewBuffer(nil)
	structName := util.GetPkgPath(r)
	//todo r.SetPkgPath(structName) 问题：grpc 到http 模板渲染时，无法得到struct结构，具体值的类型无法确定（map也不行），目前内部服务采用http,考虑内部可能是grpc的。
	if err := gob.NewEncoder(buffer).Encode(r); err != nil {
		return nil, "", err
	}
	return buffer.Bytes(), structName, nil

}
func (m *service) RegisterInterceptors(prefixPath string, excluded []string, interceptors ...constrain.IInterceptor) {
	if len(prefixPath) == 0 {
		panic(errors.Errorf("prefixPath 不能为空"))
	}
	m.interceptors[prefixPath] = scopeInterceptor{
		Interceptors: interceptors,
		ExcludedPath: append(excluded, "/404"),
	}
}
func (m *service) CheckRoute(isApi bool, route string) (*RouteInfo, bool) {
	var routeInfo *RouteInfo
	var ok bool
	if isApi {
		if routeInfo, ok = m.Routes[route]; !ok {
			if routeInfo, ok = m.Routes["/api/"]; !ok {
				return nil, false
			}
		}
	} else {
		if routeInfo, ok = m.ViewRoutes[route]; !ok {
			if routeInfo, ok = m.ViewRoutes[""]; !ok {
				return nil, false
			}
		}
	}
	return routeInfo, ok
}
func (m *service) CreateHandle(isApi bool, r *http.Request) (constrain.IRouteInfo, error) {
	var pathTemplate string
	var err error

	/*pathTemplate, err = getPathTemplate(r)
	if err != nil {
		return nil, err
	}*/

	currentRoute := mux.CurrentRoute(r)
	pathTemplate, err = currentRoute.GetPathTemplate()
	if err != nil {
		return nil, err
	}

	routeInfo, ok := m.CheckRoute(isApi, pathTemplate)
	if !ok {
		return nil, errors.Errorf("No path found %s", r.URL.String())
	}
	//apiHandler := reflect.New(routeInfo.GetHandlerType()).Interface()
	//return apiHandler, routeInfo.GetWithoutAuth(), nil
	return routeInfo, nil
}
func (m *service) ExecuteInterceptors(context constrain.IContext, routeHandler any) (bool,error) {
	//todo 权限控制采用拦截器来处理
	/*if withoutAuth {
		return false, nil
	}*/

	//

	//interceptors 是有状态的，不支持mapping
	for k := range m.interceptors {
		l := len(k)
		route := context.Route()
		if l > 0 && l <= len(route) {
			if strings.EqualFold(k, route[:l]) {
				//判断是否要执行
				isExcluded := false
				for _, excludedPath := range m.interceptors[k].ExcludedPath {
					excludedPathLen := len(excludedPath)
					if excludedPathLen > 0 && excludedPathLen <= len(route) {
						routePath := route[:excludedPathLen]
						if strings.EqualFold(excludedPath, routePath) {
							isExcluded = true
							break
						}
					}
				}

				if !isExcluded {
					//执行拦截器
					for i := range m.interceptors[k].Interceptors {
						if m.mappingCallback != nil {
							err := m.mappingCallback.Mapping(context, m.interceptors[k].Interceptors[i])
							if err != nil {
								return false,err
							}
						}
						isWriteHttpResponse,err := m.interceptors[k].Interceptors[i].Execute(context)
						if err != nil {
							return isWriteHttpResponse,err
						}
						//return nil
					}
				}

			}
		}
	}

	//通过拦截器处理
	/*if !withoutAuth {
		if context.UID() == 0 {
			return true, result.NewCodeWithMessage(result.AuthError, "用户没有授权")
		}
	}*/

	return false,nil

}

func (m *service) RegisterRoute(pathTemplate string, handler constrain.IHandler) {

	if strings.EqualFold(pathTemplate, "*") {
		pathTemplate = "/api/"
		//pathTemplate = filepath.ToSlash(filepath.Join("/", "api", "/"))
	} else {
		pathTemplate = "/api/" + strings.TrimLeft(pathTemplate, "/")
		//pathTemplate = filepath.ToSlash(filepath.Join("/", "api", strings.Trim(pathTemplate, "/")))
	}
	if _, ok := m.Routes[pathTemplate]; ok {
		panic(errors.New(fmt.Sprintf("存在相同的路由:%s", pathTemplate)))
	}
	/*var _withoutAuth bool
	if len(withoutAuth) > 0 {
		_withoutAuth = withoutAuth[0]
	}*/
	m.router.HandleFunc(pathTemplate, func(writer http.ResponseWriter, request *http.Request) {

	})
	m.Routes[pathTemplate] = &RouteInfo{
		HandlerType: reflect.TypeOf(handler).Elem(),
		//WithoutAuth: _withoutAuth,
	}
}

// RegisterView path 为 * 号时，匹配所有没有定义的路由
func (m *service) RegisterView(pathTemplate string, handler constrain.IViewHandler) {

	if strings.EqualFold(pathTemplate, "*") {
		pathTemplate = ""
	} else {
		pathTemplate = "/" + strings.TrimLeft(pathTemplate, "/")
	}
	if _, ok := m.ViewRoutes[pathTemplate]; ok {
		panic(errors.New(fmt.Sprintf("存在相同的路由:%s", pathTemplate)))
	}
	/*var _withoutAuth bool
	if len(withoutAuth) > 0 {
		_withoutAuth = withoutAuth[0]
	}*/
	m.router.HandleFunc(pathTemplate, func(writer http.ResponseWriter, request *http.Request) {

	})
	m.ViewRoutes[pathTemplate] = &RouteInfo{
		HandlerType: reflect.TypeOf(handler).Elem(),
		//WithoutAuth: _withoutAuth,
	}
	//gobext.Register(result)
}

func New(router *mux.Router, mappingCallback constrain.IMappingCallback) constrain.IRoute {
	return &service{router: router, Routes: map[string]*RouteInfo{}, ViewRoutes: map[string]*RouteInfo{}, mappingCallback: mappingCallback, interceptors: make(map[string]scopeInterceptor)}
}
