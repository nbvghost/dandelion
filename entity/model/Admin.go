package model

import (
	"github.com/nbvghost/gpa/types"
	"time"
)

type Admin struct {
	types.Entity
	OID         types.PrimaryKey `gorm:"column:OID;uniqueIndex:admin_idx_unique_id"`
	Account     string           `gorm:"column:Account;not null;uniqueIndex:admin_idx_unique_id"`
	PassWord    string           `gorm:"column:PassWord;not null"`
	Authority   string           `gorm:"column:Authority;default:''"` //json 权限
	LastLoginAt time.Time        `gorm:"column:LastLoginAt"`
}

func (Admin) TableName() string {
	return "Admin"
}
