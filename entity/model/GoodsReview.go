package model

import (
	"github.com/nbvghost/gpa/types"
)

// GoodsReview helpful
type GoodsReview struct {
	types.Entity
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
