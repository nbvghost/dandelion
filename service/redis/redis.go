package redis

import (
	"context"
	"github.com/nbvghost/dandelion/service/etcd"
	"time"
)

type IRedis interface {
	Get(ctx context.Context, key string) (string, error)
	GetEx(ctx context.Context, key string, expiration time.Duration) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetEtcd() etcd.IEtcd
}
