package constrain

import (
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/etcd"

	"github.com/nbvghost/dandelion/server/serviceobject"
)

type IEtcd interface {
	Close() error
	Register(desc *serviceobject.ServerDesc) (*serviceobject.ServerDesc, error)
	SelectInsideServer(appName key.MicroServerKey) (string, error)
	GetDNSName(localName key.MicroServerKey) (string, error)
	GetDNSLocalName(domainName string) (key.MicroServerKey, error)
	RegisterRedis(config config.RedisOptions) error
	ObtainRedis() (*config.RedisOptions, error)
	RegisterPostgresql(dsn string, serverName string) error
	RegisterDNS(dns []etcd.ServerDNS) error
	ObtainPostgresql(serverName string) (string, error)
}
