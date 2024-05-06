package order

import (
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/payment"
	"time"
)

func (m OrdersService) RefundShip(OrdersID dao.PrimaryKey, ShipKey, ShipName, ShipNo string) (error, string) {
	Orm := db.Orm()

	//var ordersGoods model.OrdersGoods
	orders := dao.GetByPrimaryKey(Orm, &model.Orders{}, OrdersID).(*model.Orders)

	if orders.RefundID > 0 {
		err := dao.UpdateByPrimaryKey(Orm, &model.OrdersGoodsRefund{}, orders.RefundID, map[string]interface{}{
			"Status": model.RefundStatusRefundShip,
			"ShipInfo": sqltype.ShipInfo{
				No:   ShipNo,
				Name: ShipName,
				Key:  ShipKey,
			},
		})
		if err != nil {
			return err, ""
		}
	}

	return nil, "快递信息填写成功"
}

// RefundComplete 后台执行的退款
func (m OrdersService) RefundComplete(context constrain.IContext,ordersID dao.PrimaryKey, refundType uint) (string, error) {
	tx := db.Orm().Begin()

	//var ordersGoods model.OrdersGoods
	//ordersGoods := dao.GetByPrimaryKey(tx, entity.OrdersGoods, OrdersGoodsID).(*model.OrdersGoods)

	//var orders model.Orders
	orders := dao.GetByPrimaryKey(tx, entity.Orders, ordersID).(*model.Orders)
	if orders.IsZero() {
		tx.Rollback()
		return "", errors.New("找不到订单数据")
	}

	//ordersPackage := service.GetOrdersPackageByOrderNo(orders.OrdersPackageNo)

	//RefundPrice := int64(ordersGoods.SellPrice) - int64(math.Floor(((float64(ordersGoods.SellPrice)*float64(ordersGoods.Quantity))/float64(orders.GoodsMoney)*float64(orders.DiscountMoney))+0.5))
	/*RefundPrice := ordersGoods.SellPrice * uint(ordersGoods.Quantity)
	if RefundPrice < 0 {
		RefundPrice = 0
	}*/
	//var RefundInfo model.RefundInfo
	//util.JSONToStruct(ordersGoods.RefundInfo, &RefundInfo)
	//RefundInfo.RefundPrice = RefundPrice
	//orders.RefundInfo.Status = sqltype.RefundStatusRefundComplete

	pm:=payment.NewPayment(context,orders.OID,orders.PayMethod)
	err := dao.UpdateByPrimaryKey(tx, &model.OrdersGoodsRefund{}, orders.RefundID, map[string]interface{}{"Status": model.RefundStatusRefundComplete})
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = pm.Refund(orders, nil, "用户申请退款")
	if err != nil {
		tx.Rollback()
		return "", err
	}
	//扣除佣金
	err = m.AfterSettlementUserBrokerage(tx, orders)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	/*ogs, err := service.FindOrdersGoodsByOrdersID(tx, orders.OrdersID)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	haveRefunc := false
	//totalBrokerage := uint(0)
	for i := range ogs {
		value := ogs[i].(*model.OrdersGoods)
		//totalBrokerage = totalBrokerage + (value.TotalBrokerage * uint(value.Quantity))
		if !(value.Status == model.OrdersGoodsStatusOGRefundComplete) && !(value.Status == model.OrdersGoodsStatusOGNone) {
			haveRefunc = true
			break
		}
	}

	if haveRefunc == false {
		//orders 所有的子单品订单，已经全部退款成功。改orders为完成

		//err := dao.UpdateByPrimaryKey(tx, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusOrderOk})
		err = dao.UpdateByPrimaryKey(tx, &model.Orders{}, orders.ID, map[string]interface{}{"Status": model.OrdersStatusRefundOk})
		if err != nil {
			tx.Rollback()
			return "", err
		}
		//扣除佣金
		err = service.AfterSettlementUserBrokerage(tx, orders)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}*/

	tx.Commit()

	//err := dao.UpdateByPrimaryKey(Orm, OrdersGoodsID, &model.OrdersGoods{}, map[string]interface{}{"Status": model.OrdersStatusOGRefundOk})
	return "已经同意,并已退款", nil
}
func (m OrdersService) RefundAgree(OrdersID dao.PrimaryKey) (error, string) {
	Orm := db.Orm()
	orders := dao.GetByPrimaryKey(Orm, entity.Orders, OrdersID).(*model.Orders)
	//orders.RefundInfo.Status = sqltype.RefundStatusRefundAgree
	err := dao.UpdateByPrimaryKey(Orm, &model.OrdersGoodsRefund{}, orders.RefundID, map[string]interface{}{"Status": model.RefundStatusRefundAgree})
	return err, "已经同意"
}
func (m OrdersService) RefundReject(OrdersID dao.PrimaryKey) (error, string) {
	//Orm := db.Orm()
	//err := dao.UpdateByPrimaryKey(Orm, entity.OrdersGoods, OrdersGoodsID, map[string]interface{}{"Status": model.OrdersGoodsStatusOGRefundNo})
	Orm := db.Orm()
	orders := dao.GetByPrimaryKey(Orm, entity.Orders, OrdersID).(*model.Orders)
	//orders.RefundInfo.Status = sqltype.RefundStatusRefundReject
	err := dao.UpdateByPrimaryKey(Orm, &model.OrdersGoodsRefund{}, orders.RefundID, map[string]interface{}{"Status": model.RefundStatusRefundReject})
	return err, "已经拒绝"
}
func (m OrdersService) AskRefund(OrdersID dao.PrimaryKey, OrdersGoodsID dao.PrimaryKey, HasGoods bool, Reason string) (error, string) {
	tx := db.Orm().Begin()

	//var ordersGoods model.OrdersGoods
	//ordersGoods := dao.GetByPrimaryKey(tx, entity.OrdersGoods, OrdersID).(*model.OrdersGoods)

	//var orders model.Orders
	orders := dao.GetByPrimaryKey(tx, entity.Orders, OrdersID).(*model.Orders)
	if orders.ID == 0 {
		tx.Rollback()
		return errors.New("订单不存在"), ""
	}

	//orders.RefundInfo.HasGoods = HasGoods
	//orders.RefundInfo.Reason = Reason
	//orders.RefundInfo.AskTime = time.Now()
	//orders.RefundInfo.Status = sqltype.RefundStatusRefund

	//下单状态,如果订单状态为，已经发货状态或正在退款中,其它状态无须退款
	if (orders.Status == model.OrdersStatusDeliver) || (orders.Status == model.OrdersStatusRefund) {
		var err error

		ordersGoodsRefund := &model.OrdersGoodsRefund{
			OID:           orders.OID,
			OrdersID:      OrdersID,
			OrdersGoodsID: OrdersGoodsID,
			Status:        model.RefundStatusRefund,
			ShipInfo:      sqltype.ShipInfo{},
			HasGoods:      HasGoods,
			Reason:        Reason,
			ApplyAt:       time.Now(),
		}
		err = dao.Create(tx, ordersGoodsRefund)
		if err != nil {
			tx.Rollback()
			return err, ""
		}

		if OrdersGoodsID > 0 {
			err = dao.UpdateByPrimaryKey(tx, &model.OrdersGoods{}, OrdersGoodsID, map[string]any{"RefundID": ordersGoodsRefund.ID})
			if err != nil {
				tx.Rollback()
				return err, ""
			}
		}

		changeOrders := make(map[string]any)
		changeOrders["Status"] = model.OrdersStatusRefund
		changeOrders["RefundID"] = ordersGoodsRefund.ID
		if orders.Status == model.OrdersStatusDeliver {
			err = dao.UpdateByPrimaryKey(tx, entity.Orders, orders.ID, changeOrders)
		} else {
			err = dao.UpdateByPrimaryKey(tx, entity.Orders, orders.ID, changeOrders)
		}
		if err != nil {
			tx.Rollback()
			return err, ""
		}

		tx.Commit()
		return nil, "已经申请，等待商家确认"

		/*err := dao.UpdateByPrimaryKey(tx, entity.OrdersGoods, OrdersGoodsID, map[string]interface{}{"RefundInfo": util.StructToJSON(&RefundInfo), "Status": model.OrdersGoodsStatusOGAskRefund})
		if err != nil {
			tx.Rollback()
			return err, ""
		} else {

		}*/

	}
	tx.Rollback()
	return errors.New("不允许申请退款"), ""
}

