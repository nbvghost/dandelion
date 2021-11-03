package etcd

import (
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/library/result"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/nbvghost/dandelion/service/iservice"
	"github.com/nbvghost/dandelion/service/serviceobject"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)

type server struct {
	config *config.ServerConfig
	client *clientv3.Client
	nodes  sync.Map
	once   sync.Once
}

func (m *server) Close() error {
	return m.client.Close()
}

func (m *server) SyncConfig(ctx context.Context, key string, callback func(kvs []*clientv3.Event), opts ...clientv3.OpOption) {
	channel := m.getClient().Watch(ctx, key, opts...)
	var compactRevision int64
	for c := range channel {
		if compactRevision != c.CompactRevision {
			callback(c.Events)
		}
	}
}
func (m *server) SelectServer(appName string) (string, error) {
	ctx := context.Background()
	resp, err := m.getClient().Get(ctx, appName, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", result.NewCodeWithError(result.Error, fmt.Errorf("没有可以用的服务节点:%s", appName))
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return string(resp.Kvs[r.Intn(len(resp.Kvs))].Value), nil
}
func (m *server) getClient() *clientv3.Client {
	m.once.Do(func() {
		client, err := clientv3.New(*m.config.Etcd)
		if err != nil {
			panic(err)
		}
		m.client = client
	})
	return m.client
}
func (m *server) Register(desc serviceobject.ServerDesc) error {

	var err error
	client := m.getClient()
	m.client = client

	ctx := context.Background()

	if err = client.Sync(ctx); err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s:%d", desc.ServerName, desc.IP, desc.Port)

	_, err = client.Get(ctx, key)
	if err != nil {
		log.Fatalln(err)
	}

	leaseGrant, err := client.Grant(ctx, 10)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = client.Put(ctx, key, fmt.Sprintf("%s:%d", desc.IP, desc.Port), clientv3.WithLease(leaseGrant.ID))
	if err != nil {
		log.Fatalln(err)
	}

	channel, err := client.KeepAlive(ctx, leaseGrant.ID)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			<-channel
		}
	}()

	return nil
}
func NewServer(config *config.ServerConfig) iservice.IEtcd {
	return &server{
		config: config,
	}
}
