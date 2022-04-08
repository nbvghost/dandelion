package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"log"
	"sync"
	"time"
)

type client struct {
	sync.RWMutex
	once   sync.Once
	client *redis.Client
	etcd   constrain.IEtcd
	redis  config.RedisOptions
}

func (m *client) GetEtcd() constrain.IEtcd {
	return m.etcd
}
func (m *client) TryLock(parentCtx context.Context, key string, timeouts ...time.Duration) (bool, func()) {
	timeout := time.Duration(0)
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	_ctx, ctxCancel := context.WithCancel(parentCtx)

	cancel := func() {
		ctxCancel()
		if err := m.getClient().Del(parentCtx, key).Err(); err != nil {
			log.Println(err)
		}
	}
	start := time.Now()

	for time.Now().Sub(start) <= timeout || timeout == 0 {
		ok := m.getClient().SetNX(_ctx, key, "lock", time.Minute)
		if ok.Val() {
			//获取锁成功
			go func() {
				t := time.NewTicker(time.Minute - (time.Second - 10))
				defer t.Stop()
				for {
					select {
					case <-_ctx.Done():
						return
					case <-t.C:
						expireOK := m.getClient().Expire(_ctx, key, time.Minute)
						if !expireOK.Val() {
							log.Println("lock设置key过期时间失败")
						}
					}
				}
			}()
			return true, cancel
		} else {
			//获取锁失败
			if timeout == 0 {
				break
			}
			time.Sleep(time.Second)
		}
	}
	return false, nil
}
func (m *client) Del(ctx context.Context, keys ...string) (int64, error) {
	return m.getClient().Del(ctx, keys...).Result()
}
func (m *client) Get(ctx context.Context, key string) (string, error) {
	return m.getClient().Get(ctx, key).Result()
}
func (m *client) GenerateUID(ctx context.Context) uint64 {
	key := NewUIDKey()
	mUID := m.getClient().Get(ctx, key)
	v, _ := mUID.Uint64()
	if v == 0 {
		v, _ = m.getClient().IncrBy(ctx, key, 100000).Uint64()
	}
	v, _ = m.getClient().Incr(ctx, key).Uint64()
	return v
}
func (m *client) GetEx(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return m.getClient().GetEx(ctx, key, expiration).Result()
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
			op := m.redis.ToOptions()
			c := redis.NewClient(&op)
			if r := c.Ping(context.TODO()); r.Err() != nil {
				log.Fatalln(r.Err())
			}
			m.client = c
			wg.Done()
			/*go m.etcd.SyncConfig(context.TODO(), "redis", func(kvs []*clientv3.Event) {
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
			}, clientv3.WithPrefix())*/
		})
		wg.Wait()
	}
	return m.client
}
func NewClient(redis config.RedisOptions, etcd constrain.IEtcd) constrain.IRedis {
	return &client{redis: redis, etcd: etcd}
}
