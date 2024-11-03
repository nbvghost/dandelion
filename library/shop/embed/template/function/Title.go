package function

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/funcmap"
	"github.com/nbvghost/dandelion/service"
)

type Title struct {
	Organization *model.Organization `mapping:""`
}

func (g *Title) Call(ctx constrain.IContext) funcmap.IFuncResult {
	title := service.Content.GetTitle(db.Orm(), g.Organization.ID)

	return funcmap.NewStringFuncResult(title)
}
