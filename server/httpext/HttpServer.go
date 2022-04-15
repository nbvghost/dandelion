package httpext

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/library/environments"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/server/route"
	"github.com/nbvghost/tool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"reflect"
	"strconv"

	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/library/contexext"

	"github.com/nbvghost/dandelion/server/redis"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool/encryption"
)

const defaultMemory = 32 << 20

type httpServer struct {
	serverName  string
	listenAddr  string
	engine      *mux.Router
	route       constrain.IRoute
	redisClient constrain.IRedis
}

type Session struct {
	ID    string
	Token string
}

func (m *httpServer) getToken(w http.ResponseWriter, r *http.Request) string {
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

	return token
}

func (m *httpServer) getSession(parentCtx context.Context, token string) (Session, error) {
	var se Session
	se.Token = token
	var sessionText string
	sessionText, _ = m.redisClient.GetEx(parentCtx, redis.NewTokenKey(token), time.Minute*10)
	if sessionText != "" {
		if err := json.Unmarshal([]byte(sessionText), &se); err != nil {
			return se, err
		}
	}
	return se, nil
}
func (m *httpServer) mustContext(route constrain.IRoute, w http.ResponseWriter, r *http.Request) (ctx constrain.IContext) {
	var lang string
	domainPrefix, domainName := util.ParseDomain(r.Host)
	if len(domainPrefix) >= 1 {
		lang = domainPrefix[0]
	}
	if len(lang) == 0 || strings.EqualFold(lang, "dev") {
		lang = "en"
	}

	parentCtx := context.TODO()

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
	session, err = m.getSession(parentCtx, token)
	if err != nil {
		logger.Error("getSession", zap.Error(err))
	}
	contextValue := &contexext.ContextValue{
		Mapping:    route.GetMappingCallback(),
		Response:   w,
		Request:    r,
		Timeout:    Timeout,
		DomainName: domainName,
		Lang:       lang,
		RequestUrl: util.GetFullUrl(r),
		//PathTemplate: pathTemplate,
		Query: r.URL.Query(),
	}
	ctx = contexext.New(contexext.NewContext(parentCtx, contextValue), m.serverName, session.ID, r.URL.Path, m.redisClient, session.Token, logger, key.Mode(mode))
	return ctx
}

func (m *httpServer) Use(middlewareRouter *mux.Router, customizeViewRender constrain.IViewRender, middleware constrain.IMiddleware) {
	middlewareRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var isNext bool
		var err error
		var pathinfo pathInfo

		ctx := m.mustContext(m.route, w, r)
		defer func() {
			m.handleError(ctx, customizeViewRender, pathinfo.IsApi, w, r, err)
		}()

		pathinfo, err = getPathInfo(r)
		if err != nil {
			return
		}

		if isNext, err = middleware.Handle(ctx, m.route, r.URL.Path, pathinfo.IsApi, customizeViewRender, w, r); err != nil {
			return
		}
		if !isNext {
			return
		}
	})
	middlewareRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var isNext bool
			var err error
			var pathinfo pathInfo

			ctx := m.mustContext(m.route, w, r)
			defer func() {
				m.handleError(ctx, customizeViewRender, pathinfo.IsApi, w, r, err)
			}()

			pathinfo, err = getPathInfo(r)
			if err != nil {
				return
			}

			if strings.EqualFold(ctx.Mode().String(), key.ModeRelease.String()) && !environments.Release() {
				err = errors.New("正式环境访问开发环境的服务")
				return
			}
			if strings.EqualFold(ctx.Mode().String(), key.ModeDev.String()) && environments.Release() {
				err = errors.New("开发环境访问正式环境的服务")
				return
			}

			if isNext, err = middleware.Handle(ctx, m.route, pathinfo.PathTemplate, pathinfo.IsApi, customizeViewRender, w, r); err != nil {
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
func (m *httpServer) handleError(ctx constrain.IContext, customizeViewRender constrain.IViewRender, isApi bool, w http.ResponseWriter, r *http.Request, err error) {
	var bytes []byte
	contextValue := contexext.FromContext(ctx)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if isApi {
			//ginContext.JSON(http.StatusOK, action.NewError(err))
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			bytes, err = json.Marshal(action.NewError(err))
			if err != nil {
				log.Println(err)
			}
			w.Write(bytes)
		} else {
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

func NewHttpServer(engine *mux.Router, route constrain.IRoute, redisClient constrain.IRedis, serverName string, listenAddr string) *httpServer {

	return &httpServer{listenAddr: listenAddr, route: route, engine: engine, redisClient: redisClient, serverName: serverName}
}

type pathInfo struct {
	PathTemplate string
	Path         string
	IsApi        bool
}

func getPathInfo(r *http.Request) (pathInfo pathInfo, err error) {
	currentRoute := mux.CurrentRoute(r)

	var pathTemplate string
	pathTemplate, err = currentRoute.GetPathTemplate()
	if err != nil {
		return
	}

	var apiPath = "/api/"
	if len(pathTemplate) >= len(apiPath) && strings.EqualFold(pathTemplate[0:len(apiPath)], apiPath) {
		pathInfo.IsApi = true
	}

	pathInfo.PathTemplate = pathTemplate
	pathInfo.Path = r.URL.Path
	return
}
