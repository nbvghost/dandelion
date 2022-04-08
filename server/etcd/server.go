package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/etcd"
	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/pkg/errors"

	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/nbvghost/dandelion/server/serviceobject"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type server struct {
	etcd              clientv3.Config
	client            *clientv3.Client
	nodes             sync.Map
	dnsDomainToServer map[string]key.MicroServerKey
	dnsServerToDomain map[string][]string
	dnsDomains        []string
	sync.RWMutex
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

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

//服务间通信通过这个方法获取
func (m *server) SelectInsideServer(appName key.MicroServerKey) (string, error) {
	ctx := context.Background()
	resp, err := m.getClient().Get(ctx, string(appName), clientv3.WithPrefix())
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", action.NewCodeWithError(action.Error, errors.Errorf("没有可以用的服务节点:%s", appName))
	}

	valueByte := resp.Kvs[r.Intn(len(resp.Kvs))].Value
	var serverDesc serviceobject.ServerDesc
	if err = json.Unmarshal(valueByte, &serverDesc); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", serverDesc.IP, serverDesc.Port), nil
}

func (m *server) getClient() *clientv3.Client {

	return m.client
}
func (m *server) ObtainRedis() (*config.RedisOptions, error) {
	var err error
	client := m.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := client.Get(ctx, "Redis")
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, errors.Errorf("没有到redis节点")
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
func (m *server) parseDNS(dns []etcd.ServerDNS, check bool) error {
	defer m.Unlock()
	m.Lock()
	for _, v := range dns {
		if check {
			if _, ok := m.dnsDomainToServer[v.Name]; ok {
				return errors.Errorf("存在重复的key:Name(%s)", v.Name)
			}
		}

		m.dnsDomainToServer[v.Name] = v.LocalName
		//------------------
		list := m.dnsServerToDomain[string(v.LocalName)]
		if check {
			var has bool
			for _, n := range list {
				if strings.EqualFold(n, v.Name) {
					has = true
					break
				}
			}
			if has {
				return errors.Errorf("存在重复的value:%s的key:LocalName(%s)", v.Name, v.LocalName)
			}
		}
		list = append(list, v.Name)
		m.dnsServerToDomain[string(v.LocalName)] = list
		m.dnsDomains = append(m.dnsDomains, v.Name)
	}
	return nil
}
func (m *server) GetDNSDomains() []string {
	return m.dnsDomains
}

/*func (m *server) getDNSEnv() string {
	var env string
	if environments.Release() {
		env = "release"
	} else {
		env = "dev"
	}
	return env
}*/

//对外服务地址
func (m *server) SelectServer(localName key.MicroServerKey) (string, error) {
	//env := m.getDNSEnv()
	list, ok := m.dnsServerToDomain[string(localName)]
	if !ok || len(list) == 0 {
		return "", errors.Errorf("在获取%s服务时找不到服务地址", localName)
	}
	return list[0], nil
}
func (m *server) GetDNSLocalName(domainName string) (key.MicroServerKey, bool) {
	//env := m.getDNSEnv()
	v, ok := m.dnsDomainToServer[domainName]
	if !ok {
		v, ok = m.dnsDomainToServer[fmt.Sprintf("*.%s", domainName)]
	}
	return v, ok
}
func (m *server) watchDNS() {
	var dns []etcd.ServerDNS

	ctx := context.TODO()
	client := m.getClient()
	etcdKey := "ServerDNS"
	resp, err := client.Get(ctx, etcdKey)
	if err != nil {
		panic(err)
	}
	if len(resp.Kvs) > 0 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &dns); err != nil {
			panic(err)
		}
		if err := m.parseDNS(dns, false); err != nil {
			panic(err)
		}
	}
	c := client.Watch(ctx, etcdKey)
	go func() {
		for resp := range c {
			for _, ev := range resp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

				if err := json.Unmarshal(ev.Kv.Value, &dns); err != nil {
					panic(err)
				}
				if err := m.parseDNS(dns, false); err != nil {
					panic(err)
				}

			}
		}
	}()
}

func (m *server) RegisterDNS(dns []etcd.ServerDNS) error {
	copyServer := &server{dnsServerToDomain: map[string][]string{}, dnsDomainToServer: map[string]key.MicroServerKey{}}
	if err := copyServer.parseDNS(dns, true); err != nil {
		return err
	}
	client := m.getClient()

	etcdKey := "ServerDNS"

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
func (m *server) ObtainPostgresql(serverName string) (string, error) {
	var err error
	client := m.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := client.Get(ctx, fmt.Sprintf("%s/%s", "Postgresql", serverName))
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", errors.Errorf("没有到Postgresql节点")
	}
	return string(resp.Kvs[0].Value), err
}
func (m *server) RegisterPostgresql(dsn string, serverName string) error {
	var err error
	client := m.getClient()
	ctx := context.Background()

	_, err = client.Put(ctx, fmt.Sprintf("%s/%s", "Postgresql", serverName), dsn)
	if err != nil {
		return err
	}
	return nil
}
func (m *server) Register(desc *serviceobject.ServerDesc) (*serviceobject.ServerDesc, error) {
	var err error
	client := m.getClient()

	ctx := context.Background()

	/*if err = client.Sync(ctx); err != nil {
		return err
	}*/

	var ip = desc.IP
	var port = desc.Port
	if ip == "" {
		ip = util.NetworkIP()
		if ip == "" {
			return nil, errors.New("无法获取本机ip")
		}
	}
	if port == 0 {
		port, err = util.RandomNetworkPort()
		if err != nil {
			return nil, err
		}
	}

	desc.IP = ip
	desc.Port = port

	etcdKey := fmt.Sprintf("%s/%s:%d", desc.Name, ip, port)

	_, err = client.Get(ctx, etcdKey)
	if err != nil {
		return nil, err
	}

	var vBytes []byte
	vBytes, err = json.Marshal(&desc)
	if err != nil {
		return nil, err
	}

	leaseGrant, err := client.Grant(ctx, 10)
	if err != nil {
		return nil, err
	}
	_, err = client.Put(ctx, etcdKey, string(vBytes), clientv3.WithLease(leaseGrant.ID))
	if err != nil {
		return nil, err
	}

	channel, err := client.KeepAlive(ctx, leaseGrant.ID)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			leaseKeepAliveResponse := <-channel
			if leaseKeepAliveResponse == nil {
				leaseGrant, err = client.Grant(ctx, 10)
				if err != nil {
					log.Println(err)
					return
				}
				_, err = client.Put(ctx, etcdKey, string(vBytes), clientv3.WithLease(leaseGrant.ID))
				if err != nil {
					log.Println(err)
				}
				channel, err = client.KeepAlive(ctx, leaseGrant.ID)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}()

	return desc, nil
}
func NewServer(config clientv3.Config) constrain.IEtcd {
	client, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}
	s := &server{etcd: config, client: client, dnsServerToDomain: map[string][]string{}, dnsDomainToServer: map[string]key.MicroServerKey{}}
	s.watchDNS()
	return s
}
