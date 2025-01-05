package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type MediaTarget string

const (
	MediaTargetContractOrder MediaTarget = "ContractOrder"
)

type Media struct {
	dao.Entity
	OID      dao.PrimaryKey `gorm:"column:OID;index"`
	TargetID dao.PrimaryKey `gorm:"column:TargetID"`
	Target   MediaTarget    `gorm:"column:Target"`
	Title    string         `gorm:"column:Title"`
	Src      string         `gorm:"column:Src"`
	Size     int            `gorm:"column:Size"`
	Width    int            `gorm:"column:Width"`
	Height   int            `gorm:"column:Height"`
	FileName string         `gorm:"column:FileName"`
	Format   string         `gorm:"column:Format"`
	SHA256   string         `gorm:"column:SHA256;unique"`
}

func (Media) TableName() string {
	return "Media"
}
