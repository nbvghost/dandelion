package route

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/util"

	"net/http"
	"reflect"

	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/library/gobext"
)

type Info struct {
	HandlerType reflect.Type
	WithoutAuth bool
}

func (m *Info) GetHandlerType() reflect.Type {
	return m.HandlerType
}

func (m *Info) GetWithoutAuth() bool {
	return m.WithoutAuth
}

func NewViewResult(name string) constrain.IViewResult {
	return &extends.ViewBase{Name: name}
}

type service struct {
	Routes     map[string]*Info
	ViewRoutes map[string]*Info

	redis           constrain.IRedis
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
	//todo r.SetPkgPath(structName)
	if err := gob.NewEncoder(buffer).Encode(r); err != nil {
		return nil, "", err
	}
	return buffer.Bytes(), structName, nil

}
func (m *service) RegisterInterceptors(prefixPath string, interceptors ...constrain.IInterceptor) {
	if len(prefixPath) == 0 {
		panic(fmt.Errorf("prefixPath 不能为空"))
	}
	if _, ok := m.interceptors[prefixPath]; !ok {
		m.interceptors[prefixPath] = make([]constrain.IInterceptor, 0)
	}
	m.interceptors[prefixPath] = append(m.interceptors[prefixPath], interceptors...)
}
func (m *service) CheckRoute(isApi bool, route string) (*Info, bool) {
	var routeInfo *Info
	var ok bool
	if isApi {
		if routeInfo, ok = m.Routes[route]; !ok {
			return nil, false
		}
	} else {
		if routeInfo, ok = m.ViewRoutes[route]; !ok {
			return nil, false
		}
	}
	return routeInfo, ok
}
func (m *service) Handle(context constrain.IContext, isApi bool, route string, binddataFunc func(apiHandler interface{}) error) (bool, bool, interface{}, error) {
	var hasRoute bool
	routeInfo, ok := m.CheckRoute(isApi, route)
	if !ok {
		return true, hasRoute, nil, action.NewCodeWithError(action.NotFoundRoute, fmt.Errorf("没有找到路由:%s", route))
	}
	hasRoute = true

	apiHandler := reflect.New(routeInfo.GetHandlerType()).Interface()
	if binddataFunc != nil {
		err := binddataFunc(apiHandler)
		if err != nil {
			return true, hasRoute, apiHandler, err
		}
	}

	if m.mappingCallback != nil {
		if err := m.mappingCallback.Before(context, apiHandler); err != nil {
			return true, hasRoute, nil, err
		}
	}

	if routeInfo.WithoutAuth {
		return false, hasRoute, apiHandler, nil
	}

	//interceptors 是有状态的，不支持mapping
	for k := range m.interceptors {
		l := len(k)
		route := context.Route()
		if l > 0 && l <= len(route) {
			if k == route[:l] {
				for i := range m.interceptors[k] {
					if m.mappingCallback != nil {
						if err := m.mappingCallback.Before(context, m.interceptors[k][i]); err != nil {
							return true, hasRoute, nil, err
						}
					}
					broken, err := m.interceptors[k][i].Execute(context)
					if err != nil {
						return true, hasRoute, nil, err
					}
					if broken {
						return true, hasRoute, nil, nil
					}
				}
			}
		}
	}

	if !routeInfo.GetWithoutAuth() {
		if context.UID() == 0 {
			return true, hasRoute, nil, action.NewCodeWithError(action.AuthError, errors.New("用户没有授权"))
		}
	}

	return false, hasRoute, apiHandler, nil

}

func (m *service) RegisterRoute(path string, handler constrain.IHandler, withoutAuth ...bool) {
	path = "/" + path
	if _, ok := m.Routes[path]; ok {
		panic(errors.New(fmt.Sprintf("存在相同的路由:%s", path)))
	}
	var _withoutAuth bool
	if len(withoutAuth) > 0 {
		_withoutAuth = withoutAuth[0]
	}
	m.router.HandleFunc("/api"+path, func(writer http.ResponseWriter, request *http.Request) {

	})
	m.Routes[path] = &Info{
		HandlerType: reflect.TypeOf(handler).Elem(),
		WithoutAuth: _withoutAuth,
	}
}

func (m *service) RegisterView(path string, handler constrain.IViewHandler, result constrain.IViewResult, withoutAuth ...bool) {
	path = "/" + path
	if _, ok := m.ViewRoutes[path]; ok {
		panic(errors.New(fmt.Sprintf("存在相同的路由:%s", path)))
	}
	var _withoutAuth bool
	if len(withoutAuth) > 0 {
		_withoutAuth = withoutAuth[0]
	}
	m.router.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {

	})
	m.ViewRoutes[path] = &Info{
		HandlerType: reflect.TypeOf(handler).Elem(),
		WithoutAuth: _withoutAuth,
	}
	gobext.Register(result)
}

func New(router *mux.Router, redis constrain.IRedis, mappingCallback constrain.IMappingCallback) IRoute {
	return &service{router: router, Routes: map[string]*Info{}, ViewRoutes: map[string]*Info{}, redis: redis, mappingCallback: mappingCallback, interceptors: make(map[string][]constrain.IInterceptor)}
}
