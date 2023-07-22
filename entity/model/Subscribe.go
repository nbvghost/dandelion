package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type Subscribe struct {
	dao.Entity
	OID   dao.PrimaryKey `gorm:"column:OID"`
	Email string         `gorm:"column:Email;unique"`
}

func (Subscribe) TableName() string {
	return "Subscribe"
}
