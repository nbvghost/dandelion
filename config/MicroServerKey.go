package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type serverType string

const (
	ServerTypeHttp serverType = "http"
	ServerTypeGrpc serverType = "grpc"
	ServerTypeDNS  serverType = "dns"
)

type MicroServer struct {
	Name       string
	ServerType serverType
}

func ParseMicroServer(domainName string) (*MicroServer) {
	domainName = strings.Split(domainName, ":")[0]
	return &MicroServer{
		Name:       domainName,
		ServerType: "http",
	}
}

func (m MicroServer) GetAddress() (string, error) {
	filename, err := m.getFileName()
	if err != nil {
		return "", err
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
func (m MicroServer) getFileName() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s/dandelion/%s.%s", dir, m.ServerType, m.Name)

	if d := filepath.Dir(fileName); len(d) > 0 {
		var stat os.FileInfo
		stat, err = os.Stat(d)
		if err != nil {
			err = os.MkdirAll(d, os.ModePerm)
			if err != nil {
				return "", err
			}
		} else {
			if !stat.IsDir() {
				return "", fmt.Errorf("无法创建目录:%s", d)
			}
		}
	}

	return fileName, nil
}

var (
	MicroServerSSO        = MicroServer{Name: "sso.service", ServerType: ServerTypeHttp}
	MicroServerOSS        = MicroServer{Name: "oss.service", ServerType: ServerTypeHttp}
	MicroServerADMIN      = MicroServer{Name: "admin.service", ServerType: ServerTypeHttp}
	MicroServerSITE       = MicroServer{Name: "site.service", ServerType: ServerTypeHttp}
	MicroServersShop      = MicroServer{Name: "shop.service", ServerType: ServerTypeHttp}
	MicroServerMANAGER    = MicroServer{Name: "manager.service", ServerType: ServerTypeHttp}
	MicroServerMimiServer = MicroServer{Name: "mimi.service", ServerType: ServerTypeGrpc}
	MicroServerMiniapp    = MicroServer{Name: "miniapp.service", ServerType: ServerTypeHttp}
)
