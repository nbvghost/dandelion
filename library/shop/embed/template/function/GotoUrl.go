package function

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/funcmap"
	"github.com/nbvghost/dandelion/repository"
)

type GotoUrl struct {
	Organization *model.Organization   `mapping:""`
	Type         model.ContentTypeType `arg:""`
	TemplateName string                `arg:""`
}

func (g *GotoUrl) Call(ctx constrain.IContext) funcmap.IFuncResult {

	contentItem := repository.ContentItemDao.GetContentItemByTypeTemplateName(db.Orm(), g.Organization.ID, g.Type, g.TemplateName)
	return funcmap.NewStringFuncResult(fmt.Sprintf("/page/%s", contentItem.Uri))
}
