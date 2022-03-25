package model

import "github.com/nbvghost/dandelion/entity/base"

//等级
type Rank struct {
	base.BaseModel
	GrowMaxValue uint   `gorm:"column:GrowMaxValue"`
	Title        string `gorm:"column:Title"`
	//VoucherID     uint `gorm:"column:VoucherID"`
}

func (Rank) TableName() string {
	return "Rank"
}
