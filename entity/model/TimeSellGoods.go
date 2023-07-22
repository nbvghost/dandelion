package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// 限时抢购商品
type TimeSellGoods struct {
	dao.Entity
	OID          dao.PrimaryKey `gorm:"column:OID"`
	TimeSellHash string         `gorm:"column:TimeSellHash"`
	GoodsID      dao.PrimaryKey `gorm:"column:GoodsID"`
	Disable      bool           `gorm:"column:Disable"` //限时抢购中，单个商品是暂时
}

func (TimeSellGoods) TableName() string {
	return "TimeSellGoods"
}
