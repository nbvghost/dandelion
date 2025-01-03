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
	AddMapping(mapping IMapping) IMappingCallback
	//todo 好像没有用
	//ViewAfter(context IContext, r IViewResult) error
}

type IService interface {
}
type IServiceContext interface {
	context.Context
	Redis() IRedis
	Etcd() IEtcd
	Logger() *zap.Logger
	SyncCache() *sync.Map
	Destroy()
}
type IContext interface {
	IServiceContext
	UID() dao.PrimaryKey
	AppName() string
	Route() string
	Token() string
	Mode() key.Mode
	Mapping(v interface{})
	//DomainName() string
	//SelectFileServer() string
}
