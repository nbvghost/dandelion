package model

import (
	"github.com/nbvghost/gpa/types"
)

//等级
type Rank struct {
	types.Entity
	GrowMaxValue uint   `gorm:"column:GrowMaxValue"`
	Title        string `gorm:"column:Title"`
	//VoucherID     uint `gorm:"column:VoucherID"`
}

func (Rank) TableName() string {
	return "Rank"
}
