package service

import (
	"dandelion/app/service/dao"

	"github.com/jinzhu/gorm"
)

type OrdersGoodsService struct {
	dao.BaseDao
}

func (service OrdersGoodsService) FindByOrdersID(DB *gorm.DB, OrdersID uint64) ([]dao.OrdersGoods, error) {
	var ogs []dao.OrdersGoods
	err := service.FindWhere(DB, &ogs, &dao.OrdersGoods{OrdersID: OrdersID})
	return ogs, err
}
