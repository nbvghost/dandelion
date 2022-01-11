package extends

import (
	"github.com/nbvghost/dandelion/entity/model"
)

type GoodsInfo struct {
	Goods          model.Goods
	Specifications []model.Specification
	Attributes     []*model.GoodsAttributes
	Discounts      []Discount
	GoodsType      model.GoodsType
	GoodsTypeChild model.GoodsTypeChild
}
