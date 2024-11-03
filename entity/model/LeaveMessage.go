package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type LeaveMessage struct {
	dao.Entity
	OID      dao.PrimaryKey `gorm:"column:OID"`
	Name     string         `gorm:"column:Name"`
	Email    string         `gorm:"column:Email"`
	Content  string         `gorm:"column:Content"`
	ClientIP string         `gorm:"column:ClientIP"`
	Extend   string         `gorm:"column:Extend;type:JSON"`
}

func (u LeaveMessage) TableName() string {
	return "LeaveMessage"
}
