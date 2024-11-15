package function

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/funcmap"
	"go.uber.org/zap"
	"net/url"
)

type URLPathJoin struct {
	Values []string `arg:"..."` //不定参数只能有一个
}

func (g *URLPathJoin) Call(ctx constrain.IContext) funcmap.IFuncResult {
	path, err := url.JoinPath("/", g.Values...)
	if err != nil {
		ctx.Logger().With(zap.NamedError("URLPathJoin", err))
	}
	return funcmap.NewStringFuncResult(path)
}
