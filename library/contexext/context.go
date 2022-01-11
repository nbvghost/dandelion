package contexext

import (
	"context"
	"github.com/nbvghost/dandelion/constrain"
	"net/url"

	"github.com/nbvghost/dandelion/service/redis"
	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/gweb"
)

type handlerContext struct {
	uid        types.PrimaryKey
	parent     context.Context
	redis      redis.IRedis
	appName    string
	route      string
	query      url.Values
	attributes *gweb.Attributes
}

func (m *handlerContext) Attributes() *gweb.Attributes {
	return m.attributes
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
func (m *handlerContext) Redis() redis.IRedis {
	return m.redis
}
func (m *handlerContext) SelectServer(appName string) (string, error) {
	return m.redis.GetEtcd().SelectServer(appName)
}
func (m *handlerContext) SelectFileServer() string {
	return m.redis.GetEtcd().SelectFileServer()
}

func New(parent context.Context, appName, uid string, route string, query url.Values, redis redis.IRedis) constrain.IContext {

	return &handlerContext{parent: parent, uid: types.NewFromString(uid), query: query, route: route, redis: redis, appName: appName, attributes: &gweb.Attributes{}}
}
