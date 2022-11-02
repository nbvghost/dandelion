package model

import (
	"github.com/nbvghost/gpa/types"
)

type Manager struct {
	types.Entity
	Account  string `gorm:"column:Account;not null;unique"`
	PassWord string `gorm:"column:PassWord;not null"`
}

func (Manager) TableName() string {
	return "Manager"
}
