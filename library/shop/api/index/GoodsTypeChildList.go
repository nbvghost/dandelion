package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type GoodsTypeChildList struct {
	Get struct {
		GoodsTypeID dao.PrimaryKey `uri:"GoodsTypeID"`
	} `method:"get"`
}

func (m *GoodsTypeChildList) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//GoodsTypeID, _ := strconv.ParseUint(context.PathParams["GoodsTypeID"], 10, 64)
	//GoodsTypeID := object.ParseUint(context.PathParams["GoodsTypeID"])
	results := service.Goods.GoodsType.ListGoodsTypeChild(ctx, dao.PrimaryKey(m.Get.GoodsTypeID))
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}, nil
}
