package extends

import (
	"github.com/nbvghost/dandelion/entity/model"
)

type GoodsInfo struct {
	Goods          model.Goods
	Specifications []*model.Specification
	Attributes     []GoodsAttributes
	Discounts      []Discount
	GoodsType      model.GoodsType
	GoodsTypeChild model.GoodsTypeChild
	SkuLabels      []SkuLabel
}

type SkuLabel struct {
	Label *model.GoodsSkuLabel
	Data  []*model.GoodsSkuLabelData
}
