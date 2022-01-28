package httpext

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/route"
	"html/template"

	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/library/contexext"

	"github.com/nbvghost/dandelion/library/funcmap"
	"github.com/nbvghost/dandelion/library/gobext"
	"github.com/nbvghost/dandelion/service/redis"
	"github.com/nbvghost/dandelion/service/serviceobject"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool/encryption"
)

const defaultMemory = 32 << 20

type httpServer struct {
	serverName string
	port       int
	engine     *mux.Router
}
type HttpMiddleware struct {
	currentRoute *mux.Route
	context      context.Context
	redis        redis.IRedis
	serverName   string

	response     *serviceobject.GrpcResponse
	session      gweb.Session
	isApi        bool
	pathTemplate string
	token        string
	route        route.IRoute
}

func (m *HttpMiddleware) Path(w http.ResponseWriter, r *http.Request) (bool, error) {
	//log.Println(writer, request)
	//log.Println(re.GetHostTemplate())
	//log.Println(re.GetPathTemplate())
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
		http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", MaxAge: conf.Config.SessionExpires})
	} else {
		token = cookie.Value
	}
	m.token = token
	sessionText, _ := m.redis.GetEx(m.context, redis.NewTokenKey(token), time.Minute*10)

	if sessionText != "" {
		session := gweb.Session{}
		if err = json.Unmarshal([]byte(sessionText), &session); err != nil {
			return false, err
		}
		//ginContext.Set("Session", session)
		m.session = session
	}
	return true, nil
}

func (m *HttpMiddleware) Handle(w http.ResponseWriter, r *http.Request) (bool, error) {

	var err error

	ctx := contexext.New(context.TODO(), m.serverName, m.session.ID, m.pathTemplate, r.URL.Query(), m.redis)
	ctx.Attributes().Put("Token", m.token)

	Timeout, _ := strconv.ParseUint(r.Header.Get("Timeout"), 10, 64)

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/x-www-form-urlencoded":
		if err = r.ParseForm(); err != nil {
			return false, err
		}
	case "application/json":
	case "":
	default:
		err = fmt.Errorf("不支持提交内容：%s", contentType)
		return false, err
	}

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return false, err
	}

	grpcRequest := &serviceobject.GrpcRequest{
		AppName:    m.serverName,
		Route:      m.pathTemplate,
		HttpMethod: r.Method,
		Timeout:    Timeout,
		Header:     url.Values(r.Header).Encode(),
		Query:      r.URL.Query().Encode(),
		Uri:        mux.Vars(r),
		Form:       r.Form.Encode(),
		Body:       bodyBytes,
		UID:        m.session.ID,
		IsApi:      m.isApi,
	}

	var routeInfo constrain.IRouteInfo
	routeInfo, err = m.route.GetInfo(grpcRequest)
	if err != nil {
		return false, err
	}

	var broken bool
	broken, err = m.route.ExecuteInterceptor(ctx, routeInfo, w, r)
	if err != nil {
		return false, err
	}
	if broken {
		return broken, nil
	}

	m.response, err = m.route.Handle(ctx, routeInfo, grpcRequest)
	if err != nil {
		if v, ok := err.(*action.ActionResult); ok {
			err = fmt.Errorf(v.Message)
		} else {
			err = fmt.Errorf(err.Error())
		}
		return false, err
	}

	if m.isApi {
		var dataBytes []byte
		var head *result.Head
		if dataBytes, head, err = result.UnmarshalResult(m.response.GetData()); err != nil {
			return false, err
		}
		switch head.Mine {
		case result.MIME_APPLICATION_JSON:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(dataBytes)
		case result.MIME_TEXT_PLAIN:
			//ginContext.String(http.StatusOK, string(dataBytes))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(dataBytes)
		default:
			err = fmt.Errorf("不支持数据格式:%d", head.Mine)
			return false, err

		}

	} else {
		v := gobext.NewGob(m.response.Name)
		if err = gob.NewDecoder(bytes.NewReader(m.response.Data)).Decode(v); err != nil {
			return false, err
		}
		if err = m.render(ctx, w, v); err != nil {
			return false, err
		}
		//ginContext.Data(http.StatusOK, "text/html; charset=utf-8", bytes)

	}
	return true, nil
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
					"Stack":     string(debug.Stack()),
				}
				if m.response != nil {
					d["Data"] = string(m.response.Data)
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
func (m *HttpMiddleware) render(context constrain.IContext, writer http.ResponseWriter, viewData interface{}) error {
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

type MiddlewareFunc func(context context.Context, currentRoute *mux.Route, serverName string) constrain.IMiddleware

func NewHttpMiddleware(context context.Context, serverName string, currentRoute *mux.Route, route route.IRoute, redis redis.IRedis) constrain.IMiddleware {
	return &HttpMiddleware{
		currentRoute: currentRoute,
		redis:        redis,
		context:      context,
		serverName:   serverName,
		route:        route,
	}
}
func (m *httpServer) Use(middlewareRouter *mux.Router, middleware MiddlewareFunc) {
	middlewareRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var middleware = middleware(context.TODO(), mux.CurrentRoute(r), m.serverName)

			var isNext bool
			var err error
			defer middleware.Defer(w, r, err)

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
	srv := &http.Server{
		Handler:      m.engine,
		Addr:         fmt.Sprintf(":%d", m.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}

func NewHttpServer(engine *mux.Router, serverName string, port int) *httpServer {
	return &httpServer{port: port, engine: engine, serverName: serverName}
}
