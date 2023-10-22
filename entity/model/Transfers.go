package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type Transfers struct {
	dao.Entity
	OID        dao.PrimaryKey `gorm:"column:OID"`
	OrderNo    string         `gorm:"column:OrderNo;unique"` //订单号
	UserID     dao.PrimaryKey `gorm:"column:UserID"`         //
	StoreID    dao.PrimaryKey `gorm:"column:StoreID"`
	Amount     uint           `gorm:"column:Amount"`     //提现金额
	ReUserName string         `gorm:"column:ReUserName"` //提现用户真实的名字。
	Desc       string         `gorm:"column:Desc"`       //提现说明
	IP         string         `gorm:"column:IP"`         //IP
	OpenId     string         `gorm:"column:OpenId"`     //OpenId
	IsPay      uint           `gorm:"column:IsPay"`      //是否支付成功,0=未支付，1，支付成功,2:关闭
	/*
		WAIT_PAY: 待付款确认。需要付款出资商户在商家助手小程序或服务商助手小程序进行付款确认
		ACCEPTED:已受理。批次已受理成功，若发起批量转账的30分钟后，转账批次单仍处于该状态，可能原因是商户账户余额不足等。商户可查询账户资金流水，若该笔转账批次单的扣款已经发生，则表示批次已经进入转账中，请再次查单确认
		PROCESSING:转账中。已开始处理批次内的转账明细单
		FINISHED:已完成。批次内的所有转账明细单都已处理完成
		CLOSED:已关闭。可查询具体的批次关闭原因确认
	*/
	Status string `gorm:"column:Status"` //同接口中的batch_status
}

func (Transfers) TableName() string {
	return "Transfers"
}
