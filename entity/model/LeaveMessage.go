package model

import (
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/gpa/types"
)

type LeaveMessage struct {
	types.Entity
	OID           types.PrimaryKey      `gorm:"column:OID"`
	Name          string                `gorm:"column:Name"`
	Email         string                `gorm:"column:Email"`
	SocialAccount sqltype.SocialAccount `gorm:"column:SocialAccount;type:JSON"`
	Content       string                `gorm:"column:Content"`
}

func (u LeaveMessage) TableName() string {
	return "LeaveMessage"
}
