package utils

import (
	"net"
)

func NetworkIP() string {
	interList, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, inter := range interList {
		//过滤掉系统网络接口没有启用或是回环的地址
		if inter.Flags&net.FlagUp != net.FlagUp || inter.Flags&net.FlagLoopback == net.FlagLoopback {
			continue
		}
		addrList, err := inter.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrList {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String()
				}
			}
		}
	}
	return ""
}
