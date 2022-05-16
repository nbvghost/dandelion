package model

import (
	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//限时抢购商品
type TimeSellGoods struct {
	base.BaseModel
	OID          types.PrimaryKey `gorm:"column:OID"`
	TimeSellHash string           `gorm:"column:TimeSellHash"`
	GoodsID      types.PrimaryKey `gorm:"column:GoodsID"`
	Disable      bool             `gorm:"column:Disable"` //限时抢购中，单个商品是暂时
}

func (TimeSellGoods) TableName() string {
	return "TimeSellGoods"
}
