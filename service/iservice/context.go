package iservice

import (
	"context"
	"github.com/nbvghost/gpa/types"
)

type IContext interface {
	Redis() IRedis
	Context() context.Context
	UID() types.PrimaryKey
}
