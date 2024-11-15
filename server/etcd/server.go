package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/library/environments"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	etcdResolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc/resolver"

	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/pkg/errors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type server struct {
	etcd            clientv3.Config
	client          *clientv3.Client
	resolverBuilder resolver.Builder

	//dnsDomainToServer map[string]config.MicroServer
	dnsDomainToServer *sync.Map
	//dnsServerToDomain map[string][]string
	dnsServerToDomain *sync.Map
	//dnsLocker         sync.RWMutex
	domains           *sync.Map
	//serverMap         map[key.MicroServer][]serviceobject.ServerDesc
	//serverLocker      sync.RWMutex
}

func (m *server) GetMicroServer(domainName string) (*config.MicroServer, error) {
	//m.dnsLocker.RLock()
	//defer m.dnsLocker.RUnlock()
	domainName = strings.Split(domainName, ":")[0]
	//v, ok := m.dnsDomainToServer[domainName]
	v, ok := m.dnsDomainToServer.Load(domainName)
	if !ok {
		v, ok = m.dnsDomainToServer.Load(domainName) //m.dnsDomainToServer[fmt.Sprintf("*.%s", domainName)]
	}
	if !ok {
		return &config.MicroServer{}, errors.Errorf("dns:没有找到(%s)的服务节点", domainName)
	}
	return v.(*config.MicroServer), nil
}

func (m *server) Close() error {
	return m.client.Close()
}

