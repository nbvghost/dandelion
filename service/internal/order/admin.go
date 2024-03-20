package order

import (
	"context"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/pkg/errors"
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
			var err error
			err = service.Wx.Refund(context, orders, nil, "用户取消", wxConfig)
			if err != nil {
				return "", err
			}
			return "退款申请成功", nil
		}
	}
	return "", errors.New("当前订单状态，无法取消订单")
}
