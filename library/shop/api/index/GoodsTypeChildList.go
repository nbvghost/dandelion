package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/goods"
	"github.com/nbvghost/gpa/types"
)

type GoodsTypeChildList struct {
	GoodsService goods.GoodsService
	Get          struct {
		GoodsTypeID types.PrimaryKey `uri:"GoodsTypeID"`
	} `method:"get"`
}

func (m *GoodsTypeChildList) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//GoodsTypeID, _ := strconv.ParseUint(context.PathParams["GoodsTypeID"], 10, 64)
	//GoodsTypeID := object.ParseUint(context.PathParams["GoodsTypeID"])
	results := m.GoodsService.ListGoodsTypeChild(types.PrimaryKey(m.Get.GoodsTypeID))
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}, nil
}
