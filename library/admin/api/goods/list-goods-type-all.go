package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type ListGoodsTypeAll struct {
	Organization *model.Organization `mapping:""`
}

func (g *ListGoodsTypeAll) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	gts := service.Goods.Goods.GoodsTypeService.ListGoodsTypeByOIDForAdmin(g.Organization.ID)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", gts)}, err
}
