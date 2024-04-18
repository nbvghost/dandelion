package constrain

import (
	"google.golang.org/grpc"

	"github.com/nbvghost/dandelion/config"
)

type IEtcd interface {
	Close() error
	Register(desc *config.MicroServerConfig) (*config.MicroServerConfig, error)
	SelectInsideServer(appName config.MicroServer) (string, error)
	SelectOutsideServer(appName config.MicroServer) (string, error)
	SelectInsideGrpcServer(appName config.MicroServer) (*grpc.ClientConn, error)
	//GetDNSName(localName key.MicroServer) (string, error)
	//GetDNSLocalName(domainName string) (config.MicroServer, error)
	ObtainRedis() (*config.RedisOptions, error)
	ObtainPostgresql(serverName string) (string, error)
}

type IEtcdAdmin interface {
	RegisterRedis(config config.RedisOptions) error
	RegisterPostgresql(dsn string, serverName string) error
	RegisterDNS(dns []ServerDNS) error
	AddDNS(newDNS []ServerDNS) error
}
