package httpext

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/environments"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/server/route"
	"html/template"
	"reflect"
	"strconv"

	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/library/contexext"

	"github.com/nbvghost/dandelion/library/funcmap"
	"github.com/nbvghost/dandelion/server/redis"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool/encryption"
)

const defaultMemory = 32 << 20

type httpServer struct {
	serverName string
	listenAddr string
	engine     *mux.Router
}
type HttpMiddleware struct {
	currentRoute *mux.Route
	context      context.Context
	redis        constrain.IRedis
	serverName   string

	session      gweb.Session
	isApi        bool
	pathTemplate string
	token        string
	route        route.IRoute
	viewRender   constrain.IViewRender
}

func (m *HttpMiddleware) Path(w http.ResponseWriter, r *http.Request) (bool, error) {
	//log.Println(writer, request)
	//log.Println(mux.CurrentRoute(r).GetHostTemplate())
	//log.Println(mux.CurrentRoute(r).GetPathTemplate())
	//log.Println(re.GetPathRegexp())
	//log.Println(re.GetName())
	//log.Println(re.GetHandler())
	//log.Println(re.GetQueriesRegexp())
	//log.Println(re.GetQueriesTemplates())

	pathTemplate, err := m.currentRoute.GetPathTemplate()
	if err != nil {
		return false, err
	}

	paths := strings.Split(pathTemplate, "/")
	var path string
	if len(paths) >= 2 {
		endpointName := paths[1]
		if endpointName == "api" {
			//api请求

			path = "/" + strings.Join(paths[2:], "/")
			m.isApi = true
			m.pathTemplate = path
			return true, nil
		}
	}
	m.isApi = false
	m.pathTemplate = pathTemplate
	return true, nil

}
func (m *HttpMiddleware) Cookie(w http.ResponseWriter, r *http.Request) (bool, error) {
	var err error
	var cookie *http.Cookie
	var token string

	cookie, err = r.Cookie("token")
	if err != nil || strings.EqualFold(cookie.Value, "") {
		token = encryption.CipherEncrypter(encryption.NewSecretKey(conf.Config.SecureKey), fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")))
		http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Hour * 24)})
	} else {
		token = cookie.Value
	}
	m.token = token

	var sessionText string
	sessionText, _ = m.redis.GetEx(m.context, redis.NewTokenKey(token), time.Minute*10)
	if sessionText != "" {
		if err = json.Unmarshal([]byte(sessionText), &m.session); err != nil {
			return false, err
		}
	}
	return true, nil
}
func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func (m *HttpMiddleware) Handle(w http.ResponseWriter, r *http.Request) (bool, error) {
	var err error

	domainName := util.ParseDomain(r.Host)

	Timeout, _ := strconv.ParseUint(r.Header.Get("Timeout"), 10, 64)

	ctx := contexext.New(contexext.NewContext(&contexext.ContextValue{
		Mapping:    m.route.GetMappingCallback(),
		Response:   w,
		Request:    r,
		Timeout:    Timeout,
		DomainName: domainName,
		IsApi:      m.isApi,
	}), m.serverName, m.session.ID, m.pathTemplate, r.URL.Query(), m.redis, m.token)

	//todo ctx.Attributes().Put("Token", m.token)

	var broken, hasRoute bool
	var apiHandler interface{}
	broken, hasRoute, apiHandler, err = m.route.Handle(ctx, m.isApi, m.pathTemplate, func(apiHandler interface{}) error {
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

		err := binding.Default(r.Method, filterFlags(r.Header.Get("Content-Type"))).Bind(r, vv.Addr().Interface())
		if err != nil {
			return err
		}
		err = binding.Header.Bind(r, vv.Addr().Interface())
		if err != nil {
			return err
		}

		uriVars := mux.Vars(r)
		uriMap := make(map[string][]string)
		for uriKey := range uriVars {
			uriMap[uriKey] = []string{uriVars[uriKey]}
		}
		err = binding.Uri.BindUri(uriMap, vv.Addr().Interface())
		if err != nil {
			return err
		}
		err = binding.Query.Bind(r, vv.Addr().Interface())
		if err != nil {
			return err
		}
		return err
	})
	if !hasRoute {
		if m.isApi {
			if err != nil {
				err = fmt.Errorf(err.Error())
			} else {
				err = action.NewCodeWithError(action.NotFoundRoute, fmt.Errorf("没有找到路由:%s", r.URL.Path))
			}
		} else {
			viewResult := route.NewViewResult(strings.TrimPrefix(r.URL.Path, "/"))
			if m.viewRender != nil {
				if err = m.viewRender.Render(ctx, r, w, viewResult); err != nil {
					return false, err
				}
				return true, nil
			}
			vr := &viewRender{}
			if err = vr.Render(ctx, r, w, viewResult); err != nil {
				return false, err
			}
		}
		return true, nil
	}
	if err != nil {
		if v, ok := err.(*action.ActionResult); ok {
			err = fmt.Errorf(v.Message)
		} else {
			err = fmt.Errorf(err.Error())
		}
		return false, err
	}
	if broken {
		return false, nil
	}

	if m.isApi {
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
				return false, fmt.Errorf("对Api访问的类型：%v不支持", apiHandler)
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

			if m.viewRender != nil {
				if err = m.viewRender.Render(ctx, r, w, viewResult); err != nil {
					return false, err
				}
				return true, nil
			}
			vr := &viewRender{}
			if err = vr.Render(ctx, r, w, viewResult); err != nil {
				return false, err
			}
		} else {
			return false, fmt.Errorf("对视图访问的类型：%v不支持", apiHandler)
		}
	}
	return true, nil
}
func (m *HttpMiddleware) IsAPI() bool {
	return m.isApi
}
func (m *HttpMiddleware) Defer(w http.ResponseWriter, r *http.Request, err error) {
	var bytes []byte

	if err != nil {
		if m.isApi {
			//ginContext.JSON(http.StatusOK, action.NewError(err))
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			bytes, err = json.Marshal(action.NewError(err))
			if err != nil {
				log.Println(err)
			}
			w.Write(bytes)
		} else {
			t, errTemplate := template.New("").Parse(html_404)
			if errTemplate == nil {
				d := map[string]interface{}{
					"ErrorText": err.Error(),
					"Mode":      environments.Release(),
					"Stack":     string(debug.Stack()),
				}
				errTemplate = t.Execute(w, d)
			}
			if errTemplate != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errTemplate.Error()))
			}
		}
	}
}

