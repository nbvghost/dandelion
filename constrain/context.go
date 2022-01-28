package constrain

import (
	"context"
	"github.com/nbvghost/dandelion/service/redis"
	"github.com/nbvghost/gpa/types"
)

//用于注入的回调方法
type ICallback interface {
	Before(context IContext, handler interface{}) error
	ViewAfter(context IContext, r IViewResult) error
}
type IContext interface {
	context.Context
	Redis() redis.IRedis
	UID() types.PrimaryKey
	//Query() url.Values
	//AppName() string
	Route() string
	//SelectFileServer() string
	//SelectServer(appName string) (string, error)
}
