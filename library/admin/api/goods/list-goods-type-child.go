package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ListGoodsTypeChild struct {
}

func (g *ListGoodsTypeChild) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	//var gts []model.GoodsTypeChild
	//err = g.Goods.FindAll(db.GetDB(ctx), &gts)
	gts := dao.Find(db.GetDB(ctx), entity.GoodsTypeChild).List()
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", gts)}, err
}
