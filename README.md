### dandelion 多用户商城
-- 我仔细看了看蒲公英，它并不是一个娇气的花
目前主要功能：

* 1、后台管理接口
* 2、活动数据接口
* 3、组织信息管理
* 4、配制管理
* 5、内容管理与发布
* 6、商品管理
* 7、日志管理
* 8、订单管理
* 9、用户管理



安装：

```sh
go get -u github.com/nbvghost/dandelion
```


## 例示代码：
### 启动程序：
```go
package main

import (
	"flag"
	"github.com/nbvghost/stone/internal/enter/interceptor"
	"github.com/nbvghost/stone/internal/mappingEntity"
	"github.com/nbvghost/tool/encryption"
	"log"
	"net/http"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/nbvghost/stone/internal/enter/api"
	"github.com/nbvghost/stone/internal/enter/view"

	"github.com/nbvghost/stone/internal/service"

	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/environments"
	"github.com/nbvghost/dandelion/library/mapping"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/server/etcd"
	"github.com/nbvghost/dandelion/server/httpext"
	"github.com/nbvghost/dandelion/server/redis"
	"github.com/nbvghost/dandelion/server/route"
	"github.com/nbvghost/dandelion/server/serviceobject"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	flag.Parse()

	serverConfig := config.MicroServerConfig{
		MicroServer: key.MicroServer{Name: "stone", ServerType: key.ServerTypeHttp},
		IP:          environments.IP(),
		Port:        environments.Port(),
	}

	etcdConfig := clientv3.Config{
		Endpoints:   environments.EtcdEndpoints(),
		DialTimeout: 30 * time.Second,
		Username:    environments.EtcdUsername(),
		//Password:    environments.EtcdPassword(),
	}	

	etcdService := etcd.NewServer(etcdConfig)

	err := service.Init(etcdService, "stone")
	if err != nil {
		log.Fatalln(err)
	}

	redisOption, err := etcdService.ObtainRedis()
	if err != nil {
		log.Fatalln(err)
	}

	redisClient := redis.NewClient(*redisOption, etcdService)

	if environments.Release() {
		//service.StartTask(redisClient)
	}
	//patch.StartPatch()

	var serverDesc *serviceobject.ServerDesc
	serverDesc, err = etcdService.Register(serviceobject.NewServerDesc(serverConfig.MicroServer, serverConfig.Port, serverConfig.IP))
	if err != nil {
		panic(err)
	}

	//mappings := mappingEntity.New(&mappingInstance.OrganizationMapping{}, &mappingInstance.AdminMapping{})
	mappings := mapping.New()
	mappings.AddMapping(&mappingEntity.AdminMapping{})
	engine := mux.NewRouter()
	engine.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("view/assets"))))

	subRouter := engine.PathPrefix("/").Subrouter()

	r := route.New(subRouter, mappings)
	r.RegisterInterceptors("/api/", []string{"/api/account/login", "/api/image/captcha", "/api/account/register"}, &interceptor.AdminInterceptor{})
	api.Register(r.RegisterRoute)
	view.Register(r)

	httpServer := httpext.NewHttpServer(engine, subRouter, r,
		httpext.WithServerDesc(serverDesc.MicroServer.Name, serverDesc.IP, serverDesc.Port),
		httpext.WithRedisOption(redisClient), httpext.WithCustomizeViewRenderOption(&CustomizeViewRender{}))
	httpServer.Use(&httpMiddleware{})
	httpServer.Use(httpext.DefaultHttpMiddleware)
	httpServer.Listen()
}

type CustomizeViewRender struct{}

func (m *CustomizeViewRender) Render(context constrain.IContext, request *http.Request, writer http.ResponseWriter, viewData constrain.IViewResult) error {

	http.FileServer(http.Dir("view")).ServeHTTP(writer, request)
	return nil
}

type httpMiddleware struct{}

func (m *httpMiddleware) CreateContext(redisClient constrain.IRedis, router constrain.IRoute, w http.ResponseWriter, r *http.Request) constrain.IContext {
	return httpext.DefaultHttpMiddleware.CreateContext(redisClient, router, w, r)
}

func (m *httpMiddleware) Handle(ctx constrain.IContext, router constrain.IRoute, customizeViewRender constrain.IViewRender, w http.ResponseWriter, r *http.Request) (bool, error) {
	contextValue := contexext.FromContext(ctx)
	contextValue.Response.Header().Set("Access-Control-Allow-Origin", strings.TrimRight(r.Referer(), "/"))
	contextValue.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type,Page-No,Page-Size,Order-Field,Order-Method") //todo
	contextValue.Response.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,OPTIONS,GET,PUT")                             //todo
	contextValue.Response.Header().Set("Access-Control-Allow-Credentials", "true")                                                //todo

	if contextValue.IsApi && contextValue.Request.Method == http.MethodOptions {
		returnResult := &result.EmptyResult{}
		returnResult.Apply(ctx)
		return false, nil
	}
	return true, nil
}

```

### 技术支持（微信联系）
<img src="https://raw.githubusercontent.com/nbvghost/dandelion/master/add-me.jpg" alt="微信联系" width="256" />

### 使用IDEA开发
<a href="https://www.jetbrains.com/?from=dandelion"><img src="https://raw.githubusercontent.com/nbvghost/dandelion/master/icon-goland.png" alt="使用IDEA开发" width="128" height="128" align="bottom" /></a>
