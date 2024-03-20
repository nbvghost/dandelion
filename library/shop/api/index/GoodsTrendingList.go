package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/mode"
)

type GoodsTrendingList struct {
	Get struct {
		Index int `form:"index"`
	} `method:"get"`
	User *model.User `mapping:""`
}

func (m *GoodsTrendingList) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//index, _ := strconv.Atoi(context.Request.URL.Query().Get("index"))
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	params := &mode.ListQueryParam{}
	orderBy := &extends.Order{}

	pagination := service.Goods.Goods.GoodsList(params, m.User.OID, orderBy.OrderByColumn(`"CountSale"+"CountView"`, true), m.Get.Index+1, 10)

	return result.NewData(pagination), nil

	//return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: controller.Goods.HotList()}}
}
