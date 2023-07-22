package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type Transfers struct {
	dao.Entity
	OrderNo    string         `gorm:"column:OrderNo;unique"` //订单号
	UserID     dao.PrimaryKey `gorm:"column:UserID"`         //
	StoreID    dao.PrimaryKey `gorm:"column:StoreID"`
	Amount     uint           `gorm:"column:Amount"`     //提现金额
	ReUserName string         `gorm:"column:ReUserName"` //提现用户真实的名字。
	Desc       string         `gorm:"column:Desc"`       //提现说明
	IP         string         `gorm:"column:IP"`         //IP
	OpenId     string         `gorm:"column:OpenId"`     //OpenId
	IsPay      uint           `gorm:"column:IsPay"`      //是否支付成功,0=未支付，1，支付成功，2过期
}

func (Transfers) TableName() string {
	return "Transfers"
}
