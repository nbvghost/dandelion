package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/company"
)

type StockList struct {
	StoreService company.StoreService
	Store        *model.Store `mapping:""`
	Get          struct {
		Order string `uri:"Order"`
	} `method:"Get"`
}

func (m *StockList) Handle(context constrain.IContext) (constrain.IResult, error) {

	list := m.StoreService.ListStoreStock(m.Store.ID)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: list}}, nil

}
