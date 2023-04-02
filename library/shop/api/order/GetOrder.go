package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/gpa/types"
)

type GetOrder struct {
	OrdersService order.OrdersService
	Get           struct {
		ID types.PrimaryKey `uri:"ID"`
	} `method:"get"`
}

func (m *GetOrder) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])

	pack := struct {
		Orders          model.Orders
		OrdersGoodsList []model.OrdersGoods
		CollageUsers    []model.User
	}{}
	pack.Orders = m.OrdersService.GetOrdersByID(types.PrimaryKey(m.Get.ID))

	pack.OrdersGoodsList, _ = m.OrdersService.FindOrdersGoodsByOrdersID(singleton.Orm(), pack.Orders.ID)

	//:todo ----
	//og := pack.OrdersGoodsList[0]
	//pack.CollageUsers = controller.Orders.FindOrdersGoodsByCollageUser(og.CollageNo)
	//SELECT u.* FROM Orders o,OrdersGoods og,USER u WHERE og.CollageNo='9d262ef3926bc83f41258410239ce5ba' AND o.ID=og.OrdersID AND u.ID=o.UserID;

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: pack}}, nil

}
