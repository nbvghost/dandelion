package extends

import "github.com/nbvghost/dandelion/entity/model"

type OrdersGoodsInfo struct {
	OrdersGoods model.OrdersGoods
	Discounts   []Discount
}
