package model

import "github.com/nbvghost/gpa/types"

type DNSType string

const (
	DNSTypeA DNSType = "A"
)

type DNS struct {
	types.Entity
	OID    types.PrimaryKey `gorm:"column:OID;not null;uniqueIndex:uniqueIndexTypeDomainOID"`
	Type   DNSType          `gorm:"column:Type;not null;uniqueIndex:uniqueIndexTypeDomainOID"`
	Domain string           `gorm:"column:Domain;not null;uniqueIndex:uniqueIndexTypeDomainOID"`
}

func (DNS) TableName() string {
	return "DNS"
}
