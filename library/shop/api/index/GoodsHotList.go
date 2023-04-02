package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/goods"
)

type GoodsHotList struct {
	GoodsService goods.GoodsService
	Get          struct {
		Index int `form:"index"`
	} `method:"get"`
	User *model.User `mapping:""`
}

func (m *GoodsHotList) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//index, _ := strconv.Atoi(context.Request.URL.Query().Get("index"))
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: m.GoodsService.GoodsList(`"CountSale" desc`, m.Get.Index, 10, `"Hide"=? and "OID"=?`, 0, m.User.OID)}}, nil

	//return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: controller.Goods.HotList()}}
}
