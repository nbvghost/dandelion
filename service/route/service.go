package route

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nbvghost/dandelion/library/result"
	icontext "github.com/nbvghost/dandelion/service/context"
	"github.com/nbvghost/dandelion/service/iservice"
	"github.com/nbvghost/dandelion/service/serviceobject"
	"github.com/nbvghost/tool/object"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type service struct {
	Routes map[string]Info
	redis  iservice.IRedis
}

var validate = validator.New()

func (m *service) Handle(ctx context.Context, desc *serviceobject.GrpcRequest) (*serviceobject.GrpcResponse, error) {
	if v, ok := m.Routes[desc.Route]; !ok {

		return nil, result.NewCodeWithError(result.NotFoundRoute, errors.New("没有找到路由"))

	} else {

		var err error

		handlerValue := reflect.New(v.HandlerType)

		if !v.WithoutAuth {
			if desc.UID == "" {
				return nil, result.NewCodeWithError(result.AuthError, errors.New("用户没有授权"))
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
		if paramQueryFieldName != "" {

			var values url.Values
			if values, err = url.ParseQuery(desc.Uri); err != nil {
				return nil, err
			}

			uriFieldValue := handlerValue.FieldByName(paramQueryFieldName)
			numField := uriFieldValue.NumField()
			for i := 0; i < numField; i++ {
				uriField := uriFieldValue.Type().Field(i)
				fieldValue := values.Get(uriField.Tag.Get("uri"))

				uriFieldValue.Field(i).Set(object.Convert(reflect.ValueOf(fieldValue), uriField.Type))
			}
		}

		if paramBodyFieldName != "" {

			if err = json.Unmarshal(desc.Body, handlerValue.Interface()); err != nil {
				return nil, err
			}
			if err = validate.Struct(handlerValue.Interface()); err != nil {
				return nil, result.NewCodeWithError(result.ValidateError, err)
			}
		}

		context := icontext.New(ctx, desc.UID, m.redis)

		apiHandler := handlerValue.Interface()
		var returnResult result.Result

		switch desc.HttpMethod {
		case http.MethodGet:
			if v, ok := apiHandler.(iservice.IHandlerGet); ok {
				returnResult, err = v.HandleGet(context)
			}
		case http.MethodPost:
			if v, ok := apiHandler.(iservice.IHandlerPost); ok {
				returnResult, err = v.HandlePost(context)
			}
		case http.MethodHead:
			if v, ok := apiHandler.(iservice.IHandlerHead); ok {
				returnResult, err = v.HandleHead(context)
			}
		case http.MethodPut:
			if v, ok := apiHandler.(iservice.IHandlerPut); ok {
				returnResult, err = v.HandlePut(context)
			}
		case http.MethodPatch:
			if v, ok := apiHandler.(iservice.IHandlerPatch); ok {
				returnResult, err = v.HandlePatch(context)
			}
		case http.MethodDelete:
			if v, ok := apiHandler.(iservice.IHandlerDelete); ok {
				returnResult, err = v.HandleDelete(context)
			}
		case http.MethodConnect:
			if v, ok := apiHandler.(iservice.IHandlerConnect); ok {
				returnResult, err = v.HandleConnect(context)
			}
		case http.MethodOptions:
			if v, ok := apiHandler.(iservice.IHandlerOptions); ok {
				returnResult, err = v.HandleOptions(context)
			}
		case http.MethodTrace:
			if v, ok := apiHandler.(iservice.IHandlerTrace); ok {
				returnResult, err = v.HandleTrace(context)
			}
		default:
			return nil, result.NewCodeWithError(result.HttpError, errors.New(fmt.Sprintf("错误的http方法:%s", desc.HttpMethod)))

		}

		if err != nil {
			return nil, err
		}
		if returnResult == nil {
			returnResult, err = apiHandler.(iservice.IHandler).Handle(context)
		}
		if err != nil {
			return nil, err
		}

		var data []byte
		if data, err = returnResult.Apply(); err != nil {
			return nil, result.NewError(err)
		}

		return &serviceobject.GrpcResponse{
			Error: 0,
			Data:  data,
		}, nil

	}

}

type Info struct {
	HandlerType reflect.Type
	WithoutAuth bool
}

func (m *service) RegisterRoute(path string, handler iservice.IHandler, withoutAuth ...bool) {
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

func New(redis iservice.IRedis) iservice.IRoute {
	return &service{Routes: map[string]Info{}, redis: redis}
}
