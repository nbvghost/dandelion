package contexext

import (
	"context"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"

	"github.com/nbvghost/gpa/types"
)

type handlerContext struct {
	uid     types.PrimaryKey
	parent  context.Context
	redis   constrain.IRedis
	appName string
	route   string
	query   url.Values
	token   string
	logger  *zap.Logger
	mode    key.Mode
}

type ContextKey struct {
}
type ContextValue struct {
	Mapping    constrain.IMappingCallback
	Timeout    uint64
	Response   http.ResponseWriter
	Request    *http.Request
	DomainName string
	IsApi      bool
	Lang       string
	RequestUrl string //

}

func NewContext(v *ContextValue) context.Context {
	return context.WithValue(context.TODO(), ContextKey{}, v)
}

func FromContext(context constrain.IContext) *ContextValue {
	m := context.Value(ContextKey{})
	v, _ := m.(*ContextValue)
	return v
}

func (m *handlerContext) Deadline() (deadline time.Time, ok bool) {
	return m.parent.Deadline()
}

func (m *handlerContext) Done() <-chan struct{} {
	return m.parent.Done()
}

func (m *handlerContext) Err() error {
	return m.parent.Err()
}

func (m *handlerContext) Value(key interface{}) interface{} {
	return m.parent.Value(key)
}

func (m *handlerContext) Query() url.Values {
	return m.query
}
func (m *handlerContext) Route() string {
	return m.route
}
func (m *handlerContext) AppName() string {
	return m.appName
}
func (m *handlerContext) UID() types.PrimaryKey {
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
func (m *handlerContext) SelectInsideServer(appName key.MicroServerKey) (string, error) {
	return m.redis.GetEtcd().SelectInsideServer(appName)
}
func (m *handlerContext) GetDNSName(localName key.MicroServerKey) (string, error) {
	return m.redis.GetEtcd().GetDNSName(localName)
}
func (m *handlerContext) Token() string {
	return m.token
}
func (m *handlerContext) Mode() key.Mode {
	return m.mode
}
func New(parent context.Context, appName, uid string, route string, query url.Values, redis constrain.IRedis, token string, logger *zap.Logger, mode key.Mode) constrain.IContext {
	return &handlerContext{parent: parent, uid: types.NewFromString(uid), query: query, route: route, redis: redis, appName: appName, token: token, logger: logger, mode: mode}
}
