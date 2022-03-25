package view

import (
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/gobext"
)

type Index struct {
	extends.ViewBase
	TrendingList extends.GoodsInfoPagination
	SaleList     extends.GoodsInfoPagination
	NewList      extends.GoodsInfoPagination
	StarList     extends.GoodsInfoPagination
}

func init() {
	gobext.Register(Index{})
}
