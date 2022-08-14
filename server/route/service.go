package route

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/nbvghost/dandelion/library/result"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/pkg/errors"

	"github.com/nbvghost/dandelion/library/gobext"
)

type RouteInfo struct {
	HandlerType reflect.Type
	WithoutAuth bool
}

func (m *RouteInfo) GetHandlerType() reflect.Type {
	return m.HandlerType
}

func (m *RouteInfo) GetWithoutAuth() bool {
	return m.WithoutAuth
}

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

type service struct {
	Routes     map[string]*RouteInfo
	ViewRoutes map[string]*RouteInfo

	//redis           constrain.IRedis
	mappingCallback constrain.IMappingCallback
	interceptors    map[string][]constrain.IInterceptor
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
func (m *service) RegisterInterceptors(prefixPath string, interceptors ...constrain.IInterceptor) {
	if len(prefixPath) == 0 {
		panic(errors.Errorf("prefixPath 不能为空"))
	}
	if _, ok := m.interceptors[prefixPath]; !ok {
		m.interceptors[prefixPath] = make([]constrain.IInterceptor, 0)
	}
	m.interceptors[prefixPath] = append(m.interceptors[prefixPath], interceptors...)
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
		return nil, errors.Errorf("没有找到路由映射:%s   =>   %s", r.URL.String(), pathTemplate)
	}
	//apiHandler := reflect.New(routeInfo.GetHandlerType()).Interface()
	//return apiHandler, routeInfo.GetWithoutAuth(), nil
	return routeInfo, nil
}
func (m *service) Handle(context constrain.IContext, withoutAuth bool, routeHandler any) (bool, error) {
	if m.mappingCallback != nil {
		m.mappingCallback.Before(context, routeHandler)
	}

	if withoutAuth {
		return false, nil
	}

	//interceptors 是有状态的，不支持mapping
	for k := range m.interceptors {
		l := len(k)
		route := context.Route()
		if l > 0 && l <= len(route) {
			if k == route[:l] {
				for i := range m.interceptors[k] {
					if m.mappingCallback != nil {
						m.mappingCallback.Before(context, m.interceptors[k][i])
					}
					broken, err := m.interceptors[k][i].Execute(context)
					if err != nil {
						return true, err
					}
					if broken {
						return true, nil
					}
				}
			}
		}
	}

	if !withoutAuth {
		if context.UID() == 0 {
			return true, result.NewCodeWithError(result.AuthError, errors.New("用户没有授权"))
		}
	}

	return false, nil

}

func (m *service) RegisterRoute(pathTemplate string, handler constrain.IHandler, withoutAuth ...bool) {
	if strings.EqualFold(pathTemplate, "*") {
		pathTemplate = "/api/"
	} else {
		pathTemplate = "/api/" + pathTemplate
	}
	if _, ok := m.Routes[pathTemplate]; ok {
		panic(errors.New(fmt.Sprintf("存在相同的路由:%s", pathTemplate)))
	}
	var _withoutAuth bool
	if len(withoutAuth) > 0 {
		_withoutAuth = withoutAuth[0]
	}
	m.router.HandleFunc(pathTemplate, func(writer http.ResponseWriter, request *http.Request) {

	})
	m.Routes[pathTemplate] = &RouteInfo{
		HandlerType: reflect.TypeOf(handler).Elem(),
		WithoutAuth: _withoutAuth,
	}
}

// RegisterView path 为 * 号时，匹配所有没有定义的路由
func (m *service) RegisterView(pathTemplate string, handler constrain.IViewHandler, result constrain.IViewResult, withoutAuth ...bool) {
	if strings.EqualFold(pathTemplate, "*") {
		pathTemplate = ""
	} else {
		pathTemplate = "/" + pathTemplate
	}
	if _, ok := m.ViewRoutes[pathTemplate]; ok {
		panic(errors.New(fmt.Sprintf("存在相同的路由:%s", pathTemplate)))
	}
	var _withoutAuth bool
	if len(withoutAuth) > 0 {
		_withoutAuth = withoutAuth[0]
	}
	m.router.HandleFunc(pathTemplate, func(writer http.ResponseWriter, request *http.Request) {

	})
	m.ViewRoutes[pathTemplate] = &RouteInfo{
		HandlerType: reflect.TypeOf(handler).Elem(),
		WithoutAuth: _withoutAuth,
	}
	gobext.Register(result)
}

func New(router *mux.Router, mappingCallback constrain.IMappingCallback) constrain.IRoute {
	return &service{router: router, Routes: map[string]*RouteInfo{}, ViewRoutes: map[string]*RouteInfo{}, mappingCallback: mappingCallback, interceptors: make(map[string][]constrain.IInterceptor)}
}
