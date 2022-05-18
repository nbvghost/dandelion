package model

import (
	"github.com/nbvghost/gpa/types"
)

type LeaveMessage struct {
	types.Entity
	OID      types.PrimaryKey       `gorm:"column:OID"`
	Name     string                 `gorm:"column:Name"`
	Email    string                 `gorm:"column:Email"`
	Content  string                 `gorm:"column:Content"`
	ClientIP string                 `gorm:"column:ClientIP"`
	Extend   map[string]interface{} `gorm:"column:Extend;type:JSON"`
}

func (u LeaveMessage) TableName() string {
	return "LeaveMessage"
}
