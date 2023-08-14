package model

import (
	"strings"

	"github.com/nbvghost/dandelion/library/dao"
)

type OrdersGoodsStatus string

const (
	OrdersGoodsStatusOGNone           OrdersGoodsStatus = ""                 // none
	OrdersGoodsStatusOGAskRefund      OrdersGoodsStatus = "OGAskRefund"      // OGAskRefund=申请，申请退货退款
	OrdersGoodsStatusOGRefundNo       OrdersGoodsStatus = "OGRefundNo"       // OGRefundOK=拒绝子商品，确认退货款
	OrdersGoodsStatusOGRefundOk       OrdersGoodsStatus = "OGRefundOk"       // OGRefundOK=允许子商品，确认退货款
	OrdersGoodsStatusOGRefundInfo     OrdersGoodsStatus = "OGRefundInfo"     // OGRefundInfo=用户填写信息，允许子商品，确认退货款
	OrdersGoodsStatusOGRefundComplete OrdersGoodsStatus = "OGRefundComplete" // OGRefund=完成子商品，用户邮寄商品，商家待收货
)

type OrdersGoods struct {
	dao.Entity
	OID            dao.PrimaryKey    `gorm:"column:OID"`
	OrdersGoodsNo  string            `gorm:"column:OrdersGoodsNo;unique"`    //
	Status         OrdersGoodsStatus `gorm:"column:Status"`                  //OGAskRefund，OGRefundNo，OGRefundOk，OGRefundInfo，OGRefundComplete
	RefundInfo     string            `gorm:"column:RefundInfo;type:text"`    //RefundInfo json 退款退货信息
	OrdersID       dao.PrimaryKey    `gorm:"column:OrdersID"`                //
	Goods          string            `gorm:"column:Goods;type:text"`         //josn
	Specification  string            `gorm:"column:Specification;type:text"` //json
	GoodsSkus      string            `gorm:"column:GoodsSkus;type:JSON"`     //json
	Discounts      string            `gorm:"column:Discounts;type:JSON"`
	Quantity       uint              `gorm:"column:Quantity"`       //数量
	CostPrice      uint              `gorm:"column:CostPrice"`      //单价-原价
	SellPrice      uint              `gorm:"column:SellPrice"`      //单价-销售价
	TotalBrokerage uint              `gorm:"column:TotalBrokerage"` //总佣金
	Error          string            `gorm:"column:Error"`          //
	//SpecificationID uint `gorm:"column:SpecificationID"`             //
	//CollageNo     string `gorm:"column:CollageNo"` //拼团码，每个订单都是唯一
	//TimeSellID     uint `gorm:"column:TimeSellID"`             //限时抢购ID
	//TimeSell       string `gorm:"column:TimeSell;type:text"` //json
	//GoodsID         uint `gorm:"column:GoodsID"`                     //
}

func (og OrdersGoods) AddError(err string) {

	if strings.EqualFold(og.Error, "") {
		og.Error = err
	} else {
		og.Error = og.Error + "|" + err
	}
}
func (OrdersGoods) TableName() string {
	return "OrdersGoods"
}
