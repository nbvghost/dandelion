package constrain

import (
	"google.golang.org/grpc"

	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain/key"
)

type IEtcd interface {
	Close() error
	Register(desc *config.MicroServerConfig) (*config.MicroServerConfig, error)
	SelectInsideServer(appName key.MicroServer) (string, error)
	SelectInsideGrpcServer(appName key.MicroServer) (*grpc.ClientConn, error)
	GetDNSName(localName key.MicroServer) (string, error)
	GetDNSLocalName(domainName string) (key.MicroServer, error)
	ObtainRedis() (*config.RedisOptions, error)
	ObtainPostgresql(serverName string) (string, error)
}

type IEtcdAdmin interface {
	RegisterRedis(config config.RedisOptions) error
	RegisterPostgresql(dsn string, serverName string) error
	RegisterDNS(dns []ServerDNS) error
	AddDNS(newDNS []ServerDNS) error
}
