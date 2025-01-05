package model

import (
	"github.com/nbvghost/dandelion/entity/sqltype"
	"strings"
	"time"

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
	//重新定义
	OrdersGoodsStatusProduction OrdersGoodsStatus = "PRODUCTION" //
	OrdersGoodsStatusComplete   OrdersGoodsStatus = "COMPLETE"   //
)

type RefundStatus string

const (
	RefundStatusRefund         RefundStatus = "Refund"         // 申请退款退货，等待客服确认
	RefundStatusRefundAgree    RefundStatus = "RefundAgree"    // 同意退货
	RefundStatusRefundReject   RefundStatus = "RefundReject"   // 拒绝退货
	RefundStatusRefundShip     RefundStatus = "RefundShip"     // 退货货品邮寄中
	RefundStatusRefundComplete RefundStatus = "RefundComplete" // 收到退货货品，客服放款
	RefundStatusRefundPay      RefundStatus = "RefundPay"      // 买家已经收到退款
)

type OrdersGoodsRefund struct {
	dao.Entity
	OID           dao.PrimaryKey `gorm:"column:OID"`
	OrdersID      dao.PrimaryKey `gorm:"column:OrdersID"`
	OrdersGoodsID dao.PrimaryKey `gorm:"column:OrdersGoodsID"`
	Status        RefundStatus   `gorm:"column:Status"`
	ShipInfo      OrdersShipping `gorm:"column:ShipInfo;serializer:json;type:json"`
	HasGoods      bool           `gorm:"column:HasGoods"`
	Reason        string         `gorm:"column:Reason"` //原因
	ApplyAt       time.Time      `gorm:"column:ApplyAt"`
}

func (OrdersGoodsRefund) TableName() string {
	return "OrdersGoodsRefund"
}

type OrdersGoods struct {
	dao.Entity
	OID             dao.PrimaryKey                  `gorm:"column:OID"`
	OrdersGoodsNo   string                          `gorm:"column:OrdersGoodsNo;unique"` //
	Status          OrdersGoodsStatus               `gorm:"column:Status"`               //OGAskRefund，OGRefundNo，OGRefundOk，OGRefundInfo，OGRefundComplete
	RefundID        dao.PrimaryKey                  `gorm:"column:RefundID"`             //RefundInfo json 退款退货信息
	OrdersID        dao.PrimaryKey                  `gorm:"column:OrdersID"`             //
	GoodsID         dao.PrimaryKey                  `gorm:"column:GoodsID"`
	SpecificationID dao.PrimaryKey                  `gorm:"column:SpecificationID"`
	Goods           Goods                           `gorm:"column:Goods;serializer:json;type:json"`         //josn,快照使用
	Specification   Specification                   `gorm:"column:Specification;serializer:json;type:json"` //json,快照使用
	GoodsSkus       Array[GoodsSku]                 `gorm:"column:GoodsSkus;type:json"`                     //json,只是对Specification进行分组选择使用的
	Discounts       sqltype.Array[sqltype.Discount] `gorm:"column:Discounts;type:json"`
	Quantity        uint                            `gorm:"column:Quantity"`           //数量
	CostPrice       uint                            `gorm:"column:CostPrice"`          //单价-原价
	SellPrice       uint                            `gorm:"column:SellPrice"`          //单价-销售价
	TotalBrokerage  uint                            `gorm:"column:TotalBrokerage"`     //总佣金
	Error           string                          `gorm:"column:Error"`              //
	Image           string                          `gorm:"column:Image"`              //Deprecated:当前规格的图片，如果规格没有图片，使用产品主图的第一张
	Pictures        sqltype.Array[sqltype.Image]    `gorm:"column:Pictures;type:JSON"` //规格的多张图片,当前规格的图片，如果规格没有图片，使用产品主图的第一张
	//SpecificationID uint `gorm:"column:SpecificationID"`             //
	//CollageNo     string `gorm:"column:CollageNo"` //拼团码，每个订单都是唯一
	//TimeSellID     uint `gorm:"column:TimeSellID"`             //限时抢购ID
	//TimeSell       string `gorm:"column:TimeSell;type:text"` //json
	//GoodsID         uint `gorm:"column:GoodsID"`                     //
}

/*
	func (m OrdersGoods) GetGoods() *Goods {
		var goods Goods
		_ = json.Unmarshal([]byte(m.Goods), &goods)
		return &goods
	}

	func (m OrdersGoods) GetSpecification() *Specification {
		var specification Specification
		_ = json.Unmarshal([]byte(m.Specification), &specification)
		return &specification
	}
*/
func (m *OrdersGoods) AddError(err string) {
	if strings.EqualFold(m.Error, "") {
		m.Error = err
	} else {
		m.Error = m.Error + "|" + err
	}
}
func (OrdersGoods) TableName() string {
	return "OrdersGoods"
}
