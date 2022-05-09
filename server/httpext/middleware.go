package httpext

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"strings"
)

var DefaultHttpMiddleware = &httpMiddleware{}

type httpMiddleware struct {
	context    constrain.IContext
	serverName string

	//session      Session
	//isApi        bool
	//pathTemplate string
	//customizeViewRender constrain.IViewRender
}

func (m *httpMiddleware) filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}
func (m *httpMiddleware) bindData(apiHandler any, r *http.Request) error {
	v := reflect.ValueOf(apiHandler)
	t := reflect.TypeOf(apiHandler).Elem()
	num := t.NumField()
	fieldIndex := -1
	for i := 0; i < num; i++ {
		method := t.Field(i).Tag.Get("method")
		if strings.EqualFold(method, r.Method) {
			fieldIndex = i
			break
		}
	}
	var vv reflect.Value
	if fieldIndex >= 0 {
		vv = v.Elem().Field(fieldIndex)
	} else {
		vv = v.Elem()
	}

	err := binding.Default(r.Method, m.filterFlags(r.Header.Get("Content-Type"))).Bind(r, vv.Addr().Interface())
	if err != nil {
		return err
	}
	err = binding.Header.Bind(r, vv.Addr().Interface())
	if err != nil {
		return err
	}

	uriVars := mux.Vars(r)
	if len(uriVars) > 0 {
		uriMap := make(map[string][]string)
		for uriKey := range uriVars {
			uriMap[uriKey] = []string{uriVars[uriKey]}
		}
		if len(uriMap) > 0 {
			err = binding.Uri.BindUri(uriMap, vv.Addr().Interface())
			if err != nil {
				return err
			}
		}
	}

	err = binding.Query.Bind(r, vv.Addr().Interface())
	if err != nil {
		return err
	}
	return err
}
func (m *httpMiddleware) Handle(ctx constrain.IContext, router constrain.IRoute, pathTemplate string, isApi bool, customizeViewRender constrain.IViewRender, w http.ResponseWriter, r *http.Request) (bool, error) {
	var err error
	var apiHandler any
	var broken, withoutAuth bool

	contextValue := contexext.FromContext(ctx)

	apiHandler, withoutAuth, err = router.CreateHandle(isApi, pathTemplate)
	if err != nil {
		/*if isApi {
			if err != nil {
				err = errors.Errorf(err.Error())
			} else {
				err = action.NewCodeWithError(action.NotFoundRoute, errors.Errorf("没有找到路由:%s", r.URL.Path))
			}
		} else {
			viewResult := route.NewViewResult(strings.TrimPrefix(r.URL.Path, "/"))

			viewBaseValue := reflect.ValueOf(viewResult).Elem().FieldByName("ViewBase")
			viewBase := viewBaseValue.Interface().(extends.ViewBase)

			htmlMeta := extends.NewHtmlMeta(contextValue.Lang, contextValue.RequestUrl)
			if viewBase.HtmlMetaCallback != nil {
				if err = viewBase.HtmlMetaCallback(viewBase, htmlMeta); err != nil {
					return false, err
				}
			}
			viewBase.HtmlMeta = htmlMeta
			viewBaseValue.Set(reflect.ValueOf(viewBase))

			if customizeViewRender != nil {
				if err = customizeViewRender.Render(ctx, r, w, viewResult); err != nil {
					return false, err
				}
				return true, nil
			}
			vr := &viewRender{}
			if err = vr.Render(ctx, r, w, viewResult); err != nil {
				return false, err
			}
		}*/
		return false, err
	}

	if err = m.bindData(apiHandler, r); err != nil {
		return false, err
	}

	broken, err = router.Handle(ctx, withoutAuth, apiHandler)
	if err != nil {
		if v, ok := err.(*action.ActionResult); ok {
			err = errors.Errorf(v.Message)
		} else {
			err = errors.Errorf(err.Error())
		}
		return false, err
	}
	if broken {
		return false, nil
	}

	if isApi {
		var handle func(ctx constrain.IContext) (constrain.IResult, error)
		switch r.Method {
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
			return false, action.NewCodeWithError(action.HttpError, errors.New(fmt.Sprintf("错误的http方法:%s", r.Method)))

		}
		if handle == nil {
			return false, action.NewCodeWithError(action.HttpError, errors.New(fmt.Sprintf("找不到http方法:%s的handle", r.Method)))
		}
		var returnResult constrain.IResult
		returnResult, err = handle(ctx)
		if err == nil && returnResult == nil {
			returnResult = &result.EmptyResult{}
		} else {
			if err != nil {
				return false, err
			}
			if returnResult == nil {
				return false, errors.Errorf("对Api访问的类型：%v不支持", apiHandler)
			}
			if err != nil {
				return false, err
			}
		}
		returnResult.Apply(ctx)

	} else {
		if v, ok := apiHandler.(constrain.IViewHandler); ok {
			var viewResult constrain.IViewResult
			viewResult, err = v.Render(ctx)
			if err != nil {
				return false, err
			}
			if viewResult == nil {
				return false, errors.New("没有返回数据")
			}
			result := viewResult.GetResult()
			if result != nil {
				result.Apply(ctx)
				return true, nil
			}

			viewBaseValue := reflect.ValueOf(viewResult).Elem().FieldByName("ViewBase")
			viewBase := viewBaseValue.Interface().(extends.ViewBase)

			htmlMeta := extends.NewHtmlMeta(contextValue.Lang, contextValue.RequestUrl)
			if viewBase.HtmlMetaCallback != nil {
				if err = viewBase.HtmlMetaCallback(viewBase, htmlMeta); err != nil {
					return false, err
				}
			}
			viewBase.HtmlMeta = htmlMeta
			viewBaseValue.Set(reflect.ValueOf(viewBase))

			if customizeViewRender != nil {
				if err = customizeViewRender.Render(ctx, r, w, viewResult); err != nil {
					return false, err
				}
				return true, nil
			}
			vr := &viewRender{}
			if err = vr.Render(ctx, r, w, viewResult); err != nil {
				return false, err
			}
		} else {
			return false, errors.Errorf("对视图访问的类型：%v不支持", apiHandler)
		}
	}
	return true, nil
}
