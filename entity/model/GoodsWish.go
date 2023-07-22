package model

import "github.com/nbvghost/dandelion/library/dao"

type GoodsWish struct {
	dao.Entity
	OID             dao.PrimaryKey `gorm:"column:OID"`
	UserID          dao.PrimaryKey `gorm:"column:UserID"`
	GoodsID         dao.PrimaryKey `gorm:"column:GoodsID"`
	SpecificationID dao.PrimaryKey `gorm:"column:SpecificationID"`
	Quantity        uint           `gorm:"column:Quantity"`
	Comment         string         `gorm:"column:Comment"`
}

func (u GoodsWish) TableName() string {
	return "GoodsWish"
}
