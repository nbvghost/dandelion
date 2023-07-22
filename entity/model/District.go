package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// 省市
type District struct {
	dao.Entity
	Code string `gorm:"column:Code;primary_key;unique"`
	Name string `gorm:"column:Name"`
}

func (District) TableName() string {
	return "District"
}
