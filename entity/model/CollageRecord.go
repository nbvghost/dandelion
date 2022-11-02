package model

import (
	"github.com/nbvghost/gpa/types"
)

//拼团记录
type CollageRecord struct {
	types.Entity
	OrderNo       string           `gorm:"column:OrderNo"`
	OrdersGoodsNo string           `gorm:"column:OrdersGoodsNo"`
	No            string           `gorm:"column:No"`
	UserID        types.PrimaryKey `gorm:"column:UserID"`
	Collager      types.PrimaryKey `gorm:"column:Collager"`
	//IsPay         uint `gorm:"column:IsPay"` //是否支付成功：0=未支付，1=支付成功
}

func (CollageRecord) TableName() string {
	return "CollageRecord"
}
