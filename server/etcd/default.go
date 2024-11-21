package etcd

import (
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/environments"
)

type DefaultEtcd struct{}

func (m *DefaultEtcd) CheckDomain(domainName string) error {
	return nil
}

func (m *DefaultEtcd) Close() error {
	return nil
}

func (m *DefaultEtcd) Register(desc *config.MicroServerConfig) (*config.MicroServerConfig, error) {
	return desc, nil
}

func (m *DefaultEtcd) SelectInsideServer(appName *config.MicroServer) (string, error) {
	return environments.GetENV(fmt.Sprintf("microserver.%s", appName.Name), "127.0.0.1"), nil
}

func (m *DefaultEtcd) SelectOutsideServer(appName *config.MicroServer) (string, error) {
	return environments.GetENV(fmt.Sprintf("microserver.%s", appName.Name), "127.0.0.1"), nil
}

func (m *DefaultEtcd) GetMicroServer(domainName string) (*config.MicroServer, error) {
	return &config.MicroServer{}, nil
}

func (m *DefaultEtcd) ObtainRedis() (*config.RedisOptions, error) {
	return &config.RedisOptions{
		Addr: environments.GetENV("REDIS", "127.0.0.1:6379"),
		DB:   0,
	}, nil
}

func (m *DefaultEtcd) ObtainPostgresql(serverName string) (string, error) {
	return environments.GetENV("POSTGRESQL", "host=172.19.32.72 user=postgres password=274455411 dbname=dandelion port=5432 sslmode=disable TimeZone=Asia/Shanghai"), nil
}
func NewDefaultEtcd() constrain.IEtcd {

	return &DefaultEtcd{}

}

var Default constrain.IEtcd = &DefaultEtcd{}
