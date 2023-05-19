package constrain

import (
	"context"
	"go.uber.org/zap"
	"sync"

	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/gpa/types"
)

// 用于注入的回调方法
type IMappingCallback interface {
	Before(context IContext, handler interface{})
	// Deprecated: 好像没有用
	ViewAfter(context IContext, r IViewResult) error
	AddMapping(mapping IMapping) IMappingCallback
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
	SyncCache() *sync.Map
	Destroy()
	//DomainName() string
	//SelectFileServer() string
}
