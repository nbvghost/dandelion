package environments

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"strings"
)

var env environments

type environments struct {
	release       bool
	listenIP      string
	listenPort    int
	etcdAble      bool
	redisAble     bool
	etcdEndpoints []string
	etcdUsername  string
	etcdPassword  string
}

func init() {
	flag.BoolVar(&env.release, "release", true, "release")
	flag.StringVar(&env.listenIP, "ip", "", "ip")
	flag.IntVar(&env.listenPort, "port", 0, "port")

	etcdEndpoints := GetENV("ETCD_ENDPOINTS", "127.0.0.1:23791,127.0.0.1:23792,127.0.0.1:23793")
	env.etcdEndpoints = strings.Split(etcdEndpoints, ",")
	env.etcdUsername = GetENV("ETCD_USERNAME", "")
	env.etcdPassword = GetENV("ETCD_PASSWORD", "")
	env.etcdAble = GetENV("ETCD_ABLE", "true") == "true"
	env.redisAble = GetENV("REDIS_ABLE", "true") == "true"
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
func EtcdAble() bool {
	return env.etcdAble
}
func RedisAble() bool {
	return env.redisAble
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

var envMap = map[string]string{}

func GetENV(key, defaultValue string) string {
	if v, ok := envMap[key]; ok {
		return v
	}
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	log.Println(fmt.Sprintf("env %s %s", key, value))
	envMap[key] = value
	return value
}
