package config

import clientv3 "go.etcd.io/etcd/client/v3"

type Config struct {
	ServerName string
	Port       int
	IP         string
	Etcd       clientv3.Config
}
