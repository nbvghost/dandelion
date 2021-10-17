package iservice

import (
	"github.com/nbvghost/gweb"
)

type IRoute interface {
	RegisterRoute(path string, handler gweb.IHandler)
}
