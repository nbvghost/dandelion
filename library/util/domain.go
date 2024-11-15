package util

import (
	"regexp"
	"strings"
)

var ipRegexp = regexp.MustCompile(`(\d)+\.(\d)+\.(\d)+\.(\d)+`)

// 分析域名
// 返回域名前缀，域名
func ParseDomain(host string) ([]string, string) {
	hosts := strings.Split(host, ":")
	var domainName string
	var domainPrefix []string
	domains := strings.Split(hosts[0], ".")
	if len(domains) == 1 {
		domainName = domains[0]
	} else {
		if !ipRegexp.MatchString(hosts[0]) {
			domainName = domains[len(domains)-2] + "." + domains[len(domains)-1]
			for i := 0; i < len(domains)-2; i++ {
				domainPrefix = append(domainPrefix, domains[i])
			}
		} else {
			return nil, ""
		}
	}
	return domainPrefix, domainName
}
