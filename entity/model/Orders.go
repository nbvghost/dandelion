package model

import (
	"errors"
	"runtime/debug"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/library/dao"
)

type OrdersStatus string

func (m OrdersStatus) CanUpdateOrder() bool {
	switch m {
	case OrdersStatusPending:
		fallthrough
	case OrdersStatusConfirm:
		fallthrough
	case OrdersStatusOrder:
		fallthrough
	case OrdersStatusPay:
		return true
	}
	return false
}

const (
	OrdersStatusPending  OrdersStatus = "Pending"  // 待定
	OrdersStatusConfirm  OrdersStatus = "Confirm"  // 订单确认
	OrdersStatusOrder    OrdersStatus = "Order"    // 下单成功，待付款
	OrdersStatusPay      OrdersStatus = "Pay"      // 支付成功，待发货
	OrdersStatusDeliver  OrdersStatus = "Deliver"  // 发货成功，待收货
	OrdersStatusRefund   OrdersStatus = "Refund"   // 订单退款退货中->所有子商品状态为空或OGRefundOK->返回Deliver状态
	OrdersStatusOrderOk  OrdersStatus = "OrderOk"  // 订单确认完成
	OrdersStatusCancel   OrdersStatus = "Cancel"   // 订单等待取消
	OrdersStatusCancelOk OrdersStatus = "CancelOk" // 订单已经取消
	OrdersStatusDelete   OrdersStatus = "Delete"   // 删除
	OrdersStatusClosed   OrdersStatus = "Closed"   // 已经关闭
)

type OrdersIsPay uint

const (
	OrdersIsPayUnPay  OrdersIsPay = 0 //未支付
	OrdersIsPayPayed  OrdersIsPay = 1 //支付成功
	OrdersIsPayExpire OrdersIsPay = 2 //过期
)

// //是否支付成功,0=未支付，1，支付成功，2过期
type OrdersPostType int

const (
	OrdersPostTypePost    OrdersPostType = 1 //邮寄
	OrdersPostTypeOffline OrdersPostType = 2 //线下使用
)

type OrdersPayMethod string

const (
	OrdersPayMethodWechat  OrdersPayMethod = "wechatpay" //邮寄
	OrdersPayMethodPaypal  OrdersPayMethod = "paypal"    //paypal
	OrdersPayMethodOffline OrdersPayMethod = "offline"   //线下转帐
)

type OrdersType string

const (
	OrdersTypeGoods        OrdersType = "Goods"        //商品购买订单
	OrdersTypeGoodsPackage OrdersType = "GoodsPackage" //合并下单
	OrdersTypeSupply       OrdersType = "Supply"       //充值
)

// Orders 订单信息
type Orders struct {
	dao.Entity
	OID             dao.PrimaryKey  `gorm:"column:OID"`             //
	UserID          dao.PrimaryKey  `gorm:"column:UserID"`          //用户ID
	PrepayID        string          `gorm:"column:PrepayID"`        //微信预支付单号
	TransactionID   string          `gorm:"column:TransactionID"`   //微信交易号，只有支付成功后才有
	IsPay           OrdersIsPay     `gorm:"column:IsPay"`           //是否支付成功,0=未支付，1，支付成功，2过期
	OrderNo         string          `gorm:"column:OrderNo;unique"`  //订单号
	OrdersPackageNo string          `gorm:"column:OrdersPackageNo"` //订单号
	PayMoney        uint            `gorm:"column:PayMoney"`        //支付价
	PostType        OrdersPostType  `gorm:"column:PostType"`        //Deprecated: 1=邮寄，2=线下使用,
	Type            OrdersType      `gorm:"column:Type"`            //
	PayMethod       OrdersPayMethod `gorm:"column:PayMethod"`       //支付方式
	Status          OrdersStatus    `gorm:"column:Status"`          //状态
	Address         string          `gorm:"column:Address"`         //收货地址 json
	DeliverTime     time.Time       `gorm:"column:DeliverTime"`     //发货时间
	ReceiptTime     time.Time       `gorm:"column:ReceiptTime"`     //确认收货时间
	RefundID        dao.PrimaryKey  `gorm:"column:RefundID"`        //申请退款退货时间
	PayTime         time.Time       `gorm:"column:PayTime"`         //支付时间
	DiscountMoney   uint            `gorm:"column:DiscountMoney"`   //相关活动的折扣金额，目前只有满减。
	GoodsMoney      uint            `gorm:"column:GoodsMoney"`      //商品总价
	ExpressMoney    uint            `gorm:"column:ExpressMoney"`    //运费
	//ShipInfo        sqltype.ShipInfo `gorm:"column:ShipInfo;serializer:json;type:json"` //快递 see:OrdersShipping
	//ShipNo          string          `gorm:"column:ShipNo"`          //快递单号
	//ShipName        string          `gorm:"column:ShipName"`        //快递
}

func (u *Orders) BeforeCreate(scope *gorm.DB) (err error) {

	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))

	}
	return nil
}

func (Orders) TableName() string {
	return "Orders"
}
