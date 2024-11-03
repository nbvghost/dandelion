package constrain

import (
	"net/http"
)

type IMiddleware interface {
	//Path(w http.ResponseWriter, r *http.Request) (bool, error)
	//Cookie(w http.ResponseWriter, r *http.Request) (bool, error)
	//CreateContext(redisClient IRedis, etcdClient IEtcd, router IRoute, w http.ResponseWriter, r *http.Request) IContext
	Handle(ctx IContext, router IRoute, beforeViewRender IBeforeViewRender, afterViewRender IAfterViewRender, w http.ResponseWriter, r *http.Request) error
	//Defer(w http.ResponseWriter, r *http.Request, err error)
}

// 路由执行前，处理视图的接口
type IBeforeViewRender interface {
	View(context IContext, request *http.Request, writer http.ResponseWriter, next func()) error
	Api(context IContext, request *http.Request, writer http.ResponseWriter, next func()) error
}

// 路由执行后，处理视图的接口
type IAfterViewRender interface {
	Render(context IContext, request *http.Request, writer http.ResponseWriter, viewData IViewResult) error
}
