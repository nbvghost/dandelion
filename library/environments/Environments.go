package environments

import (
	"flag"
	"log"
	_ "net/http/pprof"
	"os"
	"strings"
)

var env environments

type environments struct {
	release       bool
	etcdEndpoints []string
	ListenIP      string
	ListenPort    int
	etcdUsername  string
	etcdPassword  string
}

func init() {
	flag.BoolVar(&env.release, "release", true, "release")
	flag.StringVar(&env.ListenIP, "ip", "", "ip")
	flag.IntVar(&env.ListenPort, "port", 0, "port")

	etcdUsername, ok := os.LookupEnv("ETCD_USERNAME")
	if !ok {
		flag.StringVar(&env.etcdUsername, "etcd_username", "", "etcd_username")
	} else {
		env.etcdUsername = etcdUsername
	}

	etcdPassword, ok := os.LookupEnv("ETCD_PASSWORD")
	if !ok {
		flag.StringVar(&env.etcdPassword, "etcd_password", "", "etcd_password")
	} else {
		env.etcdPassword = etcdPassword
	}

	etcdEndpoints, ok := os.LookupEnv("ETCD_ENDPOINTS")
	if !ok {
		etcdEndpoints = "127.0.0.1:23791,127.0.0.1:23792,127.0.0.1:23793"
	}
	etcdEndpointList := strings.Split(etcdEndpoints, ",")
	for _, v := range etcdEndpointList {
		env.etcdEndpoints = append(env.etcdEndpoints, strings.TrimSpace(v))
	}
	log.Println("FLAG release", env.release)
	log.Println("ENV ETCD_ENDPOINTS", env.etcdEndpoints)
	log.Println("ENV ETCD_USERNAME", env.etcdUsername)
	log.Println("ENV ETCD_PASSWORD", env.etcdPassword)
}
func IP() string {
	return env.ListenIP
}
func Port() int {
	return env.ListenPort
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
