package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type Language struct {
	dao.Entity
	Name        string `gorm:"column:Name"` //语系
	ChineseName string `gorm:"column:ChineseName"`
	ISOName     string `gorm:"column:ISOName"`
	SelfName    string `gorm:"column:SelfName"`
	CodeBiadu   string `gorm:"column:CodeBiadu"`
	Code6391    string `gorm:"column:Code6391"`
	Code6392T   string `gorm:"column:Code6392T"`
	Code6392B   string `gorm:"column:Code6392B"`
	Code6393    string `gorm:"column:Code6393"`
	Hot         uint   `gorm:"column:Hot"`
}

func (u Language) TableName() string {
	return "Language"
}
