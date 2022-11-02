package model

import (
	"github.com/nbvghost/gpa/types"
)

//门店库存
type StoreStock struct {
	types.Entity
	StoreID         types.PrimaryKey `gorm:"column:StoreID"`
	GoodsID         types.PrimaryKey `gorm:"column:GoodsID"`
	SpecificationID types.PrimaryKey `gorm:"column:SpecificationID"`
	Stock           uint             `gorm:"column:Stock"`
	UseStock        uint             `gorm:"column:UseStock"` //已经使用的量
}

func (u StoreStock) TableName() string {
	return "StoreStock"
}
