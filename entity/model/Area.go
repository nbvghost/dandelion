package model

import "github.com/nbvghost/dandelion/library/dao"

// Area 中国全国5级行政区划（省、市、县、镇、村）
//
// code,name,level,pcode
//
// level: 省1，市2，县3，镇4，村5
//
// code: 12位，省2位，市2位，县2位，镇3位，村3位
//
// pcode: 直接父级别的code
type Area struct {
	Code  dao.PrimaryKey `gorm:"COMMENT:Code;NOT NULL;column:Code;PRIMARY_KEY"`
	Name  string         `gorm:"column:Name"`
	Level uint           `gorm:"column:Level"`
	PCode dao.PrimaryKey `gorm:"column:PCode;index"`
}

func (Area) TableName() string {
	return "Area"
}
