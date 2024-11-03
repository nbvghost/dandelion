package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
)

type GetOrder struct {
	Get struct {
		ID dao.PrimaryKey `uri:"ID"`
	} `method:"get"`
}

func (m *GetOrder) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	pack := extends.OrdersDetail{}

	pack.Orders = repository.OrdersDao.GetOrdersByID(m.Get.ID)

	ordersGoodsList, err := service.Order.Orders.FindOrdersGoodsByOrdersID(db.Orm(), pack.Orders.ID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(ordersGoodsList); i++ {
		goods, err := service.Order.Orders.ConvertOrdersGoods(ordersGoodsList[i])
		if err != nil {
			return nil, err
		}
		pack.OrdersGoodsList = append(pack.OrdersGoodsList, goods)
	}

	ordersShippingList := make([]*model.OrdersShipping, 0)
	dao.Find(db.Orm(), &model.OrdersShipping{}).Where(`"OrderNo"=?`, pack.Orders.OrderNo).Result(&ordersShippingList)
	pack.OrdersShippingList = ordersShippingList

	//:todo ----
	//og := pack.OrdersGoodsList[0]
	//pack.CollageUsers = controller.Orders.FindOrdersGoodsByCollageUser(og.CollageNo)
	//SELECT u.* FROM Orders o,OrdersGoods og,USER u WHERE og.CollageNo='9d262ef3926bc83f41258410239ce5ba' AND o.ID=og.OrdersID AND u.ID=o.UserID;

	return result.NewData(pack), nil

}
