package model

import (
	"github.com/nbvghost/gpa/types"
)

//省市
type District struct {
	types.Entity
	Code string `gorm:"column:Code;primary_key;unique"`
	Name string `gorm:"column:Name"`
}

func (District) TableName() string {
	return "District"
}
