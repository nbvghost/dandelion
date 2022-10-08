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
	TryLock(ctx context.Context, key string, timeout ...time.Duration) (bool, func())
	GetEtcd() IEtcd
	GenerateUID(ctx context.Context) uint64
	Expire(ctx context.Context, key string, expiration time.Duration) error
	HSet(ctx context.Context, key string, value map[string]any) error
	HMGet(ctx context.Context, key string, fields ...string) ([]any, error)
	HGet(ctx context.Context, key, field string) (string, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
}
