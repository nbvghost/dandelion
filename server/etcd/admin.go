package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"log"
	"sync"
)

type adminServer struct {
	server *server
}

func (m *adminServer) RegisterRedis(config *config.RedisOptions) error {
	var err error
	client := m.server.getClient()
	ctx := context.Background()

	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = client.Put(ctx, "redis", string(b))
	if err != nil {
		return err
	}
	return nil
}
func (m *adminServer) RegisterPostgresql(dsn string, serverName string) error {
	var err error
	client := m.server.getClient()
	ctx := context.Background()
	_, err = client.Put(ctx, fmt.Sprintf("%s/%s", "postgresql", serverName), dsn)
	if err != nil {
		return err
	}
	return nil
}

func (m *adminServer) RegisterDNS(dns []constrain.ServerDNS) error {
	copyServer := &server{
		dnsDomainToServer: &sync.Map{},
		dnsServerToDomain: &sync.Map{},
		dnsLocker:         sync.RWMutex{},
	}
	if err := copyServer.parseDNS(dns, true); err != nil {
		return err
	}
	client := m.server.getClient()
	/*for i := range dns {
		if !strings.Contains(dns[i].DomainName, "*.") {
			err := m.AddDomains(dns[i].DomainName, []string{dns[i].DomainName})
			if err != nil {
				return err
			}
		}
	}*/
	etcdKey := "dns"
	ctx := context.TODO()
	jsonByte, err := json.Marshal(dns)
	if err != nil {
		return err
	}
	_, err = client.Put(ctx, etcdKey, string(jsonByte))
	if err != nil {
		return err
	}
	return nil
}

func (m *adminServer) AddDomains(domainName string, domainNames []string) error {
	ctx := context.TODO()
	client := m.server.getClient()
	serverKey := fmt.Sprintf("%s/%s", "domains", domainName)

	domainNamesJson, err := json.Marshal(domainNames)
	if err != nil {
		return err
	}
	_, err = client.Put(ctx, serverKey, string(domainNamesJson))
	if err != nil {
		return err
	}
	return nil
}
func (m *adminServer) AddDNS(newDNS []constrain.ServerDNS) error {
	etcdKey := "dns"
	ctx := context.TODO()

	var hasDns []constrain.ServerDNS

	client := m.server.getClient()

	resp, err := client.Get(ctx, etcdKey)
	if err != nil {
		return err
	}
	if len(resp.Kvs) > 0 {
		err := json.Unmarshal(resp.Kvs[0].Value, &hasDns)
		if err != nil {
			log.Println(err)
		}
	}

	hasDns = append(hasDns, newDNS...)

	copyServer := &server{dnsServerToDomain: &sync.Map{}, dnsDomainToServer: &sync.Map{}}
	if err := copyServer.parseDNS(hasDns, true); err != nil {
		return err
	}

	jsonByte, err := json.Marshal(hasDns)
	if err != nil {
		return err
	}
	_, err = client.Put(ctx, etcdKey, string(jsonByte))
	if err != nil {
		return err
	}
	return nil
}

func NewAdminServer(clientServer constrain.IEtcd) constrain.IEtcdAdmin {
	//s := NewServer(config).(*server)
	return &adminServer{server: clientServer.(*server)}
}
