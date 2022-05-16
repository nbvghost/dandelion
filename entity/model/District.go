package model

import "github.com/nbvghost/dandelion/entity/base"

//省市
type District struct {
	base.BaseModel
	Code string `gorm:"column:Code;primary_key;unique"`
	Name string `gorm:"column:Name"`
}

func (District) TableName() string {
	return "District"
}
