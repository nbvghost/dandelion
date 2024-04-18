package constrain

import (
	"github.com/nbvghost/dandelion/config"
)

type ServerDNS struct {
	DomainName string
	LocalName  config.MicroServer //用于etcd服务发现
}
