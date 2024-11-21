package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type GoodsAllList struct {
}

func (m *GoodsAllList) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: service.Goods.Goods.AllList()}}, nil
}
