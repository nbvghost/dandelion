package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type Language struct {
	dao.Entity
	Code        string `gorm:"column:Code"` //en
	Name        string `gorm:"column:Name"` //English
}

func (u Language) TableName() string {
	return "Language"
}
