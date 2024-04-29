package extends

import "github.com/nbvghost/dandelion/entity/model"

type AnalyseResult struct {
	FullCut         model.FullCut
	OrdersGoodsInfo []*OrdersGoodsMix
	FavouredPrice   uint
	FullCutAll      uint
	GoodsPrice      uint
	ExpressPrice    uint
}
