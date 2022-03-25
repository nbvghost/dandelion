package constrain

import (
	"context"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/gpa/types"
	"net/url"
)

//用于注入的回调方法
type IMappingCallback interface {
	Before(context IContext, handler interface{}) error
	// Deprecated: 好像没有用
	ViewAfter(context IContext, r IViewResult) error
}

type IContext interface {
	context.Context
	Redis() IRedis
	UID() types.PrimaryKey
	Query() url.Values
	AppName() string
	SelectInsideServer(appName key.MicroServerKey) (string, error)
	SelectServer(appName key.MicroServerKey) (string, error)
	Route() string
	Token() string
	//DomainName() string
	//SelectFileServer() string
}
