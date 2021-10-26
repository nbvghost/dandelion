package http

import (
	"context"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/nbvghost/dandelion/service/iservice"
	"github.com/nbvghost/dandelion/service/serviceobject"
)

type service struct {
	port       int
	grpcClient iservice.IGrpcClient
}

func (m *service) Listen() {
	engine := gin.Default()
	engine.Use(func(ginContext *gin.Context) {
		paths := strings.Split(ginContext.Request.URL.Path, "/")
		log.Println(paths, len(paths))

		var err error
		if len(paths) >= 2 {
			appName := paths[0]

			var path string

			args := make(map[string]interface{})

			if err = ginContext.ShouldBind(&args); err != nil {
				return
			}

			if appName == "api" {
				path = "/" + strings.Join(paths[1:], "/")

			} else {
				path = ginContext.Request.URL.Path
			}

			ctx := context.Background()
			response := &serviceobject.GrpcResponse{}
			err = m.grpcClient.Call(ctx, appName, serviceobject.GrpcRequest{
				Route:      path,
				HttpMethod: ginContext.Request.Method,
				Body:       args,
			}, response)

			if err != nil {
				return
			}
			ginContext.HTML()
			template.ParseGlob()

		}

		ginContext.Abort()

	})
	err := engine.Run(fmt.Sprintf(":%d", m.port))
	if err != nil {
		log.Fatalln(err)
	}
}
func New(port int, grpcClient iservice.IGrpcClient) *service {
	return &service{port: port, grpcClient: grpcClient}
}
