package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)



type GoodsAttributes struct {
	dao.Entity
	OID     dao.PrimaryKey `gorm:"column:OID;index"` //
	GoodsID dao.PrimaryKey `gorm:"column:GoodsID;index"`
	GroupID dao.PrimaryKey `gorm:"column:GroupID"`
	Name    string         `gorm:"column:Name"`
	Value   string         `gorm:"column:Value"`
}

func (u GoodsAttributes) TableName() string {
	return "GoodsAttributes"
}


