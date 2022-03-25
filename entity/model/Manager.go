package model

import "github.com/nbvghost/dandelion/entity/base"

type Manager struct {
	base.BaseModel
	Account  string `gorm:"column:Account;not null;unique"`
	PassWord string `gorm:"column:PassWord;not null"`
}

func (Manager) TableName() string {
	return "Manager"
}
