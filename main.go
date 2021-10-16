package main

import (
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/service/etcd"
	"github.com/nbvghost/dandelion/service/grpc"
	"github.com/nbvghost/dandelion/service/workobject"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func main() {

	conf := config.Config{
		ServerName: "shop",
		Port:       0,
		IP:         "",
		Etcd: clientv3.Config{
			Endpoints:   []string{"0.0.0.0:23791", "0.0.0.0:23792", "0.0.0.0:23793"},
			DialTimeout: 30 * time.Second,
		},
	}

	grpc.New(conf, func(desc workobject.ServerDesc) {

		if err := etcd.New(conf.Etcd, desc).Register(); err != nil {
			log.Fatalln(err)
		}

	}).Listen()

}
