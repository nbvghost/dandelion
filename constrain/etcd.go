package constrain

import (
	"google.golang.org/grpc"

	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/etcd"

	"github.com/nbvghost/dandelion/server/serviceobject"
)

type IEtcd interface {
	Close() error
	Register(desc *serviceobject.ServerDesc) (*serviceobject.ServerDesc, error)
	SelectInsideServer(appName key.MicroServer) (string, error)
	SelectInsideGrpcServer(appName key.MicroServer) (*grpc.ClientConn, error)
	GetDNSName(localName key.MicroServer) (string, error)
	GetDNSLocalName(domainName string) (key.MicroServer, error)
	RegisterRedis(config config.RedisOptions) error
	ObtainRedis() (*config.RedisOptions, error)
	RegisterPostgresql(dsn string, serverName string) error
	RegisterDNS(dns []etcd.ServerDNS) error
	AddDNS(newDNS []etcd.ServerDNS) error
	ObtainPostgresql(serverName string) (string, error)
}
