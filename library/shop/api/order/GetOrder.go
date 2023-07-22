package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/order"
)

type GetOrder struct {
	OrdersService order.OrdersService
	Get           struct {
		ID dao.PrimaryKey `uri:"ID"`
	} `method:"get"`
}

func (m *GetOrder) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	pack := struct {
		Orders          model.Orders
		OrdersGoodsList []*extends.OrdersGoods
		CollageUsers    []model.User
	}{}
	pack.Orders = m.OrdersService.GetOrdersByID(m.Get.ID)

	ordersGoodsList, err := m.OrdersService.FindOrdersGoodsByOrdersID(db.Orm(), pack.Orders.ID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(ordersGoodsList); i++ {
		goods, err := m.OrdersService.ConvertOrdersGoods(ordersGoodsList[i].(*model.OrdersGoods))
		if err != nil {
			return nil, err
		}
		pack.OrdersGoodsList = append(pack.OrdersGoodsList, goods)
	}

	//:todo ----
	//og := pack.OrdersGoodsList[0]
	//pack.CollageUsers = controller.Orders.FindOrdersGoodsByCollageUser(og.CollageNo)
	//SELECT u.* FROM Orders o,OrdersGoods og,USER u WHERE og.CollageNo='9d262ef3926bc83f41258410239ce5ba' AND o.ID=og.OrdersID AND u.ID=o.UserID;

	return result.NewData(pack), nil

}