func (m *server) getClient() *clientv3.Client {

	return m.client
}
func (m *server) ObtainRedis() (*config.RedisOptions, error) {
	var err error
	client := m.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := client.Get(ctx, "redis")
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

func (m *server) parseDNS(dns []constrain.ServerDNS, check bool) error {
	//m.dnsLocker.Lock()
	//defer m.dnsLocker.Unlock()

	//m.dnsDomainToServer = map[string]config.MicroServer{}
	//m.dnsServerToDomain = map[string][]string{}
	for _, v := range dns {
		if check {
			//if _, ok := m.dnsDomainToServer[v.DomainName]; ok {
			if _, ok := m.dnsDomainToServer.Load(v.DomainName); ok {
				return errors.Errorf("存在重复的key:DomainName(%s)", v.DomainName)
			}
		}
		//m.dnsDomainToServer[v.DomainName] = v.LocalName
		m.dnsDomainToServer.Store(v.DomainName, v.LocalName)
		//list := m.dnsServerToDomain[v.LocalName.Name]
		var list []string
		{
			v,ok:= m.dnsServerToDomain.Load(v.LocalName.Name)
			if ok{
				list=v.([]string)
			}
		}
		if check {
			var has bool
			for _, n := range list {
				if strings.EqualFold(n, v.DomainName) {
					has = true
					break
				}
			}
			if has {
				return errors.Errorf("存在重复的value:%s的key:LocalName(%s)", v.DomainName, v.LocalName)
			}
		}
		list = append(list, v.DomainName)
		//m.dnsServerToDomain[v.LocalName.Name] = list
		m.dnsServerToDomain.Store(v.LocalName.Name,list)
	}
	return nil
}

// SelectOutsideServer 对外服务地址
func (m *server) SelectOutsideServer(localName *config.MicroServer) (string, error) {
	//m.dnsLocker.RLock()
	//defer m.dnsLocker.RUnlock()
	//list, ok := m.dnsServerToDomain[localName.Name]
	list, ok := m.dnsServerToDomain.Load(localName.Name)
	if !ok {
		return "", errors.Errorf("dns:在获取%s服务时找不到服务地址", localName)
	}
	return list.([]string)[0], nil
}
func (m *server) CheckDomain(domainName string) error {
	domainName=strings.Split(domainName,":")[0]
	_, domain := util.ParseDomain(domainName)
	v, ok := m.domains.Load(domain)
	if !ok {
		return errors.Errorf("域名不允许:%s", domainName)
	}
	if !strings.Contains(v.(string), domainName) {
		return errors.Errorf("域名不在列表中:%s", domainName)
	}

	/*ctx := context.TODO()
	client := m.getClient()

	_, domain := util.ParseDomain(domainName)

	serverKey := fmt.Sprintf("%s/%s", "domains", domain)
	resp, err := client.Get(ctx, serverKey)
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return errors.Errorf("域名不允许:%s", domainName)
	}
	if !strings.Contains(string(resp.Kvs[0].Value), domainName) {
		return errors.Errorf("域名不在列表中:%s", domainName)
	}
	return nil*/
	return nil
}
func (m *server) SelectInsideServer(appName *config.MicroServer) (string, error) {
	ctx := context.TODO()
	client := m.getClient()
	/*if appName.ServerType != config.ServerTypeHttp {
		return "", errors.Errorf("服务不是http服务:%s", appName)
	}*/

	serverKey := fmt.Sprintf("%s/%s/%s/", "server", appName.ServerType, appName.Name)

	resp, err := client.Get(ctx, serverKey, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", errors.Errorf("没有可以用的服务节点:%s", appName)
	}
	v := resp.Kvs[random.Intn(len(resp.Kvs))]
	var serverDesc config.MicroServerConfig
	if err = json.Unmarshal(v.Value, &serverDesc); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", serverDesc.IP, serverDesc.Port), nil
}

/*func (m *server) SelectInsideServer(appName config.MicroServer) (string, error) {
	return appName.Name, nil
}
func (m *server) SelectOutsideServer(appName config.MicroServer) (string, error) {
	return config.GetENV(appName.Name, appName.Name), nil
}*/

func (m *server) watch() {
	var dns []constrain.ServerDNS
	var err error
	var resp *clientv3.GetResponse

	ctx := context.TODO()
	client := m.getClient()

	etcdKey := "dns"
	{
		resp, err = client.Get(ctx, etcdKey)
		if err != nil {
			panic(err)
		}
		if len(resp.Kvs) > 0 {
			if err = json.Unmarshal(resp.Kvs[0].Value, &dns); err != nil {
				panic(err)
			}
			if err = m.parseDNS(dns, false); err != nil {
				panic(err)
			}
			log.Printf("dns list:%+v", dns)
		}
	}
	etcdDomainsKey := "domains"
	{
		resp, err = client.Get(ctx, etcdDomainsKey, clientv3.WithPrefix())
		if err != nil {
			panic(err)
		}
		if len(resp.Kvs) > 0 {

			for i := range resp.Kvs {
				k := string(resp.Kvs[i].Key)
				v := string(resp.Kvs[i].Value)
				dss := strings.Split(k, "/")
				if len(dss) >= 2 {
					m.domains.Store(dss[1], v)
				} else {
					log.Printf("etcd 中domains key 格式不正解：%s", k)
				}
				log.Printf("domains list:%s\t\t\t->\t\t\t%s", k, v)
			}

		}
	}
	/*{
		resp, err = client.Get(ctx, serverKey, clientv3.WithPrefix())
		if err != nil {
			panic(err)
		}
		for _, e := range resp.Kvs {
			var serverDesc serviceobject.ServerDesc
			if err = json.Unmarshal(e.Value, &serverDesc); err != nil {
				log.Println(err)
			}
			if err = m.parseServer(serverDesc, true, false); err != nil {
				log.Println(err)
			}
			log.Printf("server desc:%+v", serverDesc)
		}

	}*/

	//serverWatch := client.Watch(ctx, serverKey, clientv3.WithPrefix())
	dnsWatch := client.Watch(ctx, etcdKey)
	domainsWatch := client.Watch(ctx, etcdDomainsKey, clientv3.WithPrefix())
	go func() {
		for {
			select {
			case domainsResp := <-domainsWatch:
				for i := range domainsResp.Events {
					ev := domainsResp.Events[i]
					k := string(ev.Kv.Key)
					v := string(ev.Kv.Value)
					dss := strings.Split(k, "/")
					if len(dss) >= 2 {
						m.domains.Store(dss[1], v)
						log.Printf("new domains list:%s\t\t\t->\t\t\t%s", dss[1], v)
					} else {
						log.Printf("etcd 中domains key 格式不正解：%s", k)
					}
				}
			case dnsResp := <-dnsWatch:
				for _, ev := range dnsResp.Events {
					//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					if err = json.Unmarshal(ev.Kv.Value, &dns); err != nil {
						log.Println(err)
					}
					if err = m.parseDNS(dns, false); err != nil {
						log.Println(err)
					}
				}
				log.Printf("new dns list:%+v", dns)
			}
		}
	}()
}

