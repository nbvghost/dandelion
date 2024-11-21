package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service"
)

type LayoutType string

const (
	StyleTypeSide LayoutType = "side"
	StyleTypeBar  LayoutType = "bar"
	StyleTypeCard LayoutType = "card"
	StyleTypeForm LayoutType = "form"
)

type TopTypeGoods string

const (
	TopTypeGoodsView TopTypeGoods = "view"
	TopTypeGoodsSale TopTypeGoods = "sale"
)

type TopGoods struct {
	Organization *model.Organization `mapping:""`

	StyleType LayoutType   `arg:""`
	TopType   TopTypeGoods `arg:""`
}

func (m *TopGoods) Template() ([]byte, error) {
	return nil, nil
}

func (m *TopGoods) Render(ctx constrain.IContext) (map[string]any, error) {

	var list []model.Goods
	if m.TopType == TopTypeGoodsView {
		list = service.Goods.Sort.HotViewList(m.Organization.ID, 8)
	} else if m.TopType == TopTypeGoodsSale {
		list = service.Goods.Sort.HotSaleList(m.Organization.ID, 8)
	}

	return map[string]any{
		"List":      list,
		"StyleType": m.StyleType,
		"TopType":   m.TopType,
	}, nil
}
