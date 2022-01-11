package constrain

import (
	"context"
	"net/url"

	"github.com/nbvghost/dandelion/service/redis"
	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/gweb"
)

//用于注入的回调方法
type ICallback interface {
	Before(context IContext, handler interface{}) error
	ViewAfter(context IContext, r IViewResult) error
}
type IContext interface {
	Redis() redis.IRedis
	Context() context.Context
	UID() types.PrimaryKey
	Query() url.Values
	AppName() string
	Route() string
	SelectFileServer() string
	Attributes() *gweb.Attributes
	SelectServer(appName string) (string, error)
}
