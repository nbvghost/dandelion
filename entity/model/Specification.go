package model

import (
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
)

// Specification 商品规格
type Specification struct {
	dao.Entity
	OID         dao.PrimaryKey                `gorm:"column:OID;index"`
	GoodsID     dao.PrimaryKey                `gorm:"column:GoodsID;index"`        //
	Name        string                        `gorm:"column:Name"`                 //
	Label       string                        `gorm:"column:Label"`                //
	LabelIndex  sqltype.Array[dao.PrimaryKey] `gorm:"column:LabelIndex;type:JSON"` //
	Num         uint                          `gorm:"column:Num"`                  //件
	Weight      uint                          `gorm:"column:Weight"`               //每件 多少重 g
	Stock       uint                          `gorm:"column:Stock"`                //
	CostPrice   uint                          `gorm:"column:CostPrice"`            //成本价
	MarketPrice uint                          `gorm:"column:MarketPrice"`          //市场价
	Brokerage   uint                          `gorm:"column:Brokerage"`            //总佣金
}

func (m *Specification) GetMarketPrice(quantity uint) uint {
	return m.MarketPrice * quantity
}
func (Specification) TableName() string {
	return "Specification"
}
