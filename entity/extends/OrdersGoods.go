package extends

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa/types"
)

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
