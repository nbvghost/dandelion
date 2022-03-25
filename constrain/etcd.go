package constrain

import (
	"context"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/etcd"

	"github.com/nbvghost/dandelion/server/serviceobject"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type IEtcd interface {
	Close() error
	Register(desc *serviceobject.ServerDesc) (*serviceobject.ServerDesc, error)
	SelectInsideServer(appName key.MicroServerKey) (string, error)
	SyncConfig(ctx context.Context, key string, callback func(kvs []*clientv3.Event), opts ...clientv3.OpOption)
	GetDNSLocalName(domainName string) (key.MicroServerKey, bool)
	SelectServer(localName key.MicroServerKey) (string, error)
	RegisterRedis(config config.RedisOptions) error
	GetDNSDomains() []string
	ObtainRedis() (*config.RedisOptions, error)
	RegisterPostgresql(dsn string, serverName string) error
	RegisterDNS(dns []etcd.ServerDNS) error
	ObtainPostgresql(serverName string) (string, error)
}
