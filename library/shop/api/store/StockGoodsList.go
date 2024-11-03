package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type StockGoodsList struct {
	User  *model.User  `mapping:""`
	Store *model.Store `mapping:""`
	Get   struct {
		GoodsID dao.PrimaryKey `uri:"GoodsID"`
	} `method:"Get"`
}

func (m *StockGoodsList) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//GoodsID
	//GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)

	list := service.Company.Store.ListStoreSpecifications(m.Store.ID, m.Get.GoodsID)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: list}}, nil
}
