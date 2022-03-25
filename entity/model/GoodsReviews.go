package model

import (
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/gpa/types"
)

//商品评论
type GoodsReviews struct {
	types.Entity
	GoodsID types.PrimaryKey
	UserID  types.PrimaryKey
	Content string
	Image   sqltype.StringArray
	Rating  uint
	Like    uint
}

func (u *GoodsReviews) TableName() string {
	return "GoodsReviews"
}
