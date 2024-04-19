package constrain

import (
	"context"
	"go.uber.org/zap"
	"sync"

	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/library/dao"
)

// 用于注入的回调方法
type IMappingCallback interface {
	Mapping(context IContext, handler interface{}) error
	// Deprecated: 好像没有用
	//ViewAfter(context IContext, r IViewResult) error
	AddMapping(mapping IMapping) IMappingCallback
}

type IContext interface {
	context.Context
	Redis() IRedis
	Etcd() IEtcd
	UID() dao.PrimaryKey
	AppName() string
	Route() string
	Token() string
	Logger() *zap.Logger
	Mode() key.Mode
	SyncCache() *sync.Map
	Destroy()
	Mapping(v interface{})
	//DomainName() string
	//SelectFileServer() string
}
