package constrain

import (
	"context"
	"time"
)

type IRedis interface {
	Del(ctx context.Context, keys ...string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	GetEx(ctx context.Context, key string, expiration time.Duration) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetEtcd() IEtcd
	GenerateUID(ctx context.Context) uint64
}