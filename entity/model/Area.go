package model

import "github.com/nbvghost/gpa/types"

//Area 中国全国5级行政区划（省、市、县、镇、村）
//
//code,name,level,pcode
//
//level: 省1，市2，县3，镇4，村5
//
//code: 12位，省2位，市2位，县2位，镇3位，村3位
//
//pcode: 直接父级别的code
type Area struct {
	types.Entity
	Code  uint   `gorm:"column:Code;index"`
	Name  string `gorm:"column:Name"`
	Level uint   `gorm:"column:Level"`
	PCode uint   `gorm:"column:PCode;index"`
}

func (Area) TableName() string {
	return "Area"
}
