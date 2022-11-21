package model

import (
	"github.com/nbvghost/gpa/types"
)

type ExpressCompany struct {
	types.Entity
	Key  string `gorm:"column:Key;unique"`
	Name string `gorm:"column:Name"`
}

func (u ExpressCompany) TableName() string {
	return "ExpressCompany"
}
