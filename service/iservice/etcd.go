package iservice

import (
	"github.com/nbvghost/dandelion/service/serviceobject"
)

type IEtcd interface {
	Close() error
	Register(desc serviceobject.ServerDesc) error
}
type IEtcdClient interface {
	SelectServer(appName string) string
}
