package model

import (
	"github.com/nbvghost/gpa/types"
)

type Subscribe struct {
	types.Entity
	OID   types.PrimaryKey `gorm:"column:OID"`
	Email string           `gorm:"column:Email;unique"`
}

func (Subscribe) TableName() string {
	return "Subscribe"
}
