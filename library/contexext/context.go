package contexext

import (
	"context"
	"errors"

	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
)

type HandlerContext struct {
	uid       dao.PrimaryKey
	parent    context.Context
	redis     constrain.IRedis
	etcd      constrain.IEtcd
	mode      key.Mode
	logger    *zap.Logger
	appName   string
	route     string
	token     string
	syncCache *sync.Map
	mapping   constrain.IMappingCallback
	DB        *gorm.DB
}

func (m *HandlerContext) GetDB() *gorm.DB {
	return m.DB
}

func (m *HandlerContext) Mapping(v interface{}) {
	if m.mapping == nil {
		log.Println()
		m.Logger().Info("mapping", zap.Error(errors.New("不支持 mapping 方法")))
		return
	}
	err := m.mapping.Mapping(m, v)
	if err != nil {
		m.Logger().With(zap.Error(err))
	}
}

type ContextKey struct{}
type ContextValue struct {
	Timeout    uint64
	Response   http.ResponseWriter
	Request    *http.Request
	DomainName string
	Lang       string
	RequestUrl string //
	//PathTemplate string //
	IsApi bool
	Query url.Values
}

func NewContext(parentCtx context.Context, v *ContextValue) context.Context {
	return context.WithValue(parentCtx, ContextKey{}, v)
}

func FromContext(ctx context.Context) *ContextValue {
	m := ctx.Value(ContextKey{})
	v, _ := m.(*ContextValue)
	return v
}
func (m *HandlerContext) Etcd() constrain.IEtcd {
	return m.etcd
}
func (m *HandlerContext) Deadline() (deadline time.Time, ok bool) {
	return m.parent.Deadline()
}
func (m *HandlerContext) Done() <-chan struct{} {
	return m.parent.Done()
}
func (m *HandlerContext) SyncCache() *sync.Map {
	return m.syncCache
}
func (m *HandlerContext) Err() error {
	return m.parent.Err()
}
func (m *HandlerContext) Value(key interface{}) interface{} {
	return m.parent.Value(key)
}
func (m *HandlerContext) Route() string {
	return m.route
}
func (m *HandlerContext) AppName() string {
	return m.appName
}
func (m *HandlerContext) UID() dao.PrimaryKey {
	return m.uid
}
func (m *HandlerContext) Context() context.Context {
	return m.parent
}
func (m *HandlerContext) Redis() constrain.IRedis {
	return m.redis
}
func (m *HandlerContext) Logger() *zap.Logger {
	return m.logger
}

/*func (m *handlerContext) SelectInsideServer(appName key.MicroServer) (string, error) {
	return m.etcd.SelectInsideServer(appName)
}
func (m *handlerContext) SelectOutsideServer(appName key.MicroServer) (string, error) {
	return m.etcd.SelectOutsideServer(appName)
}*/

func (m *HandlerContext) Token() string {
	return m.token
}
func (m *HandlerContext) Mode() key.Mode {
	return m.mode
}
func (m *HandlerContext) Destroy() {
	m.syncCache.Range(func(key, value any) bool {
		m.syncCache.Delete(key)
		return true
	})
}
func New(parent context.Context, appName, uid string, route string, mapping constrain.IMappingCallback, etcd constrain.IEtcd, redis constrain.IRedis, token string, logger *zap.Logger, mode key.Mode) constrain.IContext {
	x := &HandlerContext{parent: parent, uid: dao.NewFromString(uid), mapping: mapping, route: route, etcd: etcd, redis: redis, appName: appName, token: token, logger: logger, mode: mode, syncCache: &sync.Map{}}
	x.DB = db.GetDB(x)
	return x
}

type ServiceContext struct {
	parent    context.Context
	redis     constrain.IRedis
	etcd      constrain.IEtcd
	logger    *zap.Logger
	appName   string
	syncCache *sync.Map
	db        *gorm.DB
}

func (m *ServiceContext) GetDB() *gorm.DB {
	return m.db
}

func (m *ServiceContext) Err() error {
	return m.parent.Err()
}

func (m *ServiceContext) Value(key any) any {
	return m.parent.Value(key)
}

func (m *ServiceContext) Redis() constrain.IRedis { return m.redis }

func (m *ServiceContext) Etcd() constrain.IEtcd { return m.etcd }

func (m *ServiceContext) Logger() *zap.Logger { return m.logger }

func (m *ServiceContext) SyncCache() *sync.Map { return m.syncCache }

func (m *ServiceContext) Destroy() {
	m.syncCache.Range(func(key, value any) bool {
		m.syncCache.Delete(key)
		return true
	})
}
func (m *ServiceContext) Deadline() (deadline time.Time, ok bool) {
	return m.parent.Deadline()
}
func (m *ServiceContext) Done() <-chan struct{} {
	return m.parent.Done()
}

func NewServiceContext(parent context.Context, appName string, etcd constrain.IEtcd, redis constrain.IRedis, logger *zap.Logger) constrain.IServiceContext {
	x := &ServiceContext{parent: parent, etcd: etcd, redis: redis, appName: appName, logger: logger, syncCache: &sync.Map{}}
	x.db = db.GetDB(x)
	return x
}
