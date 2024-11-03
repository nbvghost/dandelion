package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type Language struct {
	dao.Entity
	Code        string `gorm:"column:Code"`        //en
	Name        string `gorm:"column:Name"`        //English
	SelfName    string `gorm:"column:SelfName"`    //自称
	ChineseName string `gorm:"column:ChineseName"` //中文
}

func (u Language) TableName() string {
	return "Language"
}
