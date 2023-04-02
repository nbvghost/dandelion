package model

import (
	"github.com/nbvghost/gpa/types"
)

type GoodsSkuLabelData struct {
	types.Entity
	GoodsSkuLabelID types.PrimaryKey `gorm:"column:GoodsSkuLabelID;uniqueIndex:UIGoodsSkuLabelIDLabel"`
	GoodsID         types.PrimaryKey `gorm:"column:GoodsID"`
	Label           string           `gorm:"column:Label;uniqueIndex:UIGoodsSkuLabelIDLabel"`
	Image           string           `gorm:"column:Image"`
}

func (GoodsSkuLabelData) TableName() string {
	return "GoodsSkuLabelData"
}
