package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/iservice"
	"github.com/nbvghost/dandelion/service/redis"
	"github.com/nbvghost/dandelion/service/serviceobject"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool/encryption"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type service struct {
	port       int
	grpcClient iservice.IGrpcClient
	redis      iservice.IRedis
}

func (m *service) Listen() {
	engine := gin.Default()
	engine.Use(func(ginContext *gin.Context) {
		paths := strings.Split(ginContext.Request.URL.Path, "/")
		var path string
		if len(paths) >= 3 {
			appName := paths[1]
			endpointName := paths[2]

			if endpointName == "api" {
				//api请求
				path = "/" + strings.Join(paths[3:], "/")
			} else {
				//api请求并返回相关的页面
				path = "/" + strings.Join(paths[2:], "/")
			}

			ginContext.Set("AppName", appName)
			ginContext.Set("IsApi", endpointName == "api")
			ginContext.Set("Path", path)

		} else {

			ginContext.Abort()
		}

	})
	engine.Use(func(ginContext *gin.Context) {
		var err error
		var isApi bool

		defer func() {
			if err != nil {
				if isApi {
					ginContext.JSON(http.StatusOK, result.NewError(err))
				} else {
					ginContext.Redirect(http.StatusNotFound, "/error/404")
				}
				ginContext.Abort()
			}
		}()

		cookie, err := ginContext.Request.Cookie("token")
		var token string
		if err != nil || strings.EqualFold(cookie.Value, "") {
			token = encryption.CipherEncrypter(encryption.NewSecretKey(conf.Config.SecureKey), fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")))
			http.SetCookie(ginContext.Writer, &http.Cookie{Name: "token", Value: token, Path: "/", MaxAge: conf.Config.SessionExpires})

		} else {
			token = cookie.Value
		}

		if v, ok := ginContext.Get("IsApi"); ok {
			isApi = v.(bool)
		}

		sessionText, _ := m.redis.Get(ginContext, redis.NewTokenKey(token))
		if sessionText != "" {
			session := &gweb.Session{}
			if err = json.Unmarshal([]byte(sessionText), session); err != nil {
				return
			}
			ginContext.Set("ID", session.ID)
		}

	})
	engine.Use(func(ginContext *gin.Context) {
		var isApi bool
		var appName string
		var path string
		var UID string

		if v, ok := ginContext.Get("IsApi"); ok {
			isApi = v.(bool)
		}
		if v, ok := ginContext.Get("AppName"); ok {
			appName = v.(string)
		}
		if v, ok := ginContext.Get("Path"); ok {
			path = v.(string)
		}

		if v, ok := ginContext.Get("UID"); ok {
			UID = v.(string)
		}

		var err error

		defer func() {
			if err != nil {
				if isApi {
					ginContext.JSON(http.StatusOK, result.NewError(err))
				} else {
					ginContext.Redirect(http.StatusNotFound, "/error/404")
				}
			}
		}()

		var bodyBytes []byte
		bodyBytes, err = ginContext.GetRawData()
		if err != nil {
			return
		}

		Timeout, _ := strconv.ParseUint(ginContext.Request.Header.Get("Timeout"), 10, 64)

		ctx := context.Background()
		response := &serviceobject.GrpcResponse{}
		response, err = m.grpcClient.Call(ctx, appName, &serviceobject.GrpcRequest{
			Route:      path,
			HttpMethod: ginContext.Request.Method,
			Timeout:    uint64(Timeout),
			Header:     nil,
			Body:       bodyBytes,
			Uri:        ginContext.Request.URL.RawQuery,
			UID:        UID,
		})
		if err != nil {
			return
		}

		if isApi {
			ginContext.JSON(http.StatusOK, response)
		} else {

		}
		if err != nil {
			//panic(err)
			return
		}

		log.Println(path)

		ginContext.Abort()

	})
	err := engine.Run(fmt.Sprintf(":%d", m.port))
	if err != nil {
		log.Fatalln(err)
	}
}
func New(port int, grpcClient iservice.IGrpcClient, redis iservice.IRedis) *service {
	return &service{port: port, grpcClient: grpcClient, redis: redis}
}
