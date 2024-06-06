package function

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/funcmap"
	"github.com/nbvghost/dandelion/library/util"
	"strings"
)

type Url struct {
	Values []string `arg:"..."`
}

func (g *Url) Call(ctx constrain.IContext) funcmap.IFuncResult {
	if len(g.Values)%2 != 0 {
		return funcmap.NewStringFuncResult("参数必须是偶数个")
	}
	var contextValue = contexext.FromContext(ctx)
	var u = util.GetHost(contextValue.Request) + contextValue.Request.URL.Path

	params := contextValue.Request.URL.Query()
	var label string
	for i := range g.Values {
		if (i+1)%2 == 1 {
			label = strings.ToLower(g.Values[i])
		} else {
			params.Set(label, strings.ToLower(g.Values[i]))
		}
	}
	return funcmap.NewStringFuncResult(u + "?" + params.Encode())
}
