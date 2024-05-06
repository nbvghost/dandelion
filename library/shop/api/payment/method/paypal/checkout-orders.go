package paypal

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/library/shop/api/payment/method/paypal/internal"
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

	unit := internal.CheckoutOrdersUnit{
		ReferenceId: orders.OrderNo, //fmt.Sprintf("%d-%d", info.OrdersGoods.Goods.ID, info.OrdersGoods.Specification.ID),
		Description: m.Post.AdditionalInformation,
		Amount: internal.Amount{
			CurrencyCode: "USD",
			Value:        strconv.FormatFloat(float64(confirmOrdersGoods.TotalAmount)/100.0, 'f', 2, 64),
		},
		Shipping: &internal.Shipping{
			Name:    &internal.Name{FullName: name.GetFullName()},
			Type:    "SHIPPING",
			Address: shippingAddress,
		},
	}

	if orderDetails != nil && len(orderDetails.Id) > 0 {
		err = internal.UpdateOrder(ctx, m.User.OID, &internal.UpdateOrderRequest{
			Id: orderDetails.Id,
			ChangeList: []internal.UpdateOrderChange{
				{
					Op:    "replace",
					Path:  fmt.Sprintf("/purchase_units/@reference_id=='%s'/shipping/name", orders.OrderNo),
					Value: &internal.Name{FullName: name.GetFullName()},
				},
				{
					Op:    "replace",
					Path:  fmt.Sprintf("/purchase_units/@reference_id=='%s'/shipping/address", orders.OrderNo),
					Value: shippingAddress,
				},
			},
		})
		if err != nil {
			return nil, err
		}
		return result.NewData(orderDetails), err
	} else {
		checkoutOrders, err := internal.CheckoutOrders(ctx, m.User.OID, &internal.CheckoutOrdersRequest{
			Intent:        "CAPTURE",
			PurchaseUnits: []internal.CheckoutOrdersUnit{unit},
		})
		if err != nil {
			return nil, err
		}
		err = dao.UpdateByPrimaryKey(db.Orm(), &model.Orders{}, orders.ID, map[string]any{"PrepayID": checkoutOrders.Id})
		return result.NewData(checkoutOrders), err
	}

}
