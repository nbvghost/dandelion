package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// 拼团记录
type CollageRecord struct {
	dao.Entity
	OrderNo       string         `gorm:"column:OrderNo"`
	OrdersGoodsNo string         `gorm:"column:OrdersGoodsNo"`
	No            string         `gorm:"column:No"`
	UserID        dao.PrimaryKey `gorm:"column:UserID"`
	Collager      dao.PrimaryKey `gorm:"column:Collager"`
	//IsPay         uint `gorm:"column:IsPay"` //是否支付成功：0=未支付，1=支付成功
}

func (CollageRecord) TableName() string {
	return "CollageRecord"
}
