package util

import (
	"regexp"
	"strings"
)

var ipRegexp = regexp.MustCompile(`(\d)+\.(\d)+\.(\d)+\.(\d)+`)

func ParseDomain(host string) string {
	hosts := strings.Split(host, ":")
	var domainName string
	domains := strings.Split(hosts[0], ".")
	if len(domains) == 1 {
		domainName = domains[0]
	} else {
		if !ipRegexp.MatchString(hosts[0]) {
			domainName = domains[len(domains)-2] + "." + domains[len(domains)-1]
		} else {
			return ""
		}
	}
	return domainName
}
