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
	serverName        string
	listenAddr        string
	engine            *mux.Router
	route             constrain.IRoute
	redisClient       constrain.IRedis
	errorHandleResult constrain.IResultError
	router            *mux.Router
}

func (m *httpServer) ApiErrorHandle(result constrain.IResultError) {
	m.errorHandleResult = result
}

func (m *httpServer) Use(middleware constrain.IMiddleware, customizeViewRender constrain.IViewRender) {
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
				m.handleError(ctx, customizeViewRender, w, r, err)
			}()

			var isNext bool
			if isNext, err = middleware.Handle(ctx, m.route, customizeViewRender, w, r); err != nil {
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

func NewHttpServer(engine *mux.Router, router *mux.Router, route constrain.IRoute, redisClient constrain.IRedis, serverName string, listenAddr string) *httpServer {
	s := &httpServer{listenAddr: listenAddr, router: router, route: route, engine: engine, redisClient: redisClient, serverName: serverName}
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var isNext bool
		var err error

		ctx := DefaultHttpMiddleware.CreateContent(redisClient, route, w, r)

		defer func() {
			s.handleError(ctx, nil, w, r, err)
		}()

		/*var pathTemplate string
		pathTemplate, err = getPathTemplate(r)
		if err != nil {
			return
		}
		ctxValue := contexext.FromContext(ctx)*/

		if isNext, err = DefaultHttpMiddleware.Handle(ctx, route, nil, w, r); err != nil {
			return
		}
		if !isNext {
			return
		}
	})
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
