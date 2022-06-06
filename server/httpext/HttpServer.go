package httpext

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/action"
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
	serverDesc          *serverDesc
	engine              *mux.Router
	route               constrain.IRoute
	redisClient         constrain.IRedis
	errorHandleResult   constrain.IResultError
	router              *mux.Router
	customizeViewRender constrain.IViewRender
}

func (m *httpServer) ApiErrorHandle(result constrain.IResultError) {
	m.errorHandleResult = result
}

func (m *httpServer) Use(middleware constrain.IMiddleware) {
	m.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error

			var ctx constrain.IContext
			ctxValue := contexext.FromContext(r.Context())
			if ctxValue != nil {
				ctx = r.Context().(constrain.IContext)
			} else {
				ctx = middleware.CreateContent(m.redisClient, m.route, w, r)
				r = r.WithContext(ctx)
			}

			defer func() {
				m.handleError(ctx, m.customizeViewRender, w, r, err)
			}()

			var isNext bool
			if isNext, err = middleware.Handle(ctx, m.route, m.customizeViewRender, w, r); err != nil {
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
	listenAddr := fmt.Sprintf("%s:%d", m.serverDesc.ip, m.serverDesc.port)
	log.Printf("HttpServer Listen:%s", listenAddr)
	srv := &http.Server{
		Handler:      m.engine,
		Addr:         listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}
func (m *httpServer) handleError(ctx constrain.IContext, customizeViewRender constrain.IViewRender, w http.ResponseWriter, r *http.Request, err error) {
	var bytes []byte
	contextValue := contexext.FromContext(ctx)

	if err != nil {
		ctx.Logger().Error(err.Error())
		if contextValue.IsApi {
			if m.errorHandleResult != nil {
				m.errorHandleResult.Apply(ctx, err)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			var e error
			bytes, e = json.Marshal(action.NewError(err))
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

			if customizeViewRender != nil {
				if err = customizeViewRender.Render(ctx, r, w, viewResult); err != nil {
					ctx.Logger().Error("render", zap.Error(err))
				}
				return
			}
			vr := &viewRender{}
			if err = vr.Render(ctx, r, w, viewResult); err != nil {
				ctx.Logger().Error("render", zap.Error(err))
			}

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

type serverDesc struct {
	serverName string
	ip         string
	port       int
}

func (s *serverDesc) apply(server *httpServer) {
	server.serverDesc = s
}
func WithServerDesc(serverName string, ip string, port int) Option {
	return &serverDesc{serverName: serverName, ip: ip, port: port}
}

type emptyOption struct {
	applyFunc func(server *httpServer)
}

func newOption(apply func(server *httpServer)) *emptyOption {
	return &emptyOption{applyFunc: apply}
}
func (e *emptyOption) apply(server *httpServer) {
	e.applyFunc(server)
}

func WithRedisOption(redisClient constrain.IRedis) Option {
	return newOption(func(server *httpServer) {
		server.redisClient = redisClient
	})
}
func WithCustomizeViewRenderOption(customizeViewRender constrain.IViewRender) Option {
	return newOption(func(server *httpServer) {
		server.customizeViewRender = customizeViewRender
	})
}

func NewHttpServer(engine *mux.Router, router *mux.Router, route constrain.IRoute, ops ...Option) *httpServer {
	if route == nil {
		panic(errors.New("参数route不能为空"))
	}
	s := &httpServer{router: router, route: route, engine: engine}
	for i := range ops {
		ops[i].apply(s)
	}

	if router != nil {
		router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var isNext bool
			var err error

			ctx := DefaultHttpMiddleware.CreateContent(s.redisClient, route, w, r)

			defer func() {
				s.handleError(ctx, s.customizeViewRender, w, r, err)
			}()

			/*var pathTemplate string
			pathTemplate, err = getPathTemplate(r)
			if err != nil {
				return
			}
			ctxValue := contexext.FromContext(ctx)*/

			if isNext, err = DefaultHttpMiddleware.Handle(ctx, route, s.customizeViewRender, w, r); err != nil {
				return
			}
			if !isNext {
				return
			}
		})
	}
	return s
}

func getPathTemplate(r *http.Request) (pathTemplate string, err error) {
	currentRoute := mux.CurrentRoute(r)
	pathTemplate, err = currentRoute.GetPathTemplate()
	if err != nil {
		return
	}
	return
}
