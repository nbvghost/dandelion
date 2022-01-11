package model

import "github.com/nbvghost/gpa/types"

type GoodsAttributesGroup struct {
	types.Entity
	GoodsID types.PrimaryKey `gorm:"column:GoodsID;index"`
	Name    string           `gorm:"column:Name"`
}

func (u GoodsAttributesGroup) TableName() string {
	return "GoodsAttributesGroup"
}
