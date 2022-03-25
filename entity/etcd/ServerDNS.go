package etcd

import (
	"github.com/nbvghost/dandelion/constrain/key"
)

type ServerDNS struct {
	Name      string
	LocalName key.MicroServerKey //用于etcd服务发现
	Env       string
}
