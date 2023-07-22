package model

import "github.com/nbvghost/dandelion/library/dao"

type DNSType string

const (
	DNSTypeA DNSType = "A"
)

type DNS struct {
	dao.Entity
	OID    dao.PrimaryKey `gorm:"column:OID;not null;uniqueIndex:uniqueIndexTypeDomainOID"`
	Type   DNSType        `gorm:"column:Type;not null;uniqueIndex:uniqueIndexTypeDomainOID"`
	Domain string         `gorm:"column:Domain;not null;uniqueIndex:uniqueIndexTypeDomainOID"`
}

func (DNS) TableName() string {
	return "DNS"
}
