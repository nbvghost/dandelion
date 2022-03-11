package constrain

import (
	"context"
	"github.com/nbvghost/dandelion/config"

	"github.com/nbvghost/dandelion/server/serviceobject"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type MicroServerKey string

const (
	MicroServerKeySSO MicroServerKey = "sso"
	MicroServerKeyOSS MicroServerKey = "oss"
)

type IEtcd interface {
	Close() error
	Register(desc serviceobject.ServerDesc) error
	GetListenPort() int
	SelectServer(appName MicroServerKey) (string, error)
	SelectFileServer() string
	SyncConfig(ctx context.Context, key string, callback func(kvs []*clientv3.Event), opts ...clientv3.OpOption)

	RegisterRedis(config config.RedisOptions) error
	ObtainRedis() (*config.RedisOptions, error)
	RegisterPostgresql(dsn string, serverName string) error
	ObtainPostgresql(serverName string) (string, error)
}
