package etcd

import (
	"github.com/nbvghost/dandelion/constrain/key"
)

type ServerDNS struct {
	DomainName string
	LocalName  key.MicroServer //用于etcd服务发现
}
