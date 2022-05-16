package model

import (
	"time"

	"github.com/lib/pq"
	"github.com/nbvghost/gpa/types"
)

type Admin struct {
	types.Entity
	OID         types.PrimaryKey `gorm:"column:OID;uniqueIndex:admin_idx_unique_id"`
	Account     string           `gorm:"column:Account;not null;uniqueIndex:admin_idx_unique_id"`
	Phone       string           `gorm:"column:Phone;default:''"` //json 权限
	Name        string           `gorm:"column:Name"`
	PassWord    string           `gorm:"column:PassWord;not null"`
	Authority   string           `gorm:"column:Authority;default:''"` //json 权限
	Roles       pq.StringArray   `gorm:"column:Roles;type:text[]"`    //角色
	LastLoginAt time.Time        `gorm:"column:LastLoginAt"`
	Initiator   bool             `gorm:"column:Initiator"`
}

func (Admin) TableName() string {
	return "Admin"
}
