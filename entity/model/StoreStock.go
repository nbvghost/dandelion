package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// 门店库存
type StoreStock struct {
	dao.Entity
	StoreID         dao.PrimaryKey `gorm:"column:StoreID"`
	GoodsID         dao.PrimaryKey `gorm:"column:GoodsID"`
	SpecificationID dao.PrimaryKey `gorm:"column:SpecificationID"`
	Stock           uint           `gorm:"column:Stock"`
	UseStock        uint           `gorm:"column:UseStock"` //已经使用的量
}

func (u StoreStock) TableName() string {
	return "StoreStock"
}
