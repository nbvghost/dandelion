package model

import (
	"github.com/nbvghost/gpa/types"
)

//拼团商品
type CollageGoods struct {
	types.Entity
	OID         types.PrimaryKey `gorm:"column:OID"`
	CollageHash string           `gorm:"column:CollageHash"`
	GoodsID     types.PrimaryKey `gorm:"column:GoodsID"`
	Disable     bool             `gorm:"column:Disable"` //限时抢购中，单个商品是暂时
}

func (CollageGoods) TableName() string {
	return "CollageGoods"
}
