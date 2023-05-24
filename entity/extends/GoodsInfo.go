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
	Rating         GoodsRating
}

type SkuLabel struct {
	Label *model.GoodsSkuLabel
	Data  []*model.GoodsSkuLabelData
}
type GoodsRating struct {
	Rating      uint `gorm:"column:Rating"`
	RatingCount uint `gorm:"column:RatingCount"`
}

func (m *GoodsRating) Value() uint {
	if m.Rating == 0 || m.RatingCount == 0 {
		return 0
	}
	v := m.Rating / m.RatingCount
	if v > 5 {
		v = 5
	}
	return v
}
