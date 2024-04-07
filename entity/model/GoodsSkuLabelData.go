package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type GoodsSkuLabelData struct {
	dao.Entity
	OID             dao.PrimaryKey `gorm:"column:OID;index"`
	GoodsSkuLabelID dao.PrimaryKey `gorm:"column:GoodsSkuLabelID;uniqueIndex:UIGoodsSkuLabelIDLabel"`
	GoodsID         dao.PrimaryKey `gorm:"column:GoodsID;index"`
	Label           string         `gorm:"column:Label;uniqueIndex:UIGoodsSkuLabelIDLabel"`
	Image           string         `gorm:"column:Image"`
}

func (GoodsSkuLabelData) TableName() string {
	return "GoodsSkuLabelData"
}
