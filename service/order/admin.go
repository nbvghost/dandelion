package order

import (
	"context"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

// CancelOk 后台或者客服调用的接口
func (service OrdersService) CancelOk(context context.Context, OrdersID dao.PrimaryKey, wxConfig *model.WechatConfig) (string, error) {
	Orm := db.Orm()

	//var orders model.Orders
	orders := dao.GetByPrimaryKey(Orm, entity.Orders, OrdersID).(*model.Orders)
	if orders.ID == 0 {
		return "", errors.New("订单不存在")
	}

	//下单状态
	if orders.Status == model.OrdersStatusCancel || orders.Status == model.OrdersStatusPay {

		if orders.IsPay == model.OrdersIsPayPayed {

			var refund *refunddomestic.Refund
			var err error

			err = service.Wx.Refund(context, orders, nil, "用户取消", wxConfig)
			if err != nil {
				return "", err
			}

			if refund != nil {
				switch refund.Status {
				case refunddomestic.STATUS_SUCCESS.Ptr():
					return "退款成功", nil
				case refunddomestic.STATUS_CLOSED.Ptr():
					return "退款关闭", nil
				case refunddomestic.STATUS_PROCESSING.Ptr():
					return "退款处理中", nil
				case refunddomestic.STATUS_ABNORMAL.Ptr():
					return "", errors.New("退款异常")
				}
				return "", errors.New("无效的退款状态")
			}

		}
	}
	return "", errors.New("当前订单状态，无法取消订单")
}
