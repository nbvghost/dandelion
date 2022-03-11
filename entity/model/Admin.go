package model

import (
	"github.com/nbvghost/gpa/types"
	"time"
)

type Admin struct {
	types.Entity
	OID         types.PrimaryKey `gorm:"column:OID"`
	Account     string           `gorm:"column:Account;not null;unique"`
	PassWord    string           `gorm:"column:PassWord;not null"`
	Authority   string           `gorm:"column:Authority;default:''"` //json 权限
	LastLoginAt time.Time        `gorm:"column:LastLoginAt"`
}

func (Admin) TableName() string {
	return "Admin"
}
