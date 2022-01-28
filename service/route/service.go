package route

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/util"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/library/gobext"
	"github.com/nbvghost/dandelion/service/redis"
	"github.com/nbvghost/dandelion/service/serviceobject"
)

var validate = validator.New()

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

type service struct {
	Routes     map[string]*Info
	ViewRoutes map[string]*Info

	redis        redis.IRedis
	callbacks    []constrain.ICallback
	interceptors map[string][]constrain.IInterceptor
	router       *mux.Router
}

func (m *service) encodingViewData(ctx constrain.IContext, r constrain.IViewResult) ([]byte, string, error) {
	buffer := bytes.NewBuffer(nil)
	structName := util.GetPkgPath(r)
	r.SetPkgPath(structName)
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
func (m *service) ExecuteInterceptor(context constrain.IContext, info constrain.IRouteInfo, writer http.ResponseWriter, request *http.Request) (bool, error) {
	for k := range m.interceptors {
		l := len(k)
		route := context.Route()
		if l > 0 && l <= len(route) {
			if k == route[:l] {
				for i := range m.interceptors[k] {
					return m.interceptors[k][i].Execute(context, info, writer, request)
				}
			}
		}
	}
	return false, nil

}

func (m *service) GetInfo(desc *serviceobject.GrpcRequest) (constrain.IRouteInfo, error) {
	var routeInfo *Info
	var ok bool
	var err error
	if desc.IsApi {
		if routeInfo, ok = m.Routes[desc.Route]; !ok {
			err = action.NewCodeWithError(action.NotFoundRoute, errors.New("没有找到路由"))
		}
	} else {
		if routeInfo, ok = m.ViewRoutes[desc.Route]; !ok {
			err = action.NewCodeWithError(action.NotFoundRoute, errors.New("没有找到路由"))
		}
	}

	return routeInfo, err
}

func (m *service) Handle(parent constrain.IContext, routeInfo constrain.IRouteInfo, desc *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error) {
	if !routeInfo.GetWithoutAuth() {
		if desc.UID == "" {
			return nil, action.NewCodeWithError(action.AuthError, errors.New("用户没有授权"))
		}
	}

	var err error
	var apiHandler interface{}

	apiHandler, err = Bind(routeInfo.GetHandlerType(), desc)
	if err != nil {
		return nil, err
	}

	for index := range m.callbacks {
		item := m.callbacks[index]
		if err = item.Before(parent, apiHandler); err != nil {
			return nil, err
		}
	}

	if desc.IsApi {
		var handle func(ctx constrain.IContext) (constrain.IResult, error)
		switch desc.HttpMethod {
		case http.MethodGet:
			if v, ok := apiHandler.(constrain.IHandler); ok {
				handle = v.Handle
			}
		case http.MethodPost:
			if v, ok := apiHandler.(constrain.IHandlerPost); ok {
				handle = v.HandlePost
			}
		case http.MethodHead:
			if v, ok := apiHandler.(constrain.IHandlerHead); ok {
				handle = v.HandleHead
			}
		case http.MethodPut:
			if v, ok := apiHandler.(constrain.IHandlerPut); ok {
				handle = v.HandlePut
			}
		case http.MethodPatch:
			if v, ok := apiHandler.(constrain.IHandlerPatch); ok {
				handle = v.HandlePatch
			}
		case http.MethodDelete:
			if v, ok := apiHandler.(constrain.IHandlerDelete); ok {
				handle = v.HandleDelete
			}
		case http.MethodConnect:
			if v, ok := apiHandler.(constrain.IHandlerConnect); ok {
				handle = v.HandleConnect
			}
		case http.MethodOptions:
			if v, ok := apiHandler.(constrain.IHandlerOptions); ok {
				handle = v.HandleOptions
			}
		case http.MethodTrace:
			if v, ok := apiHandler.(constrain.IHandlerTrace); ok {
				handle = v.HandleTrace
			}
		default:
			return nil, action.NewCodeWithError(action.HttpError, errors.New(fmt.Sprintf("错误的http方法:%s", desc.HttpMethod)))

		}
		if handle == nil {
			return nil, action.NewCodeWithError(action.HttpError, errors.New(fmt.Sprintf("找不到http方法:%s的handle", desc.HttpMethod)))
		}
		var returnResult constrain.IResult
		returnResult, err = handle(parent)
		if err == nil && returnResult == nil {
			return &serviceobject.GrpcResponse{
				Code: 0,
				Data: []byte("{}"),
				Name: "",
			}, nil
		}
		if err != nil {
			return nil, err
		}
		if returnResult == nil {
			//returnResult, err = apiHandler.(constrain.IHandler).Handle(parent)
			return nil, fmt.Errorf("对Api访问的类型：%v不支持", apiHandler)
		}
		if err != nil {
			return nil, err
		}

		var data []byte
		if data, err = returnResult.Apply(parent); err != nil {
			return nil, action.NewError(err)
		}

		return &serviceobject.GrpcResponse{
			Code: 0,
			Data: data,
			Name: "",
		}, nil
	} else {
		if v, ok := apiHandler.(constrain.IViewHandler); ok {
			var r constrain.IViewResult
			r, err = v.Render(parent)
			if err != nil {
				return nil, err
			}

			var data []byte
			var structName string

			if data, structName, err = m.encodingViewData(parent, r); err != nil {
				return nil, action.NewError(err)
			}
			for index := range m.callbacks {
				item := m.callbacks[index]
				if err = item.ViewAfter(parent, r); err != nil {
					return nil, err
				}
			}
			return &serviceobject.GrpcResponse{
				Code: 0,
				Data: data,
				Name: structName,
			}, nil
		} else {
			return nil, fmt.Errorf("对视图访问的类型：%v不支持", apiHandler)
		}
	}

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

func New(router *mux.Router, redis redis.IRedis, callbacks ...constrain.ICallback) IRoute {
	return &service{router: router, Routes: map[string]*Info{}, ViewRoutes: map[string]*Info{}, redis: redis, callbacks: callbacks, interceptors: make(map[string][]constrain.IInterceptor)}
}
