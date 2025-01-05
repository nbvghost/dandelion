package order

import (
	"errors"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Change struct {

	//OrdersID   uint   `form:"OrdersGoodsID"`
	//PayMoney   float64 `form:"PayMoney"`
	//ShipID     dao.PrimaryKey `form:"ShipID"`
	//ShipNo     string         `form:"ShipNo"`

	Put struct {
		Action     string         `form:"Action"`
		OrdersID   dao.PrimaryKey `form:"OrdersID"`
		RefundType uint           `form:"RefundType"`
		//ShipName   string         `form:"ShipName"`
		ShipNo    string `form:"ShipNo"`
		ShipKey   string `form:"ShipKey"`
		ShipTitle string `form:"ShipTitle"`
		//HasGoods bool           `form:"HasGoods"`
		//Reason   string         `form:"Reason"`
		PayMoney float64 `form:"PayMoney"`
	} `method:"put"`

	WechatConfig *model.WechatConfig `mapping:""`
}

func (m *Change) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Change) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.Orm()

	switch m.Put.Action {
	case "RefundComplete":
		//OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		//RefundType, _ := strconv.ParseUint(context.Request.FormValue("RefundType"), 10, 64)
		info, err := service.Order.Orders.RefundComplete(context, dao.PrimaryKey(m.Put.OrdersID), m.Put.RefundType)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err
	case "RefundAgree":
		//OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		err, info := service.Order.Orders.RefundAgree(dao.PrimaryKey(m.Put.OrdersID))
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err
	case "RefundReject":
		//OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		err, info := service.Order.Orders.RefundReject(dao.PrimaryKey(m.Put.OrdersID))
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err
	case "Cancel":
		//ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		info, err := service.Order.Orders.Cancel(context, dao.PrimaryKey(m.Put.OrdersID))
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err
	case "CancelOk":
		//ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		//RefundType, _ := strconv.ParseUint(context.Request.FormValue("RefundType"), 10, 64) //退款资金来源	 0=未结算资金退款,1=可用余额退款
		info, err := service.Order.Orders.CancelOk(context, dao.PrimaryKey(m.Put.OrdersID))
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err
	case "Deliver":
		err := service.Order.Orders.Deliver(m.Put.ShipTitle, m.Put.ShipKey, m.Put.ShipNo, m.Put.OrdersID)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "发货成功", nil)}, err
	case "PayMoney":
		err := dao.UpdateByPrimaryKey(Orm, entity.Orders, dao.PrimaryKey(m.Put.OrdersID), map[string]interface{}{"PayMoney": uint(m.Put.PayMoney * 100)})
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
		//success, message := m.Orders.ChangeOrdersPayMoney(m.Put.PayMoney, dao.PrimaryKey(m.Put.OrdersID), m.WechatConfig)
		//return &result.JsonResult{Data: &result.ActionResult{Code: success, Message: message, Data: nil}}, err

	}

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("999"), "OK", nil)}, err
}
