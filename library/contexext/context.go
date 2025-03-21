package contexext

import (
	"context"
	"errors"
	"github.com/nbvghost/dandelion/library/dao"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
)

type handlerContext struct {
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
}

func (m *handlerContext) Mapping(v interface{}) {
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
func (m *handlerContext) Etcd() constrain.IEtcd {
	return m.etcd
}
func (m *handlerContext) Deadline() (deadline time.Time, ok bool) {
	return m.parent.Deadline()
}
func (m *handlerContext) Done() <-chan struct{} {
	return m.parent.Done()
}
func (m *handlerContext) SyncCache() *sync.Map {
	return m.syncCache
}
func (m *handlerContext) Err() error {
	return m.parent.Err()
}
func (m *handlerContext) Value(key interface{}) interface{} {
	return m.parent.Value(key)
}
func (m *handlerContext) Route() string {
	return m.route
}
func (m *handlerContext) AppName() string {
	return m.appName
}
func (m *handlerContext) UID() dao.PrimaryKey {
	return m.uid
}
func (m *handlerContext) Context() context.Context {
	return m.parent
}
func (m *handlerContext) Redis() constrain.IRedis {
	return m.redis
}
func (m *handlerContext) Logger() *zap.Logger {
	return m.logger
}

/*func (m *handlerContext) SelectInsideServer(appName key.MicroServer) (string, error) {
	return m.etcd.SelectInsideServer(appName)
}
func (m *handlerContext) SelectOutsideServer(appName key.MicroServer) (string, error) {
	return m.etcd.SelectOutsideServer(appName)
}*/

func (m *handlerContext) Token() string {
	return m.token
}
func (m *handlerContext) Mode() key.Mode {
	return m.mode
}
func (m *handlerContext) Destroy() {
	m.syncCache.Range(func(key, value any) bool {
		m.syncCache.Delete(key)
		return true
	})
}
func New(parent context.Context, appName, uid string, route string, mapping constrain.IMappingCallback, etcd constrain.IEtcd, redis constrain.IRedis, token string, logger *zap.Logger, mode key.Mode) constrain.IContext {
	return &handlerContext{parent: parent, uid: dao.NewFromString(uid), mapping: mapping, route: route, etcd: etcd, redis: redis, appName: appName, token: token, logger: logger, mode: mode, syncCache: &sync.Map{}}
}

type serviceContext struct {
	parent    context.Context
	redis     constrain.IRedis
	etcd      constrain.IEtcd
	logger    *zap.Logger
	appName   string
	syncCache *sync.Map
}

func (m *serviceContext) Err() error {
	return m.parent.Err()
}

func (m *serviceContext) Value(key any) any {
	return m.parent.Value(key)
}

func (m *serviceContext) Redis() constrain.IRedis { return m.redis }

func (m *serviceContext) Etcd() constrain.IEtcd { return m.etcd }

func (m *serviceContext) Logger() *zap.Logger { return m.logger }

func (m *serviceContext) SyncCache() *sync.Map { return m.syncCache }

func (m *serviceContext) Destroy() {
	m.syncCache.Range(func(key, value any) bool {
		m.syncCache.Delete(key)
		return true
	})
}
func (m *serviceContext) Deadline() (deadline time.Time, ok bool) {
	return m.parent.Deadline()
}
func (m *serviceContext) Done() <-chan struct{} {
	return m.parent.Done()
}

func NewServiceContext(parent context.Context, appName string, etcd constrain.IEtcd, redis constrain.IRedis, logger *zap.Logger) constrain.IServiceContext {
	return &serviceContext{parent: parent, etcd: etcd, redis: redis, appName: appName, logger: logger, syncCache: &sync.Map{}}
}
