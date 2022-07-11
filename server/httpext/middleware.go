package httpext

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/environments"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/server/redis"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool"
	"github.com/nbvghost/tool/encryption"
	"github.com/pkg/errors"
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

func (m *httpMiddleware) getSession(parentCtx context.Context, redisClient constrain.IRedis, token string) (Session, error) {
	var se Session
	se.Token = token
	var sessionText string
	sessionText, _ = redisClient.GetEx(parentCtx, redis.NewTokenKey(token), time.Minute*10)
	if sessionText != "" {
		if err := json.Unmarshal([]byte(sessionText), &se); err != nil {
			return se, err
		}
	}
	return se, nil
}
func (m *httpMiddleware) getToken(w http.ResponseWriter, r *http.Request) string {
	var err error
	var cookie *http.Cookie
	var token string

	cookie, err = r.Cookie("token")
	if err != nil || strings.EqualFold(cookie.Value, "") {
		token = r.Header.Get("X-Token")
		if len(token) == 0 {
			token = encryption.CipherEncrypter(encryption.NewSecretKey(conf.Config.SecureKey), fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")))
			http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Hour * 24)})
		}
	} else {
		token = cookie.Value
	}
	return token
}

func (m *httpMiddleware) CreateContent(redisClient constrain.IRedis, router constrain.IRoute, w http.ResponseWriter, r *http.Request) constrain.IContext {
	var lang string
	domainPrefix, domainName := util.ParseDomain(r.Host)
	if len(domainPrefix) >= 1 {
		lang = domainPrefix[0]
	}
	if len(lang) == 0 || strings.EqualFold(lang, "dev") {
		lang = "en"
	}

	parentCtx := r.Context()

	mode := r.Header.Get("Mode")

	Timeout, _ := strconv.ParseUint(r.Header.Get("Timeout"), 10, 64)
	TraceID := r.Header.Get("TraceID")
	if len(TraceID) == 0 {
		TraceID = tool.UUID()
	}

	var err error
	var logger *zap.Logger
	if environments.Release() {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}

	logger = logger.Named("HttpContext").With(zap.String("TraceID", TraceID))
	//defer logger.Sync()

	token := m.getToken(w, r)
	var session Session
	session, err = m.getSession(parentCtx, redisClient, token)
	if err != nil {
		logger.Error("getSession", zap.Error(err))
	}

	var mappingCallback constrain.IMappingCallback
	if router != nil {
		mappingCallback = router.GetMappingCallback()
	}

	contextValue := &contexext.ContextValue{
		Mapping:    mappingCallback,
		Response:   w,
		Request:    r,
		Timeout:    Timeout,
		DomainName: domainName,
		Lang:       lang,
		RequestUrl: util.GetFullUrl(r),
		//PathTemplate: pathTemplate,
		Query: r.URL.Query(),
	}

	{
		var apiPath = "/api/"
		var requestUrlPath = r.URL.Path
		if len(requestUrlPath) >= len(apiPath) && strings.EqualFold(requestUrlPath[0:len(apiPath)], apiPath) {
			contextValue.IsApi = true
		}
	}

	logger = logger.With(zap.String("Path", r.URL.String()))

	ctx := contexext.New(contexext.NewContext(parentCtx, contextValue), m.serverName, session.ID, r.URL.Path, redisClient, session.Token, logger, key.Mode(mode))
	return ctx
}
func (m *httpMiddleware) Handle(ctx constrain.IContext, router constrain.IRoute, customizeViewRender constrain.IViewRender, w http.ResponseWriter, r *http.Request) (bool, error) {
	var err error
	ctxValue := contexext.FromContext(ctx)
	var pathTemplate string
	pathTemplate, err = getPathTemplate(r)
	if err != nil {
		return false, err
	}
	{
		if strings.EqualFold(ctx.Mode().String(), key.ModeRelease.String()) && !environments.Release() {
			err = errors.New("正式环境访问开发环境的服务")
			return false, err
		}
		if strings.EqualFold(ctx.Mode().String(), key.ModeDev.String()) && environments.Release() {
			err = errors.New("开发环境访问正式环境的服务")
			return false, err
		}
	}

	w.Header().Set("Code", "0")

	var apiHandler any
	var broken, withoutAuth bool

	contextValue := contexext.FromContext(ctx)

	apiHandler, withoutAuth, err = router.CreateHandle(ctxValue.IsApi, pathTemplate)
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
		return false, err
	}
	if broken {
		return false, nil
	}

	if ctxValue.IsApi {
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
			return false, result.NewCodeWithError(result.HttpError, errors.New(fmt.Sprintf("错误的http方法:%s", r.Method)))

		}
		if handle == nil {
			return false, result.NewCodeWithError(result.HttpError, errors.New(fmt.Sprintf("找不到http方法:%s的handle", r.Method)))
		}
		var returnResult constrain.IResult
		returnResult, err = handle(ctx)
		if err == nil && returnResult == nil {
			returnResult = result.NewError(err)
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
			result := viewResult.GetResult(ctx)
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
