package httpext

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/domain/logger"
	"io"
	"log"
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
	"github.com/nbvghost/dandelion/server/route"
	"github.com/nbvghost/tool"
	"github.com/nbvghost/tool/encryption"
	"github.com/pkg/errors"
)

type httpMiddleware struct {
	serverName  string
	serviceList []constrain.IService

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
func (m *httpMiddleware) bindData(apiHandler any, ctx constrain.IContext, contextValue *contexext.ContextValue) error {

	var err error

	contentType := m.filterFlags(contextValue.Request.Header.Get("Content-Type"))

	body, err := io.ReadAll(contextValue.Request.Body)
	if err != nil {
		return err
	}

	if environments.Release() == false {
		//contextValue.Request.Body.Close()
		if strings.EqualFold(contentType, binding.MIMEMultipartPOSTForm) {
			log.Println(fmt.Sprintf("%s %s %s", m.serverName, contextValue.Request.URL.Path, contextValue.Request.Method))
		} else {
			log.Println(fmt.Sprintf("%s %s %s %s", m.serverName, contextValue.Request.URL.Path, contextValue.Request.Method, string(body)))
		}
	}

	contextValue.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var bodyValue reflect.Value
	var methodStruct reflect.Value
	{
		v := reflect.ValueOf(apiHandler)
		t := reflect.TypeOf(apiHandler).Elem()
		num := t.NumField()
		fieldIndex := -1
		for i := 0; i < num; i++ {
			method := t.Field(i).Tag.Get("method")
			if strings.EqualFold(method, contextValue.Request.Method) {
				fieldIndex = i
				break
			}
		}

		var hasBodyField bool

		if fieldIndex >= 0 {
			methodStruct = v.Elem().Field(fieldIndex)
		} else {
			methodStruct = v.Elem()
		}
		vvNum := methodStruct.NumField()
		for i := 0; i < vvNum; i++ {
			_, ok := methodStruct.Type().Field(i).Tag.Lookup("body")
			if ok {
				hasBodyField = true
				bodyValue = methodStruct.Field(i)
				if bodyValue.Kind() == reflect.Ptr {
					//bodyValue = bodyValue.Elem()
				}
				break
			}
		}

		if !hasBodyField {
			err = binding.Default(contextValue.Request.Method, contentType).Bind(contextValue.Request.Clone(contextValue.Request.Context()), methodStruct.Addr().Interface())
			if err != nil {
				ctx.Logger().With(zap.Error(err))
				return err
			}
		} else {
			if bodyValue.Kind() == reflect.Slice && bodyValue.Type().Elem().Kind() == reflect.Uint8 {
				bodyValue.SetBytes(body)
			}
		}
	}

	err = binding.Header.Bind(contextValue.Request, methodStruct.Addr().Interface())
	if err != nil {
		ctx.Logger().With(zap.Error(err))
		return err
	}

	uriVars := mux.Vars(contextValue.Request)
	if len(uriVars) > 0 {
		uriMap := make(map[string][]string)
		for uriKey := range uriVars {
			uriMap[uriKey] = []string{uriVars[uriKey]}
		}
		if len(uriMap) > 0 {
			err = binding.Uri.BindUri(uriMap, methodStruct.Addr().Interface())
			if err != nil {
				ctx.Logger().With(zap.Error(err))
				return err
			}
		}
	}

	err = binding.Query.Bind(contextValue.Request, methodStruct.Addr().Interface())
	if err != nil {
		ctx.Logger().With(zap.Error(err))
		return err
	}
	contextValue.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}

func (m *httpMiddleware) getSession(parentCtx context.Context, redisClient constrain.IRedis, token string) (Session, error) {
	var se Session
	se.Token = token
	var sessionText string
	sessionText, _ = redisClient.GetEx(parentCtx, redis.NewTokenKey(token), time.Minute*30)
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
			token = encryption.CipherEncrypter(encryption.NewSecretKey(fmt.Sprintf("%d", time.Now().UnixNano())), fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")))
			/*if environments.Release() {
				http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Hour * 23), SameSite: http.SameSiteNoneMode, Secure: true})
			} else {
				http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Hour * 23), SameSite: http.SameSiteNoneMode})
			}*/
			//http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Hour * 23), SameSite: http.SameSiteNoneMode, Secure: true})
			if environments.Release() {
				http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Hour * 24 * 15)})
			} else {
				http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Hour * 24 * 15), Secure: true, SameSite: http.SameSiteNoneMode})
			}
		}
	} else {
		token = cookie.Value
	}
	return token
}

