package order

import (
	"encoding/json"
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
)

type InfoOrders struct {
	User *model.User `mapping:""`
	Get  struct {
		OrderNo string `form:"order-no"`
	} `method:"get"`
	Put struct {
		OrderNo   string
		AddressID dao.PrimaryKey
	} `method:"Put"`
}

func (m *InfoOrders) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	orders := repository.OrdersDao.GetOrdersByOrderNo(m.Get.OrderNo)
	var address model.Address
	if err := json.Unmarshal([]byte(orders.Address), &address); err != nil {
		return nil, err
	}
	confirmOrdersGoods, err := service.Order.Orders.AnalyseOrdersGoodsListByOrders(&orders, &address)
	return result.NewData(confirmOrdersGoods), err
}
func (m *InfoOrders) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	address := dao.GetByPrimaryKey(db.Orm(), &model.Address{}, m.Put.AddressID).(*model.Address)
	if address.ID == 0 {
		return nil, errors.New("地址不能为空")
	}
	if len(m.Put.OrderNo) == 0 {
		return nil, errors.New("the parameter is invalid")
	}
	orders := repository.OrdersDao.GetOrdersByOrderNo(m.Put.OrderNo)
	if orders.ID == 0 {
		return nil, errors.New("order data does not exist")
	}

	confirmOrdersGoods, err := service.Order.Orders.AnalyseOrdersGoodsListByOrders(&orders, address)
	if err != nil {
		return nil, err
	}

	changeData := make(map[string]any)
	changeData["Address"] = util.StructToJSON(address)
	changeData["ExpressMoney"] = confirmOrdersGoods.ExpressPrice
	err = dao.UpdateByPrimaryKey(db.Orm(), &model.Orders{}, orders.ID, changeData)
	return result.NewData(confirmOrdersGoods), err
}
