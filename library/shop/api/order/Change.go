package order

import (
	"context"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"

	"github.com/pkg/errors"
)

type Change struct {
	OrdersService order.OrdersService
	WechatConfig  *model.WechatConfig `mapping:""`
	Put           struct {
		Action        string         `form:"Action"`
		OrdersGoodsID dao.PrimaryKey `form:"OrdersGoodsID"`
		ID            dao.PrimaryKey `form:"ID"`
		ShipName      string         `form:"ShipName"`
		ShipNo        string         `form:"ShipNo"`
		RefundInfo    string         `form:"RefundInfo"`
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
		err, info := m.OrdersService.RefundInfo(dao.PrimaryKey(m.Put.OrdersGoodsID), m.Put.ShipName, m.Put.ShipNo)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, nil
	case "AskRefund":
		//OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		//OrdersGoodsID := object.ParseUint(context.Request.FormValue("OrdersGoodsID"))
		//RefundInfoJson := context.Request.FormValue("RefundInfo")
		var RefundInfo model.RefundInfo
		util.JSONToStruct(m.Put.RefundInfo, &RefundInfo)
		err, info := m.OrdersService.AskRefund(dao.PrimaryKey(m.Put.OrdersGoodsID), RefundInfo)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err
	case "TakeDeliver":
		//ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		//ID := object.ParseUint(context.Request.FormValue("ID"))
		err := m.OrdersService.TakeDeliver(dao.PrimaryKey(m.Put.ID))
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "确认收货成功", nil)}, err
	case "Cancel":
		//ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		//ID := object.ParseUint(context.Request.FormValue("ID"))
		info, err := m.Cancel(ctx, dao.PrimaryKey(m.Put.ID), m.WechatConfig)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, info, nil)}, err

	}

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无法操作"), "OK", nil)}, nil
}

func (m *Change) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}

// Cancel 申请取消
func (m *Change) Cancel(ctx context.Context, OrdersID dao.PrimaryKey, wxConfig *model.WechatConfig) (string, error) {
	Orm := db.Orm()

	//var orders model.Orders
	orders := dao.GetByPrimaryKey(Orm, entity.Orders, OrdersID).(*model.Orders)
	if orders.ID == 0 {

		return "", errors.New("订单不存在")
	}

	//ordersPackage := service.GetOrdersPackageByOrderNo(orders.OrdersPackageNo)

	//下单状态
	if orders.Status == model.OrdersStatusOrder {
		if orders.IsPay == model.OrdersIsPayPayed {
			err := dao.UpdateByPrimaryKey(Orm, entity.Orders, OrdersID, map[string]interface{}{"Status": model.OrdersStatusCancel})
			return "申请取消，等待客服确认", err
		} else {
			/*transaction, err := service.Wx.OrderQuery(ctx, orders.OrderNo, wxConfig)
			if err != nil {
				return "", err
			}
			if strings.EqualFold(*transaction.TradeState, "SUCCESS") {
				//如果查询订单已经支付，由客服确认
				err := dao.UpdateByPrimaryKey(Orm, entity.Orders, OrdersID, map[string]interface{}{"Status": model.OrdersStatusCancel})
				return "申请取消，等待客服确认", err
			} else*/
			{
				//没支付的订单
				//管理商品库存
				err := m.OrdersService.OrdersStockManager(Orm, orders, false)
				if err != nil {
					return "", err
				}
				err = dao.UpdateByPrimaryKey(Orm, entity.Orders, OrdersID, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
				return "取消成功", err
				/*refund, err := service.Wx.Refund(ctx, orders, ordersPackage, orders.PayMoney, "用户取消", wxConfig)
				if err != nil {
					return "", err
				}
				log.Println("Orders", "Cancel", refund)
				if Success == false {
					Success, Message1 = service.Wx.Refund(ctx, orders, ordersPackage, orders.PayMoney, "用户取消", wxConfig)
					log.Println("Orders", "Cancel", Message1)
				}

				if Success {
					//管理商品库存
					err := service.OrdersStockManager(Orm, orders, false)
					if err != nil {
						return "", err
					}
					err = dao.UpdateByPrimaryKey(Orm, entity.Orders, OrdersID, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
					return "取消成功", err
				} else {
					//管理商品库存
					err := service.OrdersStockManager(Orm, orders, false)
					if err != nil {
						return "", err
					}
					err = dao.UpdateByPrimaryKey(Orm, entity.Orders, OrdersID, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
					return "取消成功", err


					//return errors.New(Message1), ""
					/*Success, Message2 := service.Wx.Refund(orders, orders.PayMoney, orders.PayMoney, "用户取消", 1)
					if Success {

					} else {

					}*/
			}
		}
	} else if orders.Status == model.OrdersStatusPay {
		if orders.IsPay == model.OrdersIsPayPayed {
			//已经支付的订单，发起退款
			//ordersPackage := service.GetOrdersPackageByOrderNo(orders.OrdersPackageNo)
			refund, err := m.OrdersService.Wx.Refund(ctx, orders, nil, "用户取消", wxConfig)
			if err != nil {
				return "", err
			}
			if refund.Status != refunddomestic.STATUS_SUCCESS.Ptr() {
				return "", errors.New("退款异常")
			}
			err = dao.UpdateByPrimaryKey(Orm, &model.Orders{}, OrdersID, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
			if err != nil {
				return "", err
			}
			err = m.OrdersService.OrdersStockManager(Orm, orders, false)
			if err != nil {
				return "", err
			}
			return "订单已经取消，退款资金已经按原路退回，请注意查收信息", nil

		} else {
			return "", errors.New("不允许取消订单,订单没有支付或已经过期")
		}

	} else {
		return "", errors.New("不允许取消订单")
	}
}
