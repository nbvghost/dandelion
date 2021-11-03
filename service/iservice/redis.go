package iservice

import "context"

type IRedis interface {
	Get(ctx context.Context, key string) (string, error)
}
