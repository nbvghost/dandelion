package extends

import "github.com/nbvghost/dandelion/entity/model"

type AnalyseResult struct {
	FullCut         model.FullCut
	OrdersGoodsInfo []*model.OrdersGoods
	FavouredPrice   uint
	FullCutAll      uint
	GoodsPrice      uint
	ExpressPrice    uint
}
