package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/tool/object"
)

// 卡
type CardItem struct {
	dao.Entity
	OrderNo       string         `gorm:"column:OrderNo;unique"` //订单号
	UserID        dao.PrimaryKey `gorm:"column:UserID"`         //
	Type          string         `gorm:"column:Type"`           //OrdersGoods,Voucher,ScoreGoods
	OrdersGoodsID dao.PrimaryKey `gorm:"column:OrdersGoodsID"`  //
	VoucherID     dao.PrimaryKey `gorm:"column:VoucherID"`      //
	ScoreGoodsID  dao.PrimaryKey `gorm:"column:ScoreGoodsID"`   //
	Data          string         `gorm:"column:Data;type:text"` //json数据
	Quantity      uint           `gorm:"column:Quantity"`       //数量
	UseQuantity   uint           `gorm:"column:UseQuantity"`    //已经使用数量
	ExpireTime    time.Time      `gorm:"column:ExpireTime"`     //过期时间
	PostType      OrdersPostType `gorm:"column:PostType"`       //1=邮寄，2=线下使用
}

func (cardItem CardItem) GetNameLabel(DB *gorm.DB) (Name, Label string) {

	switch cardItem.Type {
	case "OrdersGoods":
		var item OrdersGoods
		DB.First(&item, cardItem.OrdersGoodsID)
	/*	var goods Goods
		var specification Specification
		util.JSONToStruct(item.Goods, &goods)
		util.JSONToStruct(item.Specification, &specification)*/
		Name = item.Goods.Title
		Label = "规格：" + item.Specification.Label + "(" + strconv.FormatFloat(float64(item.Specification.Num)*float64(item.Specification.Weight)/1000, 'f', 2, 64) + "Kg)"
	case "Voucher":
		var item Voucher
		DB.First(&item, cardItem.VoucherID)
		Name = item.Name
		Label = "金额：" + strconv.FormatFloat(float64(item.Amount)/100, 'f', 2, 64) + "元," + "说明：" + item.Introduce

	case "ScoreGoods":
		var item ScoreGoods
		DB.First(&item, cardItem.ScoreGoodsID)
		Name = item.Name
		Label = "积分：" + object.ParseString(item.Score) + "," + "市场价：" + strconv.FormatFloat(float64(item.Price)/100, 'f', 2, 64) + "元," + "说明：" + item.Introduce
	}

	return Name, Label
}
func (CardItem) TableName() string {
	return "CardItem"
}
