package payment

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service/internal/payment/paypal"
	"github.com/nbvghost/dandelion/service/internal/payment/wechatpay"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

var _ IPayment = &wechatpay.Service{}
var _ IPayment = &paypal.Service{}

type IPayment interface {
	Order(OrderNo string, title, description string, detail, openid string, IP string, Money uint, attach string) (*serviceargument.OrderResult, error)
	OrderQuery(OrderNo string) (*serviceargument.OrderQuery, error)
	CloseOrder(OrderNo string) error
	Refund(order *model.Orders, ordersGoods *model.OrdersGoods, reason string) error
}

func NewPayment(ctx constrain.IContext, oid dao.PrimaryKey, payMethod model.OrdersPayMethod) IPayment {
	switch payMethod {
	case model.OrdersPayMethodWechat:
		return NewWechat(ctx, oid)
	case model.OrdersPayMethodPaypal:
		return NewPaypal(ctx, oid)
	}
	return NewWechat(ctx, oid)
}

func NewWechat(ctx constrain.IContext, oid dao.PrimaryKey) *wechatpay.Service {
	return &wechatpay.Service{Context: ctx, OID: oid}
}

func NewPaypal(ctx constrain.IContext, oid dao.PrimaryKey) *paypal.Service {
	return &paypal.Service{Context: ctx, OID: oid}
}