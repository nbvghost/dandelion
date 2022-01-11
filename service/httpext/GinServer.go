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

	"github.com/gin-gonic/gin"
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

type ginServer struct {
	serverName string
	port       int
	redis      redis.IRedis
	engine     *mux.Router
	route      route.IRoute
}

func pathMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
func cookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
func handleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
func (m *ginServer) Listen() {
	engine := m.engine

	engine.Use(func(ginContext *gin.Context) {
		paths := strings.Split(ginContext.Request.URL.Path, "/")
		var path string
		if len(paths) >= 2 {
			endpointName := paths[1]
			if endpointName == "api" {
				//api请求

				path = "/" + strings.Join(paths[2:], "/")
				ginContext.Set("IsApi", true)
				ginContext.Set("Path", path)
				return
			}
		}

		ginContext.Set("IsApi", false)
		ginContext.Set("Path", ginContext.Request.URL.Path)

	})

	engine.Use(func(ginContext *gin.Context) {
		var err error
		cookie, err := ginContext.Request.Cookie("token")
		var token string
		if err != nil || strings.EqualFold(cookie.Value, "") {
			token = encryption.CipherEncrypter(encryption.NewSecretKey(conf.Config.SecureKey), fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")))
			http.SetCookie(ginContext.Writer, &http.Cookie{Name: "token", Value: token, Path: "/", MaxAge: conf.Config.SessionExpires})
		} else {
			token = cookie.Value
		}
		ginContext.Set("Token", token)
		sessionText, _ := m.redis.GetEx(ginContext, redis.NewTokenKey(token), time.Minute*10)

		if sessionText != "" {
			session := &gweb.Session{}
			if err = json.Unmarshal([]byte(sessionText), session); err != nil {
				return
			}
			ginContext.Set("Session", session)
		}

	})

	engine.Use(func(ginContext *gin.Context) {
		var isApi bool
		var path string
		var session = &gweb.Session{}

		if v, ok := ginContext.Get("IsApi"); ok {
			isApi = v.(bool)
		}
		if v, ok := ginContext.Get("Path"); ok {
			path = v.(string)
		}

		if v, ok := ginContext.Get("Session"); ok {
			session = v.(*gweb.Session)
		}

		var err error
		var response *serviceobject.GrpcResponse

		defer func() {
			if err != nil {
				if isApi {
					ginContext.JSON(http.StatusOK, action.NewError(err))
				} else {
					//ginContext.Redirect(http.StatusFound, "/error/404")
					t, errTemplate := template.New("").Parse(html_404)
					buffer := bytes.NewBuffer(nil)
					if errTemplate == nil {
						d := map[string]interface{}{
							"ErrorText": err.Error(),
							"Stack":     string(debug.Stack()),
						}
						if response != nil {
							d["Data"] = string(response.Data)
						}
						errTemplate = t.Execute(buffer, d)
					}
					if errTemplate != nil {
						ginContext.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(errTemplate.Error()))
					} else {
						ginContext.Data(http.StatusNotFound, "text/html; charset=utf-8", buffer.Bytes())
					}

				}
				ginContext.Abort()
			}
		}()

		ctx := contexext.New(context.TODO(), m.serverName, session.ID, path, ginContext.Request.URL.Query(), m.redis)

		Timeout, _ := strconv.ParseUint(ginContext.Request.Header.Get("Timeout"), 10, 64)

		contentType := ginContext.Request.Header.Get("Content-Type")
		switch contentType {
		case "application/x-www-form-urlencoded":
			if err = ginContext.Request.ParseForm(); err != nil {
				return
			}
		case "application/json":
		case "":
		default:
			err = fmt.Errorf("不支持提交内容：%s", contentType)
			return
		}

		uriParams := make([]*serviceobject.Param, 0)
		for i := range ginContext.Params {
			item := ginContext.Params[i]
			uriParams = append(uriParams, &serviceobject.Param{
				Key:   item.Key,
				Value: item.Value,
			})
		}

		var bodyBytes []byte
		bodyBytes, err = ginContext.GetRawData()
		if err != nil {
			return
		}

		grpcRequest := &serviceobject.GrpcRequest{
			AppName:    m.serverName,
			Route:      path,
			HttpMethod: ginContext.Request.Method,
			Timeout:    Timeout,
			Header:     url.Values(ginContext.Request.Header).Encode(),
			Query:      ginContext.Request.URL.Query().Encode(),
			Uri:        uriParams,
			Form:       ginContext.Request.Form.Encode(),
			Body:       bodyBytes,
			UID:        session.ID,
			IsApi:      isApi,
		}

		var routeInfo constrain.IRouteInfo
		routeInfo, err = m.route.GetInfo(grpcRequest)
		if err != nil {
			return
		}

		var broken bool
		broken, err = m.route.ExecuteInterceptor(ctx, routeInfo, ginContext)
		if err != nil {
			return
		}
		if broken {
			return
		}

		response, err = m.route.Handle(ctx, routeInfo, grpcRequest)
		if err != nil {
			if v, ok := err.(*action.ActionResult); ok {
				err = fmt.Errorf(v.Message)
			} else {
				err = fmt.Errorf(err.Error())
			}
			return
		}

		if isApi {
			var dataBytes []byte
			var head *result.Head
			if dataBytes, head, err = result.UnmarshalResult(response.GetData()); err != nil {
				return
			}
			switch head.Mine {
			case result.MIME_APPLICATION_JSON:
				jsonResultMap := make(map[string]interface{})
				if err = json.Unmarshal(dataBytes, &jsonResultMap); err != nil {
					return
				}
				jsonResultMap["Code"] = response.GetCode()
				ginContext.JSON(http.StatusOK, jsonResultMap)
			case result.MIME_TEXT_PLAIN:
				ginContext.String(http.StatusOK, string(dataBytes))
			default:
				err = fmt.Errorf("不支持数据格式:%d", head.Mine)
				return

			}

		} else {
			v := gobext.NewGob(response.Name)
			if err = gob.NewDecoder(bytes.NewReader(response.Data)).Decode(v); err != nil {
				return
			}
			var bytes []byte
			if bytes, err = m.Render(ctx, v); err != nil {
				return
			}
			ginContext.Data(http.StatusOK, "text/html; charset=utf-8", bytes)

		}

	})

	err := engine.Run(fmt.Sprintf(":%d", m.port))
	if err != nil {
		log.Fatalln(err)
	}
}
func (m *ginServer) Render(context constrain.IContext, viewData interface{}) ([]byte, error) {
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
		return nil, err
	}

	var t *template.Template
	t, err = template.New("").Funcs(funcmap.NewFuncMap(context)).Parse(string(fileByte))
	if err != nil {
		return nil, err
	}
	t, err = t.ParseGlob(fmt.Sprintf("view/template/*.gohtml"))
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	err = t.Execute(buffer, map[string]interface{}{
		"Query": context.Query(),
		"Data":  viewData,
	})
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
func NewSingleServer(engine *mux.Router, route route.IRoute, redis redis.IRedis, serverName string, port int) *ginServer {
	return &ginServer{port: port, engine: engine, redis: redis, route: route, serverName: serverName}
}