type viewRender struct {
}

func (v *viewRender) Render(context constrain.IContext, request *http.Request, writer http.ResponseWriter, viewData interface{}) error {
	var err error
	var fileByte []byte

	vd := viewData.(constrain.IViewResult)

	viewName := vd.GetName()
	if len(viewName) > 0 {
		dir, _ := filepath.Split(context.Route())
		fileByte, err = ioutil.ReadFile(fmt.Sprintf("view/%s/%s.html", dir, viewName))
	} else {
		fileByte, err = ioutil.ReadFile(fmt.Sprintf("view/%s.html", strings.TrimSuffix(context.Route(), "/")))
		if err != nil {
			if _, ok := err.(*fs.PathError); ok {
				fileByte, err = ioutil.ReadFile(fmt.Sprintf("view/%s.html", "index"))
			}
		}

	}

	if err != nil {
		return err
	}

	var t *template.Template
	t, err = template.New("").Funcs(funcmap.NewFuncMap(context)).Parse(string(fileByte))
	if err != nil {
		return err
	}

	filenames, err := filepath.Glob(fmt.Sprintf("view/template/*.gohtml"))
	if err != nil {
		return err
	}
	if len(filenames) > 0 {
		t, err = t.ParseFiles(filenames...)
		if err != nil {
			return err
		}
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(writer, map[string]interface{}{
		"Query": context.Query(),
		"Data":  viewData,
	})
	if err != nil {
		return err
	}
	return nil
}

type MiddlewareFunc func(context context.Context, currentRoute *mux.Route, customizeViewRender constrain.IViewRender, serverName string) constrain.IMiddleware

func NewHttpMiddleware(context context.Context, serverName string, currentRoute *mux.Route, route route.IRoute, redis constrain.IRedis, viewRender constrain.IViewRender) constrain.IMiddleware {
	return &HttpMiddleware{
		currentRoute: currentRoute,
		redis:        redis,
		context:      context,
		serverName:   serverName,
		route:        route,
		viewRender:   viewRender,
	}
}

func (m *httpServer) Use(middlewareRouter *mux.Router, customizeViewRender constrain.IViewRender, middleware MiddlewareFunc) {
	middlewareRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var middleware = middleware(context.TODO(), mux.CurrentRoute(r), customizeViewRender, m.serverName)

		var isNext bool
		var err error
		defer func() {
			middleware.Defer(w, r, err)
		}()

		if isNext, err = middleware.Path(w, r); err != nil {
			return
		}
		if !isNext {
			return
		}

		v, ok := middleware.(*HttpMiddleware)
		if ok {
			v.pathTemplate = r.URL.Path //todo 临时处理
		}

		if isNext, err = middleware.Cookie(w, r); err != nil {
			return
		}
		if !isNext {
			return
		}
		if isNext, err = middleware.Handle(w, r); err != nil {
			return
		}
		if !isNext {
			return
		}
	})
	middlewareRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var middleware = middleware(context.TODO(), mux.CurrentRoute(r), customizeViewRender, m.serverName)

			var isNext bool
			var err error
			defer func() {
				middleware.Defer(w, r, err)
			}()

			if isNext, err = middleware.Path(w, r); err != nil {
				return
			}
			if !isNext {
				return
			}
			if isNext, err = middleware.Cookie(w, r); err != nil {
				return
			}
			if !isNext {
				return
			}
			if isNext, err = middleware.Handle(w, r); err != nil {
				return
			}
			if !isNext {
				return
			}
			next.ServeHTTP(w, r)
		})
	})

}
func (m *httpServer) Listen() {
	log.Printf("HttpServer Listen:%s", m.listenAddr)
	srv := &http.Server{
		Handler:      m.engine,
		Addr:         m.listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}

func NewHttpServer(engine *mux.Router, serverName string, listenAddr string) *httpServer {

	return &httpServer{listenAddr: listenAddr, engine: engine, serverName: serverName}
}
