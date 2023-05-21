package model

import "github.com/nbvghost/gpa/types"

type GoodsWish struct {
	types.Entity
	OID             types.PrimaryKey `gorm:"column:OID"`
	UserID          types.PrimaryKey `gorm:"column:UserID"`
	GoodsID         types.PrimaryKey `gorm:"column:GoodsID"`
	SpecificationID types.PrimaryKey `gorm:"column:SpecificationID"`
	Quantity        uint             `gorm:"column:Quantity"`
	Comment         string           `gorm:"column:Comment"`
}

func (u GoodsWish) TableName() string {
	return "GoodsWish"
}
