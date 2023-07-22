package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type GoodsSkuLabel struct {
	dao.Entity
	OID     dao.PrimaryKey `gorm:"column:OID;index"`
	GoodsID dao.PrimaryKey `gorm:"column:GoodsID;uniqueIndex:UIGoodsIDLabel"`
	Label   string         `gorm:"column:Label;uniqueIndex:UIGoodsIDLabel"`
	Name    string         `gorm:"column:Name"`
	Abel    bool           `gorm:"column:Abel"`
	Image   bool           `gorm:"column:Image"`
}

func (GoodsSkuLabel) TableName() string {
	return "GoodsSkuLabel"
}
