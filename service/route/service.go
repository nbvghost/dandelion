package route

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/nbvghost/dandelion/library/action"
	icontext "github.com/nbvghost/dandelion/library/context"
	"github.com/nbvghost/dandelion/library/gobext"
	"github.com/nbvghost/dandelion/library/handler"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/redis"
	"github.com/nbvghost/dandelion/service/serviceobject"
)

type service struct {
	Routes map[string]Info
	redis  redis.IRedis
}

var validate = validator.New()

func (m *service) encodingViewData(r result.ViewResult) ([]byte, string, error) {
	buffer := bytes.NewBuffer(nil)
	structName := gobext.GetStructName(r)
	r.SetName(structName)
	if err := gob.NewEncoder(buffer).Encode(r); err != nil {
		return nil, "", err
	}
	return buffer.Bytes(), structName, nil

}
func (m *service) Handle(ctx context.Context, desc *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error) {
	if v, ok := m.Routes[desc.Route]; !ok {

		return nil, action.NewCodeWithError(action.NotFoundRoute, errors.New("没有找到路由"))

	} else {

		var err error

		handlerValue := reflect.New(v.HandlerType)

		if !v.WithoutAuth {
			if desc.UID == "" {
				return nil, action.NewCodeWithError(action.AuthError, errors.New("用户没有授权"))
			}
		}

		var paramBodyFieldName string
		var paramQueryFieldName string

		if field, ok := v.HandlerType.FieldByName(desc.HttpMethod); ok {
			paramBodyFieldName = field.Name
		}
		if field, ok := v.HandlerType.FieldByName("QUERY"); ok {
			paramQueryFieldName = field.Name
		}
		{
			numField := v.HandlerType.NumField()
			for i := 0; i < numField; i++ {
				if tag, ok := v.HandlerType.Field(i).Tag.Lookup("param"); ok {
					if paramBodyFieldName == "" && strings.EqualFold(desc.HttpMethod, tag) {
						paramBodyFieldName = v.HandlerType.Field(i).Name
					}
					if paramQueryFieldName == "" && strings.EqualFold("QUERY", tag) {
						paramQueryFieldName = v.HandlerType.Field(i).Name
					}
				}
			}
		}
		var newFieldValue reflect.Value
		if paramQueryFieldName != "" {
			uriFieldValue := handlerValue.Elem().FieldByName(paramQueryFieldName)
			newFieldValue = reflect.New(uriFieldValue.Type())
			err = binding.Query.Bind(&http.Request{
				URL: &url.URL{
					RawQuery: desc.Query,
				},
			}, newFieldValue.Interface())
			if err != nil {
				return nil, err
			}
			uriFieldValue.Set(newFieldValue.Elem())
		}

		if paramBodyFieldName != "" {

			if err = json.Unmarshal(desc.Body, handlerValue.Interface()); err != nil {
				return nil, err
			}
			if err = validate.Struct(handlerValue.Interface()); err != nil {
				return nil, action.NewCodeWithError(action.ValidateError, err)
			}
		}

		var query interface{}
		if newFieldValue.IsValid() {
			query = newFieldValue.Interface()
		}

		context := icontext.New(ctx, desc.AppName, desc.UID, desc.Route, query, m.redis)

		apiHandler := handlerValue.Interface()
		var returnResult result.Result

		if v, ok := apiHandler.(handler.IViewHandler); ok {
			var r result.ViewResult
			r, err = v.Render(context)
			if err != nil {
				return nil, err
			}

			var data []byte
			var structName string

			if data, structName, err = m.encodingViewData(r); err != nil {
				return nil, action.NewError(err)
			}
			return &serviceobject.GrpcResponse{
				Error: 0,
				Data:  data,
				Name:  structName,
			}, nil
		} else {
			switch desc.HttpMethod {
			case http.MethodGet:
				if v, ok := apiHandler.(handler.IHandlerGet); ok {
					returnResult, err = v.HandleGet(context)
				}
			case http.MethodPost:
				if v, ok := apiHandler.(handler.IHandlerPost); ok {
					returnResult, err = v.HandlePost(context)
				}
			case http.MethodHead:
				if v, ok := apiHandler.(handler.IHandlerHead); ok {
					returnResult, err = v.HandleHead(context)
				}
			case http.MethodPut:
				if v, ok := apiHandler.(handler.IHandlerPut); ok {
					returnResult, err = v.HandlePut(context)
				}
			case http.MethodPatch:
				if v, ok := apiHandler.(handler.IHandlerPatch); ok {
					returnResult, err = v.HandlePatch(context)
				}
			case http.MethodDelete:
				if v, ok := apiHandler.(handler.IHandlerDelete); ok {
					returnResult, err = v.HandleDelete(context)
				}
			case http.MethodConnect:
				if v, ok := apiHandler.(handler.IHandlerConnect); ok {
					returnResult, err = v.HandleConnect(context)
				}
			case http.MethodOptions:
				if v, ok := apiHandler.(handler.IHandlerOptions); ok {
					returnResult, err = v.HandleOptions(context)
				}
			case http.MethodTrace:
				if v, ok := apiHandler.(handler.IHandlerTrace); ok {
					returnResult, err = v.HandleTrace(context)
				}
			default:
				return nil, action.NewCodeWithError(action.HttpError, errors.New(fmt.Sprintf("错误的http方法:%s", desc.HttpMethod)))

			}

			if err != nil {
				return nil, err
			}
			if returnResult == nil {
				returnResult, err = apiHandler.(handler.IHandler).Handle(context)
			}
			if err != nil {
				return nil, err
			}

			var data []byte
			if data, err = returnResult.Apply(context); err != nil {
				return nil, action.NewError(err)
			}

			return &serviceobject.GrpcResponse{
				Error: 0,
				Data:  data,
				Name:  "",
			}, nil
		}

	}

}

type Info struct {
	HandlerType reflect.Type
	WithoutAuth bool
}

func (m *service) RegisterRoute(path string, handler handler.IHandler, withoutAuth ...bool) {
	path = "/" + path
	if _, ok := m.Routes[path]; ok {
		panic(errors.New(fmt.Sprintf("存在相同的路由:%s", path)))
	}
	var _withoutAuth bool
	if len(withoutAuth) > 0 {
		_withoutAuth = withoutAuth[0]
	}
	m.Routes[path] = Info{
		HandlerType: reflect.TypeOf(handler).Elem(),
		WithoutAuth: _withoutAuth,
	}
}
func (m *service) RegisterView(path string, handler handler.IViewHandler, withoutAuth ...bool) {
	path = "/" + path
	if _, ok := m.Routes[path]; ok {
		panic(errors.New(fmt.Sprintf("存在相同的路由:%s", path)))
	}
	var _withoutAuth bool
	if len(withoutAuth) > 0 {
		_withoutAuth = withoutAuth[0]
	}
	m.Routes[path] = Info{
		HandlerType: reflect.TypeOf(handler).Elem(),
		WithoutAuth: _withoutAuth,
	}
}

func New(redis redis.IRedis) IRoute {
	return &service{Routes: map[string]Info{}, redis: redis}
}
