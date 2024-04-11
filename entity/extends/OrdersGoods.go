package extends

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type OrdersGoodsInfo struct {
	OrdersGoods *OrdersGoodsMix
	Discounts   []Discount
}

// ElementStatus 由于要把列表中的每一项的状态告诉前端，增加这个结构给列表中的每项做标记
type ElementStatus struct {
	IsError bool
	Error   string
}

// ConfirmOrdersGoods 商家订单计算结构，
type ConfirmOrdersGoods struct {
	//Organization     *model.Organization //如果涉及多个商家的订单的话，则使用OrdersGoodsInfos[].OrdersGoods.OID来区分
	OrdersGoodsInfos []OrdersGoodsInfo
	FavouredPrice    uint //拼团价格
	FullCutAll       uint
	GoodsPrice       uint
	ExpressPrice     uint
	FullCut          model.FullCut
	Address          *model.Address
	TotalAmount      uint
}

type OrdersGoodsMix struct {
	OID             dao.PrimaryKey                              //
	OrdersGoodsNo   string                                      //
	Status          model.OrdersGoodsStatus                     //OGAskRefund，OGRefundNo，OGRefundOk，OGRefundInfo，OGRefundComplete
	RefundInfo      string                                      //RefundInfo json 退款退货信息
	OrdersID        dao.PrimaryKey                              //
	Goods           *model.Goods                                //josn
	Specification   *model.Specification                        //json
	Discounts       []Discount                                  //
	Quantity        uint                                        //数量
	CostPrice       uint                                        //单价-原价
	SellPrice       uint                                        //单价-销售价
	TotalBrokerage  uint                                        //总佣金
	Image           string                                      //当前规格的图片，如果规格没有图片，使用产品主图的第一张
	SkuImages       []string                                    //
	SkuLabelMap     map[dao.PrimaryKey]*model.GoodsSkuLabel     //
	SkuLabelDataMap map[dao.PrimaryKey]*model.GoodsSkuLabelData //
	ElementStatus   ElementStatus
}
