package route

import (
	"github.com/nbvghost/dandelion/constrain"
)

type RegisterRoute func(path string, handler constrain.IHandler, withoutAuth ...bool)
type RegisterView func(path string, handler constrain.IViewHandler, result constrain.IViewResult, withoutAuth ...bool)
