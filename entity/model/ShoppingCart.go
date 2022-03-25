package model

import (
	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//购物车
type ShoppingCart struct {
	base.BaseModel
	UserID        types.PrimaryKey `gorm:"column:UserID"`
	GSID          string           `gorm:"column:GSID"` //GoodsID+""+SpecificationID
	Goods         string           `gorm:"column:Goods;type:text"`
	Specification string           `gorm:"column:Specification;type:text"`
	Quantity      uint             `gorm:"column:Quantity"` //数量
}

func (ShoppingCart) TableName() string {
	return "ShoppingCart"
}
