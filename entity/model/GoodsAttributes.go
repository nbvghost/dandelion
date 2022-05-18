package model

import "github.com/nbvghost/gpa/types"

type GoodsAttributes struct {
	types.Entity
	GoodsID types.PrimaryKey `gorm:"column:GoodsID;index"`
	GroupID types.PrimaryKey `gorm:"column:GroupID"`
	Name    string           `gorm:"column:Name"`
	Value   string           `gorm:"column:Value"`
}

func (u GoodsAttributes) TableName() string {
	return "GoodsAttributes"
}
