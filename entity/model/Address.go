package model

import (
	"github.com/nbvghost/gpa/types"
	"strings"
)

type Address struct {
	types.Entity
	UserID          types.PrimaryKey `gorm:"column:UserID"`
	Name            string           `gorm:"column:Name"`
	CountyName      string           `gorm:"column:CountyName"`
	ProvinceName    string           `gorm:"column:ProvinceName"`
	CityName        string           `gorm:"column:CityName"`
	Detail          string           `gorm:"column:Detail"`
	PostalCode      string           `gorm:"column:PostalCode"`
	Tel             string           `gorm:"column:Tel"`
	Company         string           `gorm:"column:Company"`
	DefaultBilling  bool             `gorm:"column:DefaultBilling"`
	DefaultShipping bool             `gorm:"column:DefaultShipping"`
}

func (addr Address) TableName() string {
	return "Address"
}

func (addr Address) IsEmpty() bool {

	return strings.EqualFold(addr.Name, "") || strings.EqualFold(addr.Tel, "") || strings.EqualFold(addr.Detail, "")
}
