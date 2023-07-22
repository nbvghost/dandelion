package model

import "github.com/nbvghost/dandelion/library/dao"

type GoodsAttributesGroup struct {
	dao.Entity
	OID     dao.PrimaryKey `gorm:"column:OID;index"` //
	GoodsID dao.PrimaryKey `gorm:"column:GoodsID;index"`
	Name    string         `gorm:"column:Name"`
}

func (u GoodsAttributesGroup) TableName() string {
	return "GoodsAttributesGroup"
}
