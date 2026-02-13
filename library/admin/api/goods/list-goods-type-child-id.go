package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type ListGoodsTypeChildID struct {
	Get struct {
		ID uint `form:"ID"`
	} `method:"get"`
}

func (g *ListGoodsTypeChildID) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	//ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
	//ID := object.ParseUint(context.Request.URL.Query().Get("ID"))
	gts := service.Goods.GoodsType.ListAllGoodsTypeChild(ctx, dao.PrimaryKey(g.Get.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", gts)}, err
}
