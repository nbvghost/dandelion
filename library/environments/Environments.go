package environments

import (
	"flag"
	"log"
	"os"
	"strings"
)

var env environments

type environments struct {
	release       bool
	etcdEndpoints []string
}

func init() {
	flag.BoolVar(&env.release, "release", true, "release")

	etcdEndpoints, ok := os.LookupEnv("etcd.endpoints")
	if !ok {
		etcdEndpoints = "127.0.0.1:23791,127.0.0.1:23792,127.0.0.1:23793"
	}
	etcdEndpointList := strings.Split(etcdEndpoints, ",")
	for _, v := range etcdEndpointList {
		env.etcdEndpoints = append(env.etcdEndpoints, strings.TrimSpace(v))
	}
	log.Println("FLAG release", env.release)
	log.Println("ENV etcd.endpoints", env.etcdEndpoints)
}
func Release() bool {
	return env.release
}
func EtcdEndpoints() []string {
	return env.etcdEndpoints
}
