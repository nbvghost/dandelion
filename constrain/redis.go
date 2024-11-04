package constrain

import (
	"context"
	"time"
)

type IRedis interface {
	Keys(ctx context.Context, key string) []string
	Del(ctx context.Context, keys ...string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	GetEx(ctx context.Context, key string, expiration time.Duration) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	TryLock(ctx context.Context, key string, wait ...time.Duration) (bool, func())
	GenerateUID(ctx context.Context, maxID int64) (uint64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	HSet(ctx context.Context, key string, value map[string]any) error
	HMGet(ctx context.Context, key string, fields ...string) ([]any, error)
	HGet(ctx context.Context, key, field string) (string, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Incr(ctx context.Context, key string) (int64, error)

	SetAdd(ctx context.Context, key string, members ...any) (int64, error)
	SetCard(ctx context.Context, key string) (int64, error)
	SetRem(ctx context.Context, key string, members ...any) (int64, error)
	SetIsMember(ctx context.Context, key string, member any) (bool, error)
}
