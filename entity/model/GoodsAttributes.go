package model

import "github.com/nbvghost/gpa/types"

type GoodsAttributes struct {
	types.Entity
	GoodsID   types.PrimaryKey `gorm:"column:GoodsID;index"`
	GroupName string           `gorm:"column:GroupName"`
	Name      string           `gorm:"column:Name"`
	Value     string           `gorm:"column:Value"`
}

func (u GoodsAttributes) TableName() string {
	return "GoodsAttributes"
}
