package model

import "github.com/nbvghost/gpa/types"

type GoodsWish struct {
	types.Entity
	GoodsID         types.PrimaryKey
	SpecificationID types.PrimaryKey
}

func (u GoodsWish) TableName() string {
	return "GoodsWish"
}
