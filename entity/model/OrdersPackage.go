package model

import (
	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//合并支付
type OrdersPackage struct {
	base.BaseModel
	OrderNo string `gorm:"column:OrderNo;unique"` //订单号
	//OrderList     string `gorm:"column:OrderList;type:text"` //json []
	TotalPayMoney uint             `gorm:"column:TotalPayMoney"` //支付价
	IsPay         uint             `gorm:"column:IsPay"`         //是否支付成功,0=未支付，1，支付成功，2过期
	PrepayID      string           `gorm:"column:PrepayID"`      //
	UserID        types.PrimaryKey `gorm:"column:UserID"`        //用户ID
}

func (OrdersPackage) TableName() string {
	return "OrdersPackage"
}