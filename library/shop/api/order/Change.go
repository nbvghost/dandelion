package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/gpa/types"

	"github.com/pkg/errors"
)

type Change struct {
	OrdersService order.OrdersService
	WechatConfig  *model.WechatConfig `mapping:""`
	Put           struct {
		Action        string           `form:"Action"`
		OrdersGoodsID types.PrimaryKey `form:"OrdersGoodsID"`
		ID            types.PrimaryKey `form:"ID"`
		ShipName      string           `form:"ShipName"`
		ShipNo        string           `form:"ShipNo"`
		RefundInfo    string           `form:"RefundInfo"`
	} `method:"put"`
}

func (m *Change) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	//context.Request.ParseForm()
	//Action := context.Request.FormValue("Action")
	switch m.Put.Action {
	case "RefundInfo":
		//OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		//OrdersGoodsID := object.ParseUint(context.Request.FormValue("OrdersGoodsID"))
		//ShipName := context.Request.FormValue("ShipName")
		//ShipNo := context.Request.FormValue("ShipNo")
		err, info := m.OrdersService.RefundInfo(types.PrimaryKey(m.Put.OrdersGoodsID), m.Put.ShipName, m.Put.ShipNo)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, nil
	case "AskRefund":
		//OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		//OrdersGoodsID := object.ParseUint(context.Request.FormValue("OrdersGoodsID"))
		//RefundInfoJson := context.Request.FormValue("RefundInfo")
		var RefundInfo model.RefundInfo
		util.JSONToStruct(m.Put.RefundInfo, &RefundInfo)
		err, info := m.OrdersService.AskRefund(types.PrimaryKey(m.Put.OrdersGoodsID), RefundInfo)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err
	case "TakeDeliver":
		//ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		//ID := object.ParseUint(context.Request.FormValue("ID"))
		err := m.OrdersService.TakeDeliver(types.PrimaryKey(m.Put.ID))
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "确认收货成功", nil)}, err
	case "Cancel":
		//ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		//ID := object.ParseUint(context.Request.FormValue("ID"))
		info, err := m.OrdersService.Cancel(ctx, types.PrimaryKey(m.Put.ID), m.WechatConfig)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err

	}

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无法操作"), "OK", nil)}, nil
}

func (m *Change) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}