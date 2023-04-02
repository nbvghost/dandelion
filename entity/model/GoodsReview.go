package model

import (
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/gpa/types"
)

// GoodsReview helpful
type GoodsReview struct {
	types.Entity
	GoodsID          types.PrimaryKey    `gorm:"column:GoodsID"`
	UserID           types.PrimaryKey    `gorm:"column:UserID"`
	Title            string              `gorm:"column:Title"`
	Content          string              `gorm:"column:Content"`
	Image            sqltype.StringArray `gorm:"column:Image"`
	Portrait         string              `gorm:"column:Portrait"`
	NickName         string              `gorm:"column:NickName"`
	Helpful          uint                `gorm:"column:Helpful"`
	IsBuy            bool                `gorm:"column:IsBuy"`
	BuySpecification string              `gorm:"column:BuySpecification"`
	Rating           uint                `gorm:"column:Rating"`
	Like             uint                `gorm:"column:Like"`
}

func (GoodsReview) TableName() string {
	return "GoodsReview"
}
