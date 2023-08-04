package redis

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
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
func (m *client) TryLock(parentCtx context.Context, key string, wait ...time.Duration) (bool, func()) {

	waitTime := time.Duration(0)
	if len(wait) > 0 {
		waitTime = wait[0]
	}

	_ctx, ctxCancel := context.WithCancel(parentCtx)

	cancel := func() {
		ctxCancel()
		if err := m.getClient().Del(parentCtx, key).Err(); err != nil {
			log.Println(err)
		}
	}
	start := time.Now()

	for time.Now().Sub(start) <= waitTime || waitTime == 0 {
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
			if waitTime == 0 {
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
func (m *client) Incr(ctx context.Context, key string) (int64, error) {
	return m.getClient().Incr(ctx, key).Result()
}
func (m *client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.Kind() == reflect.Struct || v.Kind() == reflect.Map || v.Kind() == reflect.Slice {
			marshal, err := json.Marshal(value)
			if err != nil {
				return err
			}
			return m.getClient().Set(ctx, key, string(marshal), expiration).Err()
		}
	}
	return m.getClient().Set(ctx, key, value, expiration).Err()
}
func (m *client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return m.getClient().Expire(ctx, key, expiration).Err()
}
func (m *client) HSet(ctx context.Context, key string, value map[string]any) error {
	return m.getClient().HSet(ctx, key, value).Err()
}
func (m *client) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	return m.getClient().HMGet(ctx, key, fields...).Result()
}
func (m *client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return m.getClient().Exists(ctx, keys...).Result()
}
func (m *client) HGet(ctx context.Context, key, field string) (string, error) {
	return m.getClient().HGet(ctx, key, field).Result()
}

/*func (m *client) ListPush(ctx context.Context, key string, values ...any) (int64, error) {
	return m.getClient().LPush(ctx, key, values...).Result()
}
func (m *client) ListLen(ctx context.Context, key string) (int64, error) {
	return m.getClient().LLen(ctx, key).Result()
}
func (m *client) ListIndex(ctx context.Context, key string, index int64) (string, error) {
	return m.getClient().LIndex(ctx, key, index).Result()
}
func (m *client) ListLRem(ctx context.Context, key string, value any) (int64, error) {
	return m.getClient().LRem(ctx, key,0, value).Result()
}*/

func (m *client) SetAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return m.getClient().SAdd(ctx, key, members...).Result()
}
func (m *client) SetCard(ctx context.Context, key string) (int64, error) {
	return m.getClient().SCard(ctx, key).Result()
}
func (m *client) SetRem(ctx context.Context, key string, members ...any) (int64, error) {
	return m.getClient().SRem(ctx, key, members...).Result()
}
func (m *client) SetIsMember(ctx context.Context, key string, member any) (bool, error) {
	return m.getClient().SIsMember(ctx, key, member).Result()
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
	c := &client{redis: redis, etcd: etcd}
	c.getClient()
	return c
}
