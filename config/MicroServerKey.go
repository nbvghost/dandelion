package config

import (
	"fmt"
	"github.com/nbvghost/dandelion/library/environments"
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

func ParseMicroServer(domainName string) *MicroServer {
	domainName = strings.Split(domainName, ":")[0]
	return &MicroServer{
		Name:       domainName,
		ServerType: "http",
	}
}
func (m MicroServer) SelectInsideServer() (string, error) {
	if environments.Release() {
		return GetENV(m.Name, m.Name), nil
	}
	return m.GetAddress()
}
func (m MicroServer) SelectOutsideServer() (string, error) {
	if environments.Release() {
		return GetENV(m.Name, m.Name), nil
	}
	return m.Name, nil
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

	fileName := fmt.Sprintf("%s/dandelion/%s/%s", dir, m.Name, m.ServerType)

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
	MicroServerSSO        = MicroServer{Name: "sso", ServerType: ServerTypeHttp}
	MicroServerOSS        = MicroServer{Name: "oss", ServerType: ServerTypeHttp}
	MicroServerADMIN      = MicroServer{Name: "admin", ServerType: ServerTypeHttp}
	MicroServerSITE       = MicroServer{Name: "site", ServerType: ServerTypeHttp}
	MicroServersShop      = MicroServer{Name: "shop", ServerType: ServerTypeHttp}
	MicroServerMANAGER    = MicroServer{Name: "manager", ServerType: ServerTypeHttp}
	MicroServerMimiServer = MicroServer{Name: "mimi", ServerType: ServerTypeGrpc}
	MicroServerMiniapp    = MicroServer{Name: "miniapp", ServerType: ServerTypeHttp}
)
