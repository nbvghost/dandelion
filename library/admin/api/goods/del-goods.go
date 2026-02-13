package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type DelGoods struct {
	Get struct {
		ID uint `form:"ID"`
	} `method:"get"`
}

func (m *DelGoods) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	//ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
	//ID := object.ParseUint(context.Request.URL.Query().Get("ID"))
	return &result.JsonResult{Data: service.Goods.Goods.DeleteGoods(ctx, dao.PrimaryKey(m.Get.ID))}, err
}
