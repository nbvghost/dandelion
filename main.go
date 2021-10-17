package main

import (
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/service/etcd"
	"github.com/nbvghost/dandelion/service/grpc"
	"github.com/nbvghost/dandelion/service/route"
	"github.com/nbvghost/dandelion/service/serviceobject"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func main() {

	r := route.New()
	log.SetFlags(log.LstdFlags)

	conf := config.Config{
		ServerName: "shop",
		Port:       0,
		IP:         "",
		Etcd: clientv3.Config{
			Endpoints:   []string{"0.0.0.0:23791", "0.0.0.0:23792", "0.0.0.0:23793"},
			DialTimeout: 30 * time.Second,
		},
	}

	etcdService := etcd.New(conf.Etcd)

	defer func() {
		etcdService.Close()
	}()
	grpc.New(conf, r, func(desc serviceobject.ServerDesc) {

		if err := etcdService.Register(desc); err != nil {
			log.Fatalln(err)
		}

	}).Listen()

}
