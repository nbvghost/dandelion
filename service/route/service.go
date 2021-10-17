package route

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/service/iservice"
	"github.com/nbvghost/gweb"
	"reflect"
)

type service struct {
	Routes map[string]reflect.Type
}

func (m *service) RegisterRoute(path string, handler gweb.IHandler) {
	if _, ok := m.Routes[path]; ok {
		panic(errors.New(fmt.Sprintf("存在相同的路由:%s", path)))
	}
	m.Routes[path] = reflect.TypeOf(handler).Elem()
}

func New() iservice.IRoute {
	return &service{Routes: map[string]reflect.Type{}}
}
