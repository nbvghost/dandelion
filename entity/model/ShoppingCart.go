package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// ShoppingCart 购物车
type ShoppingCart struct {
	dao.Entity
	UserID          dao.PrimaryKey `gorm:"column:UserID"`                  //
	Goods           string         `gorm:"column:Goods;type:text"`         //
	Specification   string         `gorm:"column:Specification;type:text"` //
	GoodsID         dao.PrimaryKey `gorm:"column:GoodsID"`                 //
	SpecificationID dao.PrimaryKey `gorm:"column:SpecificationID"`         //
	Quantity        uint           `gorm:"column:Quantity"`                //数量
}

func (ShoppingCart) TableName() string {
	return "ShoppingCart"
}
