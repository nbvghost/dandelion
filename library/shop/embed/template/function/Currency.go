package function

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/funcmap"
	"strconv"
)

type Currency struct {
	Price uint `arg:""` //不定参数只能有一个
}

func (g *Currency) Call(ctx constrain.IContext) funcmap.IFuncResult {

	return funcmap.NewStringFuncResult("$" + strconv.FormatFloat(float64(g.Price)/100, 'f', 2, 64))
}
