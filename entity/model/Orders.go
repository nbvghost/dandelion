package model

import (
	"errors"
	"runtime/debug"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/gpa/types"
)

//订单信息
type Orders struct {
	base.BaseModel
	OID             types.PrimaryKey       `gorm:"column:OID"`             //
	UserID          types.PrimaryKey       `gorm:"column:UserID"`          //用户ID
	PrepayID        string                 `gorm:"column:PrepayID"`        //
	IsPay           sqltype.OrdersIsPay    `gorm:"column:IsPay"`           //是否支付成功,0=未支付，1，支付成功，2过期
	OrderNo         string                 `gorm:"column:OrderNo;unique"`  //订单号
	OrdersPackageNo string                 `gorm:"column:OrdersPackageNo"` //订单号
	PayMoney        uint                   `gorm:"column:PayMoney"`        //支付价
	PostType        sqltype.OrdersPostType `gorm:"column:PostType"`        //1=邮寄，2=线下使用
	Status          string                 `gorm:"column:Status"`          //状态
	ShipNo          string                 `gorm:"column:ShipNo"`          //快递单号
	ShipName        string                 `gorm:"column:ShipName"`        //快递
	Address         string                 `gorm:"column:Address"`         //收货地址 json
	DeliverTime     time.Time              `gorm:"column:DeliverTime"`     //发货时间
	ReceiptTime     time.Time              `gorm:"column:ReceiptTime"`     //确认收货时间
	RefundTime      time.Time              `gorm:"column:RefundTime"`      //申请退款退货时间
	PayTime         time.Time              `gorm:"column:PayTime"`         //支付时间
	DiscountMoney   uint                   `gorm:"column:DiscountMoney"`   //相关活动的折扣金额，目前只有满减。
	GoodsMoney      uint                   `gorm:"column:GoodsMoney"`      //商品总价
	ExpressMoney    uint                   `gorm:"column:ExpressMoney"`    //运费
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