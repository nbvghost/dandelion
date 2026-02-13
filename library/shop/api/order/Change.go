package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"github.com/pkg/errors"
)

type Change struct {
	//WechatConfig *model.WechatConfig `mapping:""`
	Put struct {
		Action        string         `form:"Action"`
		OrdersID      dao.PrimaryKey `form:"OrdersID"`
		OrdersGoodsID dao.PrimaryKey `form:"OrdersGoodsID"`
		ID            dao.PrimaryKey `form:"ID"`
		ShipName      string         `form:"ShipName"`
		ShipNo        string         `form:"ShipNo"`
		ShipKey       string         `form:"ShipKey"`
		HasGoods      bool           `form:"HasGoods"`
		Reason        string         `form:"Reason"`
	} `method:"put"`
}

func (m *Change) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	//context.Request.ParseForm()
	//Action := context.Request.FormValue("Action")
	switch m.Put.Action {
	case "RefundShip":
		//OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		//OrdersGoodsID := object.ParseUint(context.Request.FormValue("OrdersGoodsID"))
		//ShipName := context.Request.FormValue("ShipName")
		//ShipNo := context.Request.FormValue("ShipNo")
		err, info := service.Order.Orders.RefundShip(ctx, dao.PrimaryKey(m.Put.OrdersID), m.Put.ShipKey, m.Put.ShipName, m.Put.ShipNo)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, nil
	case "AskRefund":
		//OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		//OrdersGoodsID := object.ParseUint(context.Request.FormValue("OrdersGoodsID"))
		//RefundInfoJson := context.Request.FormValue("RefundInfo")
		//var refundInfo sqltype.RefundInfo //{"HasGoods":true,"Reason":"dsfdsfds fdsfad"}
		//util.JSONToStruct(m.Put.RefundInfo, &refundInfo)
		err, info := service.Order.Orders.AskRefund(ctx, m.Put.OrdersID, m.Put.OrdersGoodsID, m.Put.HasGoods, m.Put.Reason)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err
	case "TakeDeliver":
		//ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		//ID := object.ParseUint(context.Request.FormValue("ID"))
		err := service.Order.Orders.TakeDeliver(ctx, dao.PrimaryKey(m.Put.ID))
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "确认收货成功", nil)}, err
	case "Cancel":
		//ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		//ID := object.ParseUint(context.Request.FormValue("ID"))
		info, err := service.Order.Orders.Cancel(ctx, dao.PrimaryKey(m.Put.ID))
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err

	}

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无法操作"), "OK", nil)}, nil
}

func (m *Change) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
