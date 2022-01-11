package view

import (
	"github.com/nbvghost/dandelion/domain/pagination"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/gobext"
)

type Search struct {
	extends.ViewBase
	ViewType    string
	GoodsType   []extends.AllGoodsType
	GoodsList0  []extends.GoodsInfo
	GoodsList1  []extends.GoodsInfo
	GoodsTypeID uint
	Keyword     string
	SwitchPage  pagination.SwitchPage
}

func init() {
	gobext.Register(Search{})
}
