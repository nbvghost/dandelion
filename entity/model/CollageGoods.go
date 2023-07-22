package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// 拼团商品
type CollageGoods struct {
	dao.Entity
	OID         dao.PrimaryKey `gorm:"column:OID"`
	CollageHash string         `gorm:"column:CollageHash"`
	GoodsID     dao.PrimaryKey `gorm:"column:GoodsID"`
	Disable     bool           `gorm:"column:Disable"` //限时抢购中，单个商品是暂时
}

func (CollageGoods) TableName() string {
	return "CollageGoods"
}