/*
	func (m *server) parseServer(serverDesc serviceobject.ServerDesc, isCreate, isModify bool) error {
		m.serverLocker.Lock()
		defer m.serverLocker.Unlock()

		v := m.serverMap[key.MicroServer(serverDesc.Name)]
		if !isCreate && !isModify {
			//删除已经存的
			for i, e := range v {
				if e.IP == serverDesc.IP && e.Port == serverDesc.Port {
					v = append(v[:i], v[i+1:]...)
					log.Printf("del server desc:%+v,isCreate:%v,isModify:%v", serverDesc, isCreate, isModify)
					break
				}
			}
		} else {
			v = append(v, serverDesc)
			log.Printf("new server desc:%+v,isCreate:%v,isModify:%v", serverDesc, isCreate, isModify)
		}
		m.serverMap[key.MicroServer(serverDesc.Name)] = v
		return nil
	}
*/

func (m *server) ObtainPostgresql(serverName string) (string, error) {
	var err error
	client := m.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := client.Get(ctx, fmt.Sprintf("%s/%s", "postgresql", serverName))
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", errors.Errorf("没有到PostgreSQL数据库:%s", serverName)
	}
	return string(resp.Kvs[0].Value), err
}

// Register 注册服务
func (m *server) Register(desc *config.MicroServerConfig) (*config.MicroServerConfig, error) {
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
	desc.Addr = fmt.Sprintf("%s:%d", ip, port)

	etcdKey := fmt.Sprintf("%s/%s/%s/%s:%d", "server", desc.MicroServer.ServerType, desc.MicroServer.Name, desc.IP, desc.Port)

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
					time.Sleep(time.Second * 10)
					continue
				}
				_, err = client.Put(ctx, etcdKey, string(vBytes), clientv3.WithLease(leaseGrant.ID))
				if err != nil {
					log.Println(err)
					time.Sleep(time.Second * 10)
					continue
				}
				channel, err = client.KeepAlive(ctx, leaseGrant.ID)
				if err != nil {
					log.Println(err)
					time.Sleep(time.Second * 10)
					continue
				}
			}
		}
	}()

	return desc, nil
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// SelectInsideGrpcServer  服务间通信通过这个方法获取
/*func (m *server) SelectInsideGrpcServer(appName config.MicroServer) (*grpc.ClientConn, error) {
	if appName.ServerType != config.ServerTypeGrpc {
		return nil, errors.Errorf("服务不是grpc服务:%s", appName)
	}
	ctx := context.TODO()

	d, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s/%s/%s", m.resolverBuilder.Scheme(), "server", appName.ServerType, appName.Name), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                0,
		Timeout:             0,
		PermitWithoutStream: false,
	}), grpc.WithInsecure(), grpc.WithResolvers(m.resolverBuilder))
	if err != nil {
		return nil, err
	}
	return d, nil
}*/
func NewServer(config clientv3.Config) constrain.IEtcd {
	log.Println("connetct to etcd server")
	defer log.Println("success connetct to etcd server")
	if !environments.EtcdAble() {
		log.Println("use default etcd interface")
		return NewDefaultEtcd()
	}
	client, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}

	r, err := etcdResolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}

	resolver.Register(r)

	s := &server{etcd: config,
		client:            client,
		dnsDomainToServer: &sync.Map{},
		dnsServerToDomain: &sync.Map{},
		//dnsServerToDomain: map[string][]string{},
		//dnsDomainToServer: map[string]config.MicroServer{},
		//serverMap:         map[key.MicroServer][]serviceobject.ServerDesc{},
		resolverBuilder: r,
		domains:         &sync.Map{},
	}
	s.watch()

	resp, err := client.Get(context.TODO(), "", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, v := range resp.Kvs {
		log.Println(fmt.Sprintf("%s:%s", string(v.Key), v.Value))
	}

	/*em, err := endpoints.NewManager(client, "grpc server")
	log.Println(err)
	em.AddEndpoint(client.Ctx(), "grpc server"+"/"+"addsdfdsr", endpoints.Endpoint{Addr: "adsdfdsdr"})*/

	/*etcdResolver, err := resolver.NewBuilder(client)
	log.Println(err)

	d, errr := grpc.Dial("etcd:///grpc server", grpc.WithInsecure(), grpc.WithResolvers(etcdResolver))
	log.Println(errr)
	d.Connect()
	log.Println(d.GetState(), d.Target())*/
	return s
}
