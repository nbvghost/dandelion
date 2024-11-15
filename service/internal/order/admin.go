package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/payment"
	"github.com/pkg/errors"
)

// CancelOk 后台或者客服调用的接口
// func (m OrdersService) CancelOk(context context.Context, OrdersID dao.PrimaryKey, wxConfig *model.WechatConfig) (string, error) {
func (m OrdersService) CancelOk(context constrain.IServiceContext, OrdersID dao.PrimaryKey) (string, error) {
	Orm := db.Orm()

	//var orders model.Orders
	orders := dao.GetByPrimaryKey(Orm, entity.Orders, OrdersID).(*model.Orders)
	if orders.ID == 0 {
		return "", errors.New("订单不存在")
	}

	pm := payment.NewPayment(context,orders.OID, orders.PayMethod)

	//下单状态
	if orders.Status == model.OrdersStatusCancel || orders.Status == model.OrdersStatusPay {

		if orders.IsPay == model.OrdersIsPayPayed {
			var err error
			err = pm.Refund(orders, nil, "用户取消")
			if err != nil {
				return "", err
			}
			return "退款申请成功", nil
		}
	}
	return "", errors.New("当前订单状态，无法取消订单")
}
