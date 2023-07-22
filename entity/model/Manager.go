package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type Manager struct {
	dao.Entity
	Account  string `gorm:"column:Account;not null;unique"`
	PassWord string `gorm:"column:PassWord;not null"`
}

func (Manager) TableName() string {
	return "Manager"
}
