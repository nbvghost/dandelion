package etcd

import (
	"fmt"
	"github.com/nbvghost/dandelion/service/workobject"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"log"
)

type service struct {
	desc workobject.ServerDesc
	etcd clientv3.Config
}

func (m *service) Register() error {
	client, err := clientv3.New(m.etcd)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("/%s/%s:%d/", m.desc.ServerName, m.desc.IP, m.desc.Port)

	ctx := context.Background()

	resp, err := client.Get(ctx, key)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.Count == 0 {
		leaseGrantResponse, err := client.Grant(ctx, 55)
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
				<-channel
			}
		}()
	}
	return nil
}
func New(etcd clientv3.Config, desc workobject.ServerDesc) *service {
	return &service{
		etcd: etcd,
		desc: desc,
	}
}
