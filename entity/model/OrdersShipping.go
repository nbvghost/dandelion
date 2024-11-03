package model

import "github.com/nbvghost/dandelion/library/dao"

type OrdersShipping struct {
	dao.Entity
	OID     dao.PrimaryKey `gorm:"column:OID"`
	OrderNo string         `gorm:"column:OrderNo;index"` //订单号
	Title   string         `gorm:"column:Title"`
	Image   string         `gorm:"column:Image"` //快递拍照图片
	No      string         `gorm:"column:No"`    //快递单号
	Name    string         `gorm:"column:Name"`  //快递
	Key     string         `gorm:"column:Key"`   //快递编号
}

func (OrdersShipping) TableName() string {
	return "OrdersShipping"
}
