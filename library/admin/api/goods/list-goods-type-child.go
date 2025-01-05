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

func (g *ListGoodsTypeChild) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//var gts []model.GoodsTypeChild
	//err = g.Goods.FindAll(db.Orm(), &gts)
	gts := dao.Find(db.Orm(), entity.GoodsTypeChild).List()
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", gts)}, err
}
