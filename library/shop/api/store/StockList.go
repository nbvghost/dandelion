package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type StockList struct {
	Store *model.Store `mapping:""`
	Get   struct {
		Order string `uri:"Order"`
	} `method:"Get"`
}

func (m *StockList) Handle(context constrain.IContext) (constrain.IResult, error) {

	list := service.Company.Store.ListStoreStock(m.Store.ID)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: list}}, nil

}
