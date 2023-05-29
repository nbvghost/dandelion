package route

import (
	"github.com/nbvghost/dandelion/constrain"
)

type RegisterRoute func(path string, handler constrain.IHandler)
type RegisterView func(path string, handler constrain.IViewHandler)
