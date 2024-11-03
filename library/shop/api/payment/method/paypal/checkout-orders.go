package paypal

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/library/shop/api/payment/method/paypal/internal"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"regexp"
	"strconv"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
)

var nameReg = regexp.MustCompile(`\s+`)

type CheckoutOrders struct {
	User *model.User `mapping:""`
	Post struct {
		OrderNo               string
		AddressID             dao.PrimaryKey
		AdditionalInformation string
	} `method:"post"`
}

var matchRegexp = regexp.MustCompile("^(https:)([/|.|\\w|\\s|-])*\\.(?:jpg|gif|png|jpeg|JPG|GIF|PNG|JPEG)")

func matchImageUrl(url string) bool {
	return matchRegexp.MatchString(url)
}
func (m *CheckoutOrders) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}
func (m *CheckoutOrders) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	orders := repository.OrdersDao.GetOrdersByOrderNo(m.Post.OrderNo)
	if orders.ID == 0 {
		return nil, errors.New("no order found")
	}
	var orderDetails *internal.OrderDetailsResponse
	if len(orders.PrepayID) > 0 {
		var err error
		orderDetails, err = internal.OrderDetails(ctx, m.User.OID, orders.PrepayID)
		if err != nil {
			return nil, err
		}
	}

	address := dao.GetByPrimaryKey(db.Orm(), &model.Address{}, m.Post.AddressID).(*model.Address)
	if address.ID == 0 {
		return nil, errors.New("地址不能为空")
	}

	name := &internal.Name{}
	{
		names := nameReg.Split(address.Name, -1)
		if len(names) >= 2 {
			name.Surname = names[0]
			name.GivenName = strings.Join(names[1:], " ")
		}
	}

	shippingAddress := &internal.Address{}
	shippingAddress.SetAddress(address)

	confirmOrdersGoods, err := service.Order.Orders.AnalyseOrdersGoodsListByOrders(&orders, address)
	if err != nil {
		return nil, err
	}

	dns := repository.DNSDao.GetDefaultDNS(orders.OID)

	items := make([]internal.CheckoutOrdersUnitItem, 0)

	for i := range confirmOrdersGoods.OrdersGoodsInfos {
		ordersGoods := confirmOrdersGoods.OrdersGoodsInfos[i]

		imageOSS, err := oss.ReadUrl(ctx, ordersGoods.Image)
		if err != nil {
			return nil, err
		}
		if matchImageUrl(imageOSS) == false {
			imageOSS = ""
		}

		title := ordersGoods.Goods.Title
		if len(title) > 126 {
			title = ordersGoods.Goods.Title[:126]
		}
		summary := ordersGoods.Goods.Summary
		if len(summary) > 126 {
			summary = ordersGoods.Goods.Summary[:126]
		}

		amount := strconv.FormatFloat(float64(ordersGoods.SellPrice*ordersGoods.Quantity)/100.0, 'f', 2, 64)

		items = append(items, internal.CheckoutOrdersUnitItem{
			Name:        title,
			Quantity:    fmt.Sprintf("%d", ordersGoods.Quantity),
			Description: summary,
			Sku:         fmt.Sprintf("%d-%d", ordersGoods.Goods.ID, ordersGoods.Specification.ID),
			Url:         fmt.Sprintf("https://%s/product/detail/%d", dns.Domain, ordersGoods.Goods.ID),
			Category:    "",
			ImageUrl:    imageOSS,
			//ImageUrl: "https://oss.dev.com/assets/usokay.com/goods/143/image/c000af85aeaf74aeee732ec303da31ba.png@convert_from=webp",
			UnitAmount: internal.CheckoutOrdersUnitItemUnitAmount{
				CurrencyCode: "USD",
				Value:        amount,
			},
		})

	}

	if orderDetails != nil && len(orderDetails.Id) > 0 {
		uors := make([]internal.UpdateOrderChange, 0)

		uors = append(uors, internal.UpdateOrderChange{
			Op:    "replace",
			Path:  fmt.Sprintf("/purchase_units/@reference_id=='%s'/shipping/name", orders.OrderNo),
			Value: &internal.Name{FullName: name.GetFullName()},
		})

		uors = append(uors, internal.UpdateOrderChange{
			Op:    "replace",
			Path:  fmt.Sprintf("/purchase_units/@reference_id=='%s'/shipping/address", orders.OrderNo),
			Value: shippingAddress,
		})

		err = internal.UpdateOrder(ctx, m.User.OID, &internal.UpdateOrderRequest{
			Id:         orderDetails.Id,
			ChangeList: uors,
		})
		if err != nil {
			return nil, err
		}
		return result.NewData(orderDetails), err
	} else {
		units := make([]internal.CheckoutOrdersUnit, 0)
		amount := strconv.FormatFloat(float64(confirmOrdersGoods.TotalAmount)/100.0, 'f', 2, 64)
		units = append(units, internal.CheckoutOrdersUnit{
			ReferenceId: orders.OrderNo, //fmt.Sprintf("%d-%d", info.OrdersGoods.Goods.ID, info.OrdersGoods.Specification.ID),
			Description: m.Post.AdditionalInformation,
			Amount: serviceargument.Amount{
				CurrencyCode: "USD",
				Value:        amount,
				Breakdown:    serviceargument.Breakdown{ItemTotal: serviceargument.ItemTotal{CurrencyCode: "USD", Value: amount}},
			},
			Items: items,
			Shipping: &internal.Shipping{
				Name:    &internal.Name{FullName: name.GetFullName()},
				Type:    "SHIPPING",
				Address: shippingAddress,
			},
		})
		
		checkoutOrders, err := internal.CheckoutOrders(ctx, m.User.OID, &internal.CheckoutOrdersRequest{
			Intent:        "CAPTURE",
			PurchaseUnits: units,
		})
		if err != nil {
			return nil, err
		}
		err = dao.UpdateByPrimaryKey(db.Orm(), &model.Orders{}, orders.ID, map[string]any{"PrepayID": checkoutOrders.Id})
		return result.NewData(checkoutOrders), err
	}

}
