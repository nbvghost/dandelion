package context

import (
	"context"
	"github.com/nbvghost/dandelion/service/redis"
	"github.com/nbvghost/gpa/types"
)

type IContext interface {
	Redis() redis.IRedis
	Context() context.Context
	UID() types.PrimaryKey
	Query() interface{}
	AppName() string
	Route() string
	SelectFileServer() string
}
