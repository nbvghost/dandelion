package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/goods"
	"github.com/nbvghost/gpa/types"
)

type Goods struct {
	GoodsService goods.GoodsService
	Get          struct {
		ID types.PrimaryKey `uri:"ID"`
	} `method:"get"`
}

func (m *Goods) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	Orm := db.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	goodsInfo, err := m.GoodsService.GetGoods(Orm, ctx, types.PrimaryKey(m.Get.ID))

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: goodsInfo}}, err
}