func (m *httpMiddleware) CreateContext(etcdClient constrain.IEtcd, redisClient constrain.IRedis, router constrain.IRoute, w http.ResponseWriter, r *http.Request) constrain.IContext {
	var lang string
	domainPrefix, domainName := util.ParseDomain(r.Host)
	if len(domainPrefix) >= 1 {
		lang = domainPrefix[0]
	}
	if len(lang) == 0 || strings.EqualFold(lang, "dev") {
		lang = "en"
	}

	parentCtx := r.Context()

	if len(m.serviceList) > 0 {
		for i := range m.serviceList {
			service := m.serviceList[i]
			st := reflect.TypeOf(service).Elem()
			parentCtx = context.WithValue(parentCtx, reflect.New(st).Interface(), service)
		}
	}

	mode := r.Header.Get("Mode")

	Timeout, _ := strconv.ParseUint(r.Header.Get("Timeout"), 10, 64)
	TraceID := r.Header.Get("TraceID")
	if len(TraceID) == 0 {
		TraceID = tool.UUID()
	}

	appLogger, err := logger.CreateLogger("HttpContext", TraceID)
	if err != nil {
		return nil
	}

	//logger = logger.Named("HttpContext").With(zap.String("TraceID", TraceID)).With(zap.String("DomainName", domainName))
	//defer logger.Sync()

	token := m.getToken(w, r)
	var session Session
	session, err = m.getSession(parentCtx, redisClient, token)
	if err != nil {
		appLogger.Error("getSession", zap.Error(err))
	}

	var mappingCallback constrain.IMappingCallback
	if router != nil {
		mappingCallback = router.GetMappingCallback()
	}

	contextValue := &contexext.ContextValue{
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

	appLogger = appLogger.With(zap.String("Path", r.URL.String()))

	ctx := contexext.New(contexext.NewContext(parentCtx, contextValue), m.serverName, session.ID, r.URL.Path, mappingCallback, etcdClient, redisClient, session.Token, appLogger, key.Mode(mode))
	return ctx
}
func (m *httpMiddleware) Handle(ctx constrain.IContext, router constrain.IRoute, beforeViewRender constrain.IBeforeViewRender, afterViewRender constrain.IAfterViewRender, w http.ResponseWriter, r *http.Request) error {
	var err error
	ctxValue := contexext.FromContext(ctx)

	{
		if strings.EqualFold(ctx.Mode().String(), key.ModeRelease.String()) && !environments.Release() {
			err = errors.New("正式环境访问开发环境的服务")
			return err
		}
		if strings.EqualFold(ctx.Mode().String(), key.ModeDev.String()) && environments.Release() {
			err = errors.New("开发环境访问正式环境的服务")
			return err
		}
	}

	ctxValue.Response.Header().Set("Code", "0")

	var apiHandler any
	var routeInfo constrain.IRouteInfo
	contextValue := contexext.FromContext(ctx)

	//apiHandler := reflect.New(routeInfo.GetHandlerType()).Interface()
	//return apiHandler, routeInfo.GetWithoutAuth(), nil
	//apiHandler, withoutAuth, err = router.CreateHandle(ctxValue.IsApi, r)
	routeInfo, err = router.CreateHandle(ctxValue.IsApi, r)
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
		return err
	}

	///=withoutAuth = routeInfo.GetWithoutAuth()

	//创建新的handler处理器
	apiHandler = reflect.New(routeInfo.GetHandlerType()).Interface()

	var isWriteHttpResponse bool

	if ctxValue.IsApi {
		if beforeViewRender != nil {
			var canNext bool
			err = beforeViewRender.Api(ctx, r, w, func() {
				canNext = true
			})
			if err != nil {
				return err
			}
			if !canNext {
				return nil
			}
		}

		isWriteHttpResponse, err = router.ExecuteInterceptors(ctx, apiHandler)
		if err != nil {
			if isWriteHttpResponse {
				ctx.Logger().Error("ExecuteInterceptors Api", zap.Error(err))
				return nil
			}
			return err
		}

		//注入路由mapping
		err = router.GetMappingCallback().Mapping(ctx, apiHandler)
		if err != nil {
			return err
		}

		if err = m.bindData(apiHandler, ctx, contextValue); err != nil {
			return err
		}

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
		/*case http.MethodHead:
		if v, ok := apiHandler.(constrain.IHandlerHead); ok {
			handle = v.HandleHead
		}*/
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
		/*case http.MethodConnect:
		if v, ok := apiHandler.(constrain.IHandlerConnect); ok {
			handle = v.HandleConnect
		}*/
		/*case http.MethodOptions:
		if v, ok := apiHandler.(constrain.IHandlerOptions); ok {
			handle = v.HandleOptions
		}*/
		/*case http.MethodTrace:
		if v, ok := apiHandler.(constrain.IHandlerTrace); ok {
			handle = v.HandleTrace
		}*/
		default:
			return result.NewCodeWithMessage(result.HttpError, fmt.Sprintf("错误的http方法:%s", r.Method))

		}
		if handle == nil {
			return result.NewCodeWithMessage(result.HttpError, fmt.Sprintf("找不到http方法:%s的handle", r.Method))
		}
		var returnResult constrain.IResult
		returnResult, err = handle(ctx)
		if err == nil && returnResult == nil {
			returnResult = result.NewSuccess("OK")
		} else {
			if err != nil {
				returnResult = result.NewError(err)
			}
		}
		returnResult.Apply(ctx)

	} else {
		if beforeViewRender != nil {
			var canNext bool
			err = beforeViewRender.View(ctx, r, w, func() {
				canNext = true
			})
			if err != nil {
				return err
			}
			if !canNext {
				return nil
			}
		}

		isWriteHttpResponse, err = router.ExecuteInterceptors(ctx, apiHandler)
		if err != nil {
			if isWriteHttpResponse {
				ctx.Logger().Error("ExecuteInterceptors View", zap.Error(err))
				return nil
			}
			return err
		}

		//注入路由mapping
		err = router.GetMappingCallback().Mapping(ctx, apiHandler)
		if err != nil {
			return err
		}

		if err = m.bindData(apiHandler, ctx, contextValue); err != nil {
			return err
		}

		if v, ok := apiHandler.(constrain.IViewHandler); ok {
			var viewResult constrain.IViewResult
			viewResult, err = v.Render(ctx)
			if err != nil {
				return err
			}
			if viewResult == nil {
				return errors.New("没有返回数据")
			}
			if _, okk := viewResult.(*route.NoneView); okk {
				return nil
			}
			rr := viewResult.GetResult(ctx, v)
			if rr != nil {
				rr.Apply(ctx)
				return nil
			}

			viewBaseValue := reflect.ValueOf(viewResult).Elem().FieldByName("ViewBase")
			viewBase := viewBaseValue.Interface().(extends.ViewBase)

			htmlMeta := extends.NewHtmlMeta(contextValue.Lang, contextValue.RequestUrl)
			if viewBase.HtmlMetaCallback != nil {
				if err = viewBase.HtmlMetaCallback(viewBase, htmlMeta); err != nil {
					return err
				}
			}
			viewBase.HtmlMeta = htmlMeta
			viewBaseValue.Set(reflect.ValueOf(viewBase))

			if afterViewRender == nil {
				return errors.New("没找开视图渲染器")
			}
			if err = afterViewRender.Render(ctx, r, w, viewResult); err != nil {
				return err
			}
			return nil
			/*vr := &DefaultViewRender{}
			if err = vr.Render(ctx, r, w, viewResult); err != nil {
				return err
			}*/
		} else {
			return errors.Errorf("对视图访问的类型：%v不支持", apiHandler)
		}
	}
	return nil
}