/*
{"mchid":"1652384025","out_trade_no":"a83f1d2f1c413d66322f27e7f8a699bf",
"transaction_id":"4200001990202310130267337609","out_refund_no":"a83f1d2f1c413d66322f27e7f8a699bf",
"refund_id":"50301007362023101326018445177","refund_status":"SUCCESS","success_time":"2023-10-13T15:31:43+08:00",
"amount":{"total":3600,"refund":3600,"payer_total":3600,"payer_refund":3600},"user_received_account":"支付用户零钱"}

OutTradeNo:   core.String(order.OrderNo),
OutRefundNo:  core.String(ordersGoods.OrdersGoodsNo),
*/
func (m OrdersService) OrdersRefundSuccess(orders *model.Orders) error {
	if orders.Status == model.OrdersStatusCancelOk {
		//说明已经退款
		return nil
	}
	if orders.Status == model.OrdersStatusClosed {
		//关闭了，不处理
		return nil
	}
	if orders.Status == model.OrdersStatusDelete {
		//删除了，不处理
		return nil
	}
	tx := db.Orm().Begin()

	if orders.Status == model.OrdersStatusRefund {

		//orders.RefundInfo.Status = sqltype.RefundStatusRefundPay
		err := dao.UpdateByPrimaryKey(tx, &model.Orders{}, orders.ID, map[string]interface{}{"Status": model.OrdersStatusClosed})
		if err != nil {
			tx.Rollback()
			return err
		}
		err = dao.UpdateByPrimaryKey(tx, &model.OrdersGoodsRefund{}, orders.RefundID, map[string]interface{}{"Status": model.RefundStatusRefundPay})
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		err := dao.UpdateByPrimaryKey(tx, entity.Orders, orders.ID, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	//管理商品库存
	err := m.OrdersStockManager(tx, orders, false)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = m.AfterSettlementUserBrokerage(tx, orders)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
