package utils

import (
	"net"
	"net/http"
	"strings"
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
func RemoteIP(request http.Request) string {
	//fmt.Println(context.Request)
	//fmt.Println(context.Request.Header.Get("X-Forwarded-For"))
	//fmt.Println(context.Request.RemoteAddr)
	//Ali-Cdn-Real-Ip
	IP := request.Header.Get("Ali-Cdn-Real-Ip")
	if strings.EqualFold(IP, "") {
		//_IP := context.Request.Header.Get("X-Forwarded-For")

		IP = strings.Split(request.Header.Get("X-Forwarded-For"), ",")[0]
		if strings.EqualFold(IP, "") {
			text := request.RemoteAddr
			if strings.Contains(text, "::") {
				IP = "0.0.0.0"
			} else {
				IP = strings.Split(text, ":")[0]
			}
		}
	}
	return IP
}
