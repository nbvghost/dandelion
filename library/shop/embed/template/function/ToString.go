package function

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/funcmap"
	"github.com/nbvghost/tool/object"
)

type ToString struct {
	V any `arg:""`
}

func (g *ToString) Call(ctx constrain.IContext) funcmap.IFuncResult {
	return funcmap.NewStringFuncResult(object.ParseString(g.V))
}
