package model

import (
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/library/dao"
)

type MediaTarget string

const (
	MediaTargetContractOrder MediaTarget = "ContractOrder"
)

type Media struct {
	dao.Entity
	OID      dao.PrimaryKey `gorm:"column:OID;uniqueIndex:otts"`
	TargetID dao.PrimaryKey `gorm:"column:TargetID;uniqueIndex:otts"`
	Target   MediaTarget    `gorm:"column:Target;uniqueIndex:otts"`
	SHA256   string         `gorm:"column:SHA256;uniqueIndex:otts"`
	Title    string         `gorm:"column:Title"`
	Src      string         `gorm:"column:Src"`
	Size     int            `gorm:"column:Size"`
	Width    int            `gorm:"column:Width"`
	Height   int            `gorm:"column:Height"`
	FileName string         `gorm:"column:FileName"`
	Format   string         `gorm:"column:Format"`
	Tags     pq.StringArray `gorm:"column:Tags;type:text[]"`
}

func (Media) TableName() string {
	return "Media"
}
