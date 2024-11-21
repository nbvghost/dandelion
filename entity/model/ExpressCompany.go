package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type ExpressCompany struct {
	dao.Entity
	Key  string `gorm:"column:Key;unique"`
	Name string `gorm:"column:Name"`
}

func (u ExpressCompany) TableName() string {
	return "ExpressCompany"
}
