package model

import (
	"github.com/nbvghost/gpa/types"
)

// Specification 商品规格
type Specification struct {
	types.Entity
	GoodsID     types.PrimaryKey `gorm:"column:GoodsID"`
	Label       string           `gorm:"column:Label"`
	Size        string           `gorm:"column:Size"`
	Num         uint             `gorm:"column:Num"`    //件
	Weight      uint             `gorm:"column:Weight"` //每件 多少重 g
	Stock       uint             `gorm:"column:Stock"`
	CostPrice   uint             `gorm:"column:CostPrice"`   //成本价
	MarketPrice uint             `gorm:"column:MarketPrice"` //市场价
	Brokerage   uint             `gorm:"column:Brokerage"`   //总佣金
}

func (Specification) TableName() string {
	return "Specification"
}
