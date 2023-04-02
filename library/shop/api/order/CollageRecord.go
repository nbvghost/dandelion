package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/order"
)

type CollageRecord struct {
	OrdersService order.OrdersService
	User          *model.User `mapping:""`
	Get           struct {
		Index int `form:"Index"`
	} `method:"get"`
}

func (m *CollageRecord) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//Index, _ := strconv.Atoi(context.Request.URL.Query().Get("Index"))
	list := m.OrdersService.ListCollageRecord(m.User.ID, m.Get.Index)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: list}}, nil

}
