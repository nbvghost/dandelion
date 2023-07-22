package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// 等级
type Rank struct {
	dao.Entity
	GrowMaxValue uint   `gorm:"column:GrowMaxValue"`
	Title        string `gorm:"column:Title"`
	//VoucherID     uint `gorm:"column:VoucherID"`
}

func (Rank) TableName() string {
	return "Rank"
}
