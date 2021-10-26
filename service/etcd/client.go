package etcd

import (
	"sync"

	"github.com/nbvghost/dandelion/service/iservice"
)

type client struct {
	nodes sync.Map
}

func (m *client) SelectServer(appName string) string {

	return ""
}
func NewClient() iservice.IEtcdClient {
	return &client{}
}
