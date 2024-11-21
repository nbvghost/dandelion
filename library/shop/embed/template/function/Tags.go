package function

import (
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/tag"
	"github.com/nbvghost/dandelion/library/funcmap"
)

type Tags struct {
	Tags pq.StringArray `arg:""`
}

func (g *Tags) Call(ctx constrain.IContext) funcmap.IFuncResult {
	return funcmap.NewResult(tag.ToTagsUri(g.Tags))
}
