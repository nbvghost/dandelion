package internal

import (
	"context"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type OrdersDao struct{}

func (m OrdersDao) GetOrdersByPrepayID(ctx context.Context, prepayID string) model.Orders {
	Orm := db.GetDB(ctx)
	var orders model.Orders
	Orm.Where(map[string]any{"PrepayID": prepayID}).First(&orders)
	return orders
}

func (m OrdersDao) GetOrdersByOrderNo(ctx context.Context, OrderNo string) model.Orders {
	Orm := db.GetDB(ctx)
	var orders model.Orders
	Orm.Where(map[string]any{"OrderNo": OrderNo}).First(&orders)
	return orders
}
func (m OrdersDao) GetOrdersByOrdersPackageNo(ctx context.Context, OrdersPackageNo string) []model.Orders {
	Orm := db.GetDB(ctx)
	var orders []model.Orders
	Orm.Where(&model.Orders{OrdersPackageNo: OrdersPackageNo}).Find(&orders)
	return orders
}
func (m OrdersDao) GetOrdersByID(ctx context.Context, ID dao.PrimaryKey) model.Orders {
	Orm := db.GetDB(ctx)
	var orders model.Orders
	Orm.First(&orders, ID)
	return orders
}
