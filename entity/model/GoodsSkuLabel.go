package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type GoodsSkuLabel struct {
	dao.Entity
	OID     dao.PrimaryKey `gorm:"column:OID;index"`
	GoodsID dao.PrimaryKey `gorm:"column:GoodsID;uniqueIndex:UIGoodsIDLabel"`
	Label   string         `gorm:"column:Label;uniqueIndex:UIGoodsIDLabel"` //用于显示的名字
	Name    string         `gorm:"column:Name"`                             //字段名
	Abel    bool           `gorm:"column:Abel"`
	Image   bool           `gorm:"column:Image"`
}

func (GoodsSkuLabel) TableName() string {
	return "GoodsSkuLabel"
}
