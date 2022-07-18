package constrain

import (
	"context"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/gpa/types"
	"go.uber.org/zap"
)

//用于注入的回调方法
type IMappingCallback interface {
	Before(context IContext, handler interface{})
	// Deprecated: 好像没有用
	ViewAfter(context IContext, r IViewResult) error
}

type IContext interface {
	context.Context
	Redis() IRedis
	UID() types.PrimaryKey
	AppName() string
	SelectInsideServer(appName key.MicroServer) (string, error)
	GetDNSName(localName key.MicroServer) (string, error)
	Route() string
	Token() string
	Logger() *zap.Logger
	Mode() key.Mode
	//DomainName() string
	//SelectFileServer() string
}
