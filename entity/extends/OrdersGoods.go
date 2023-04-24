package extends

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa/types"
)

// 商家订单计算结构，
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

type OrdersGoods struct {
	OID             types.PrimaryKey                              //
	OrdersGoodsNo   string                                        //
	Status          model.OrdersGoodsStatus                       //OGAskRefund，OGRefundNo，OGRefundOk，OGRefundInfo，OGRefundComplete
	RefundInfo      string                                        //RefundInfo json 退款退货信息
	OrdersID        types.PrimaryKey                              //
	Goods           *model.Goods                                  //josn
	Specification   *model.Specification                          //json
	Discounts       []Discount                                    //
	Quantity        uint                                          //数量
	CostPrice       uint                                          //单价-原价
	SellPrice       uint                                          //单价-销售价
	TotalBrokerage  uint                                          //总佣金
	SkuImages       []string                                      //
	SkuLabelMap     map[types.PrimaryKey]*model.GoodsSkuLabel     //
	SkuLabelDataMap map[types.PrimaryKey]*model.GoodsSkuLabelData //
}
