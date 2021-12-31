package redis

import (
	"context"
	"github.com/nbvghost/dandelion/service/etcd"
)

type IRedis interface {
	Get(ctx context.Context, key string) (string, error)
	GetEtcd() etcd.IEtcd
}
