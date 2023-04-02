package model

import (
	"github.com/nbvghost/gpa/types"
)

type GoodsSkuLabel struct {
	types.Entity
	GoodsID types.PrimaryKey `gorm:"column:GoodsID;uniqueIndex:UIGoodsIDLabel"`
	Label   string           `gorm:"column:Label;uniqueIndex:UIGoodsIDLabel"`
	Name    string           `gorm:"column:Name"`
	Abel    bool             `gorm:"column:Abel"`
	Image   bool             `gorm:"column:Image"`
}

func (GoodsSkuLabel) TableName() string {
	return "GoodsSkuLabel"
}
