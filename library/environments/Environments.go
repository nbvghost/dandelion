package environments

import (
	"flag"
	"github.com/nbvghost/dandelion/config"
	"log"
	_ "net/http/pprof"
	"strings"
)

var env environments

type environments struct {
	release       bool
	listenIP      string
	listenPort    int
	etcdEndpoints []string
	etcdUsername  string
	etcdPassword  string
}

func init() {
	flag.BoolVar(&env.release, "release", true, "release")
	flag.StringVar(&env.listenIP, "ip", "", "ip")
	flag.IntVar(&env.listenPort, "port", 0, "port")

	etcdEndpoints := config.GetENV("ETCD_ENDPOINTS", "127.0.0.1:23791,127.0.0.1:23792,127.0.0.1:23793")
	env.etcdEndpoints = strings.Split(etcdEndpoints, ",")
	env.etcdUsername = config.GetENV("ETCD_USERNAME", "")
	env.etcdPassword = config.GetENV("ETCD_PASSWORD", "")
}
func Print() {
	log.Println("FLAG release", env.release)
	log.Println("FLAG ip", env.listenIP)
	log.Println("FLAG port", env.listenPort)
}
func IP() string {
	return env.listenIP
}
func Port() int {
	return env.listenPort
}
func Release() bool {
	return env.release
}
func EtcdEndpoints() []string {
	return env.etcdEndpoints
}
func EtcdUsername() string {
	return env.etcdUsername
}
func EtcdPassword() string {
	return env.etcdPassword
}
