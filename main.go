package main

import "github.com/nbvghost/dandelion/service/http"

func main() {

	/*	r := route.New()
		log.SetFlags(log.LstdFlags)

		conf := config.Config{
			ServerName: "shop",
			Port:       0,
			IP:         "",
			Etcd: clientv3.Config{
				Endpoints:   []string{"172.17.114.159:23791", "172.17.114.159:23792", "172.17.114.159:23793"},
				DialTimeout: 30 * time.Second,
			},
		}
	*/

	http.New(9090).Listen()

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
