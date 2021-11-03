package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/nbvghost/dandelion/service/iservice"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"time"
)

type client struct {
	sync.RWMutex
	once   sync.Once
	client *redis.ClusterClient
	etcd   iservice.IEtcd
}

func (m *client) Get(ctx context.Context, key string) (string, error) {
	return m.getClient().Get(ctx, key).Result()
}
func (m *client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.getClient().Set(ctx, key, value, expiration).Err()
}
func (m *client) getClient() redis.Cmdable {
	m.RLock()
	defer m.RUnlock()
	if m.client == nil {
		var wg sync.WaitGroup
		wg.Add(1)
		m.once.Do(func() {
			go m.etcd.SyncConfig(context.TODO(), "redis", func(kvs []*clientv3.Event) {
				m.Lock()
				defer m.Unlock()

				var addrList []string
				for i := range kvs {
					addrList = append(addrList, string(kvs[i].Kv.Value))
				}
				c := redis.NewClusterClient(&redis.ClusterOptions{Addrs: addrList})
				if r := c.Ping(context.TODO()); r.Err() != nil {
					log.Fatalln(r.Err())
				}
				m.client = c
				wg.Done()
			}, clientv3.WithPrefix())
		})
		wg.Wait()
	}
	return m.client
}
func NewClient(etcd iservice.IEtcd) iservice.IRedis {
	return &client{etcd: etcd}
}
