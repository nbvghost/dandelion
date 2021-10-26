package etcd

import (
	"fmt"
	"log"

	"github.com/nbvghost/dandelion/service/iservice"
	"github.com/nbvghost/dandelion/service/serviceobject"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)

type server struct {
	desc   serviceobject.ServerDesc
	config clientv3.Config
	client *clientv3.Client
}

func (m *server) Close() error {
	return m.client.Close()
}

func (m *server) Register(desc serviceobject.ServerDesc) error {
	m.desc = desc

	client, err := clientv3.New(m.config)
	if err != nil {
		return err
	}

	ctx := context.Background()

	if err = client.Sync(ctx); err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s:%d", m.desc.ServerName, m.desc.IP, m.desc.Port)

	resp, err := client.Get(ctx, key)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.Count == 0 {
		leaseGrantResponse, err := client.Grant(ctx, 0)
		if err != nil {
			log.Fatalln(err)
		}
		putResponse, err := client.Put(ctx, key, fmt.Sprintf("%s:%d", m.desc.IP, m.desc.Port), clientv3.WithLease(leaseGrantResponse.ID))
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(putResponse)

		channel, err := client.KeepAlive(ctx, leaseGrantResponse.ID)
		if err != nil {
			log.Fatalln(err)
		}
		go func() {
			for {
				kvResp := <-channel
				log.Println(kvResp)
			}
		}()
	}
	return nil
}
func NewServer(etcdConfig clientv3.Config) iservice.IEtcd {
	return &server{
		config: etcdConfig,
	}
}
