package internal

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type OrdersDao struct{}

func (m OrdersDao) GetOrdersByPrepayID(prepayID string) model.Orders {
	Orm := db.Orm()
	var orders model.Orders
	Orm.Where(map[string]any{"PrepayID": prepayID}).First(&orders)
	return orders
}

func (m OrdersDao) GetOrdersByOrderNo(OrderNo string) model.Orders {
	Orm := db.Orm()
	var orders model.Orders
	Orm.Where(map[string]any{"OrderNo": OrderNo}).First(&orders)
	return orders
}
func (m OrdersDao) GetOrdersByOrdersPackageNo(OrdersPackageNo string) []model.Orders {
	Orm := db.Orm()
	var orders []model.Orders
	Orm.Where(&model.Orders{OrdersPackageNo: OrdersPackageNo}).Find(&orders)
	return orders
}
func (m OrdersDao) GetOrdersByID(ID dao.PrimaryKey) model.Orders {
	Orm := db.Orm()
	var orders model.Orders
	Orm.First(&orders, ID)
	return orders
}
