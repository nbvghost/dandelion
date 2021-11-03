package iservice

import (
	"context"
	"github.com/nbvghost/dandelion/service/serviceobject"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type IEtcd interface {
	Close() error
	Register(desc serviceobject.ServerDesc) error
	SelectServer(appName string) (string, error)
	SyncConfig(ctx context.Context, key string, callback func(kvs []*clientv3.Event), opts ...clientv3.OpOption)
}
