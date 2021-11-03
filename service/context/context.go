package context

import (
	"context"
	"github.com/nbvghost/dandelion/service/iservice"
	"github.com/nbvghost/gpa/types"
)

type HandlerContext struct {
	uid    types.PrimaryKey
	parent context.Context
	redis  iservice.IRedis
}

func (m *HandlerContext) UID() types.PrimaryKey {
	return m.uid
}
func (m *HandlerContext) Context() context.Context {
	return m.parent
}
func (m *HandlerContext) Redis() iservice.IRedis {
	return m.redis
}
func New(parent context.Context, uid string, redis iservice.IRedis) iservice.IContext {
	return &HandlerContext{parent: parent, uid: types.NewFromString(uid), redis: redis}
}
