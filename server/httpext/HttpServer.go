package httpext

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/nbvghost/dandelion/library/result"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/server/route"
	"github.com/pkg/errors"

	"github.com/nbvghost/dandelion/library/contexext"
)

const defaultMemory = 32 << 20

type Session struct {
	ID    string
	Token string
}

type httpServer struct {
	//serverDesc          *serverDesc
	engine             *mux.Router
	route              constrain.IRoute
	redisClient        constrain.IRedis
	etcdClient         constrain.IEtcd
	errorHandleResult  constrain.IResultError
	router             *mux.Router
	viewRender         constrain.IViewRender
	notFoundViewRender constrain.IViewRender
	middlewares        []constrain.IMiddleware
	defaultMiddleware  *httpMiddleware
}

func (m *httpServer) ApiErrorHandle(result constrain.IResultError) {
	m.errorHandleResult = result
}

func (m *httpServer) Use(middleware constrain.IMiddleware) {
	m.middlewares = append(m.middlewares, middleware)
}
func (m *httpServer) Listen(microServerConfig *config.MicroServerConfig) error {
	if m.etcdClient != nil {
		if microServerConfig.Port == 0 {
			var err error
			microServerConfig, err = m.etcdClient.Register(microServerConfig)
			if err != nil {
				return err
			}
		}
	}

	listenAddr := fmt.Sprintf("%s:%d", microServerConfig.IP, microServerConfig.Port)
	log.Printf("HttpServer Listen:%s", listenAddr)
	srv := &http.Server{
		Handler:      m.engine,
		Addr:         listenAddr,
		WriteTimeout: 3 * time.Minute,
		ReadTimeout:  3 * time.Minute,
	}
	return srv.ListenAndServe()
}
func (m *httpServer) handleError(ctx constrain.IContext, customizeViewRender constrain.IViewRender, w http.ResponseWriter, r *http.Request, err error) {
	var bytes []byte
	contextValue := contexext.FromContext(ctx)

	if err != nil {
		if contextValue.IsApi {
			if m.errorHandleResult != nil {
				m.errorHandleResult.Apply(ctx, err)
				return
			}
			if ar, ok := err.(*result.ActionResult); ok {
				w.Header().Set("Code", fmt.Sprintf("%d", ar.Code))
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			var e error
			bytes, e = json.Marshal(result.NewError(err))
			if e != nil {
				log.Println(e)
			}
			w.Write(bytes)
		} else {
			w.WriteHeader(http.StatusNotFound)
			d := map[string]interface{}{
				"ErrorText": err.Error(),
				"Stack":     fmt.Sprintf("%+v", errors.WithStack(err)),
			}

			if ar, ok := err.(*result.ActionResult); ok {
				w.Header().Set("Code", fmt.Sprintf("%d", ar.Code))
			}

			viewResult := route.NewViewResult("404", d)
			viewBaseValue := reflect.ValueOf(viewResult).Elem().FieldByName("ViewBase")
			viewBase := viewBaseValue.Interface().(extends.ViewBase)

			htmlMeta := extends.NewHtmlMeta(contextValue.Lang, contextValue.RequestUrl)
			if viewBase.HtmlMetaCallback != nil {
				if err = viewBase.HtmlMetaCallback(viewBase, htmlMeta); err != nil {
					ctx.Logger().Error("render", zap.Error(err))
				}
			}
			viewBase.HtmlMeta = htmlMeta
			viewBaseValue.Set(reflect.ValueOf(viewBase))

			if customizeViewRender == nil {
				ctx.Logger().Error("render", zap.Error(errors.New("没找开视图渲染器")))
				return
			}

			if err = customizeViewRender.Render(ctx, r, w, viewResult); err != nil {
				ctx.Logger().Error("render", zap.Error(err))
			}
			return

			/*vr := &viewRender{}
			if err = vr.Render(ctx, r, w, viewResult); err != nil {
				ctx.Logger().Error("render", zap.Error(err))
			}*/

			/*t, errTemplate := template.New("").Parse(html_404)
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
				w.Write([]byte(errTemplate.Error()))
			}*/
		}
	}
}

type Option interface {
	apply(server *httpServer)
}

/*type serverDesc struct {
	serverName string
	ip         string
	port       int
}

func (s *serverDesc) apply(server *httpServer) {
	server.serverDesc = s
}*/

/*func WithServerDesc(serverName string, ip string, port int) Option {
	return &serverDesc{serverName: serverName, ip: ip, port: port}
}*/

type emptyOption struct {
	applyFunc func(server *httpServer)
}

func newOption(apply func(server *httpServer)) *emptyOption {
	return &emptyOption{applyFunc: apply}
}
func (e *emptyOption) apply(server *httpServer) {
	e.applyFunc(server)
}

/*
	func WithEtcdOption(etcdClient constrain.IEtcd) Option {
		return newOption(func(server *httpServer) {
			server.etcdClient = etcdClient
		})
	}

	func WithRedisOption(redisClient constrain.IRedis) Option {
		return newOption(func(server *httpServer) {
			server.redisClient = redisClient
		})
	}
*/
func WithViewRenderOption(customizeViewRender constrain.IViewRender) Option {
	return newOption(func(server *httpServer) {
		server.viewRender = customizeViewRender
	})
}
func WithNotFoundViewRenderOption(customizeViewRender constrain.IViewRender) Option {
	return newOption(func(server *httpServer) {
		server.notFoundViewRender = customizeViewRender
	})
}

func (m *httpServer) handlerFunc(viewRender constrain.IViewRender, response http.ResponseWriter, request *http.Request) (next bool, ctxValue *contexext.ContextValue) {

	var ctx constrain.IContext
	ctxValue = contexext.FromContext(request.Context())
	if ctxValue != nil {
		ctx = request.Context().(constrain.IContext)
		ctxValue = contexext.FromContext(ctx)
	} else {
		ctx = m.defaultMiddleware.CreateContext(m.redisClient, m.etcdClient, m.route, response, request) //CreateContext(m.redisClient, m.etcdClient, m.route, response, request)
		ctxValue = contexext.FromContext(ctx)
		ctxValue.Request = request.WithContext(ctx)
	}

	defer func() {
		var err error
		if rerr := recover(); rerr != nil {
			switch rerr.(type) {
			case error:
				err = rerr.(error)
			default:
				err = fmt.Errorf("%v", rerr)
			}
			ctx.Logger().Error("http-server", zap.Error(err))
			m.handleError(ctx, viewRender, ctxValue.Response, ctxValue.Request, err)
		}

	}()

	for i := range m.middlewares {
		middleware := m.middlewares[i]
		if err := middleware.Handle(ctx, m.route, viewRender, ctxValue.Response, ctxValue.Request); err != nil {
			ctx.Logger().Error("http-server", zap.Error(err))
			m.handleError(ctx, viewRender, ctxValue.Response, ctxValue.Request, err)
			return
		}
	}
	if ctxValue.Request.Method == http.MethodOptions {
		return
	}
	if err := m.defaultMiddleware.Handle(ctx, m.route, viewRender, ctxValue.Response, ctxValue.Request); err != nil {
		ctx.Logger().Error("http-server", zap.Error(err))
		m.handleError(ctx, viewRender, ctxValue.Response, ctxValue.Request, err)
		return
	}
	return true, ctxValue
}

func NewHttpServer(etcdClient constrain.IEtcd, redisClient constrain.IRedis, engine *mux.Router, router *mux.Router, mRoute constrain.IRoute, ops ...Option) *httpServer {
	s := &httpServer{router: router, redisClient: redisClient, etcdClient: etcdClient, route: mRoute, engine: engine, defaultMiddleware: &httpMiddleware{}}
	for i := range ops {
		ops[i].apply(s)
	}

	if s.router != nil && s.route != nil {
		router.NotFoundHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if s.notFoundViewRender == nil {
				if request.URL.Path == "/404" {
					writer.WriteHeader(http.StatusNotFound)
					t, _ := template.New("404").Parse("404")
					err := t.Execute(writer, nil)
					if err != nil {
						log.Println(err)
					}
				} else {
					http.Redirect(writer, request, "/404", http.StatusPermanentRedirect)
				}
			} else {
				s.handlerFunc(s.notFoundViewRender, writer, request)
			}

		})
		//router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/*ctx := DefaultHttpMiddleware.CreateContext(s.redisClient, mRoute, w, r)
		w.WriteHeader(http.StatusNotFound)
		viewResult := route.NewViewResult("404", map[string]any{})
		viewBaseValue := reflect.ValueOf(viewResult).Elem().FieldByName("ViewBase")
		viewBase := viewBaseValue.Interface().(extends.ViewBase)
		viewBaseValue.Set(reflect.ValueOf(viewBase))
		if s.customizeViewRender != nil {
			if err := s.customizeViewRender.Render(ctx, r, w, viewResult); err != nil {
				ctx.Logger().Error("render", zap.Error(err))
			}
			return
		}
		vr := &viewRender{}
		if err := vr.Render(ctx, r, w, viewResult); err != nil {
			ctx.Logger().Error("render", zap.Error(err))
		}*/
		//})

		s.router.Use(func(next http.Handler) http.Handler {
			if s.route == nil {
				log.Println("没有启用路由功能，因为httpServer.route(constrain.IRoute)对象为空")
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			}
			return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				hasNext, ctxValue := s.handlerFunc(s.viewRender, response, request)
				if hasNext {
					next.ServeHTTP(ctxValue.Response, ctxValue.Request)
				}
			})
		})
	}

	return s
}
