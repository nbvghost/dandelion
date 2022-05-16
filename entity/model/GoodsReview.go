package model

import "github.com/nbvghost/dandelion/entity/base"

// GoodsReview helpful
type GoodsReview struct {
	base.BaseModel
	GoodsID          uint
	Content          string
	Portrait         string
	NickName         string
	Helpful          uint
	BuySpecification string
	Star             uint
}

func (GoodsReview) TableName() string {
	return "GoodsReview"
}
