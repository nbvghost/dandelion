package store

import (
	"errors"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type stockParam struct {
	StoreID            uint `form:"StoreID"`
	ID                 uint `form:"ID"`
	GoodsID            uint `form:"GoodsID"`
	SpecificationID    uint `form:"SpecificationID"`
	AddStoreStockStock int  `form:"AddStoreStockStock"`
}
type Stock struct {
	POST stockParam `method:"POST"`
	PUT  stockParam `method:"PUT"`
}

func (m *Stock) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	return m.updateStock(context, m.PUT)
}

func (m *Stock) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Stock) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	return m.updateStock(context, m.POST)
}
func (m *Stock) updateStock(ctx constrain.IContext, params stockParam) (r constrain.IResult, err error) {
	Orm := db.GetDB(ctx)
	//context.Request.ParseForm()

	//form.StoreID=parseInt($routeParams.ID);
	//form.GoodsID=$scope.SelectGoods.ID;
	//form.ID=$scope.StoreStock.ID;
	//form.SpecificationID=$scope.StoreStock.SpecificationID;
	//form.AddStoreStockStock=$scope.AddStoreStockStock;

	//StoreID := object.ParseUint(context.Request.FormValue("StoreID"))
	//GoodsID := object.ParseUint(context.Request.FormValue("GoodsID"))
	//ID := object.ParseUint(context.Request.FormValue("ID"))
	//SpecificationID := object.ParseUint(context.Request.FormValue("SpecificationID"))
	//AddStoreStockStock := object.ParseInt(context.Request.FormValue("AddStoreStockStock"))

	//item := &model.StoreStock{}
	item := dao.GetByPrimaryKey(Orm, entity.StoreStock, dao.PrimaryKey(params.ID)).(*model.StoreStock)
	if item.ID == 0 {

		if params.AddStoreStockStock < 0 {
			//return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("增加的库存不能小于0"), "", nil)}
			params.AddStoreStockStock = 0
		}

		item.StoreID = dao.PrimaryKey(params.StoreID)
		item.GoodsID = dao.PrimaryKey(params.GoodsID)
		item.SpecificationID = dao.PrimaryKey(params.SpecificationID)
		item.Stock = uint(params.AddStoreStockStock)
		err := dao.Create(Orm, item)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}, nil
	} else {

		stock := int(item.Stock) + params.AddStoreStockStock
		if stock < 0 {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("库存不能为负的"), "", nil)}, nil
		}
		if stock < int(item.UseStock) {
			//return (&result.ActionResult{}).SmartError(errors.New("库存不能为负的"), "", 0)
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("库存不能为负的"), "", nil)}, nil
		}

		item.Stock = uint(stock)
		item.SpecificationID = dao.PrimaryKey(params.SpecificationID)

		err := dao.UpdateByPrimaryKey(Orm, entity.StoreStock, dao.PrimaryKey(params.ID), map[string]interface{}{"SpecificationID": item.SpecificationID, "Stock": item.Stock})
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, nil
	}
}
