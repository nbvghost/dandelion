package funcmap

import (
	"github.com/nbvghost/dandelion/constrain"
)

type IWidget interface {
	Render(ctx constrain.IContext) (map[string]interface{}, error)
	Template() ([]byte,error)
}
