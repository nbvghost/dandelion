package main

import (
	"github.com/nbvghost/dandelion/service/etcd"
	"github.com/nbvghost/dandelion/service/grpc"
	"github.com/nbvghost/dandelion/service/http"
	"log"
)

func main() {

	//r := route.New()
	log.SetFlags(log.LstdFlags)

	conf := etcd.New("config.json")

	etcdService := etcd.NewServer(conf.Etcd)

	http.New(9090, grpc.NewClient(etcdService)).Listen()

	/*etcdService := etcd.New(conf.Etcd)

	defer func() {
		etcdService.Close()
	}()
	grpc.New(conf, r, func(desc serviceobject.ServerDesc) {

		if err := etcdService.Register(desc); err != nil {
			log.Fatalln(err)
		}

	}).Listen()*/

}
