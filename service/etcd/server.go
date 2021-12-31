package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/library/action"
	"math/rand"
	"sync"
	"time"

	"github.com/nbvghost/dandelion/service/serviceobject"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type server struct {
	etcd   clientv3.Config
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
func (m *server) SelectFileServer() string {
	return "http://127.0.0.1/file"
}
func (m *server) SelectServer(appName string) (string, error) {
	ctx := context.Background()
	resp, err := m.getClient().Get(ctx, appName, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", action.NewCodeWithError(action.Error, fmt.Errorf("没有可以用的服务节点:%s", appName))
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return string(resp.Kvs[r.Intn(len(resp.Kvs))].Value), nil
}
func (m *server) getClient() *clientv3.Client {
	m.once.Do(func() {
		client, err := clientv3.New(m.etcd)
		if err != nil {
			panic(err)
		}
		m.client = client
	})
	return m.client
}
func (m *server) ObtainRedis() (*config.RedisOptions, error) {
	var err error
	client := m.getClient()
	ctx := context.Background()

	resp, err := client.Get(ctx, "Redis")
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("没有到redis节点")
	}

	op := config.RedisOptions{}
	if err = json.Unmarshal(resp.Kvs[0].Value, &op); err != nil {
		return nil, err
	}
	return &op, nil
}
func (m *server) RegisterRedis(config config.RedisOptions) error {
	var err error
	client := m.getClient()
	ctx := context.Background()

	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = client.Put(ctx, "Redis", string(b))
	if err != nil {
		return err
	}
	return nil
}
func (m *server) ObtainPostgresql() (string, error) {
	var err error
	client := m.getClient()
	ctx := context.Background()
	resp, err := client.Get(ctx, "Postgresql")
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("没有到Postgresql节点")
	}

	return string(resp.Kvs[0].Value), err
}
func (m *server) RegisterPostgresql(dsn string) error {
	var err error
	client := m.getClient()
	ctx := context.Background()

	_, err = client.Put(ctx, "Postgresql", dsn)
	if err != nil {
		return err
	}
	return nil
}
func (m *server) Register(desc serviceobject.ServerDesc) error {
	var err error
	client := m.getClient()

	ctx := context.Background()

	/*if err = client.Sync(ctx); err != nil {
		return err
	}*/

	key := fmt.Sprintf("%s/%s:%d", desc.Name, desc.IP, desc.Port)

	_, err = client.Get(ctx, key)
	if err != nil {
		return err
	}

	leaseGrant, err := client.Grant(ctx, 10)
	if err != nil {
		return err
	}
	_, err = client.Put(ctx, key, fmt.Sprintf("%s:%d", desc.IP, desc.Port), clientv3.WithLease(leaseGrant.ID))
	if err != nil {
		return err
	}

	channel, err := client.KeepAlive(ctx, leaseGrant.ID)
	if err != nil {
		return err
	}
	go func() {
		for {
			<-channel
		}
	}()

	return nil
}
func NewServer(config clientv3.Config) IEtcd {
	return &server{
		etcd: config,
	}
}
