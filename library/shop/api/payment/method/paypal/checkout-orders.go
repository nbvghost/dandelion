package paypal

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/shop/api/payment/method/paypal/internal/network"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/configuration"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/gpa/types"
	"regexp"
	"strconv"
	"strings"
)

type CheckoutOrders struct {
	ConfigurationService configuration.ConfigurationService
	ShoppingCartService  order.ShoppingCartService
	User                 *model.User `mapping:""`
	Post                 struct {
		AddressID             types.PrimaryKey
		AdditionalInformation string
	} `method:"post"`
}

func (m *CheckoutOrders) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}
func (m *CheckoutOrders) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	address := dao.GetByPrimaryKey(singleton.Orm(), &model.Address{}, m.Post.AddressID).(*model.Address)
	if address.ID == 0 {
		return nil, errors.New("地址不能为空")
	}

	var name network.Name
	{
		names := nameReg.Split(address.Name, -1)
		if len(names) >= 2 {
			name.Surname = names[0]
			name.GivenName = strings.Join(names[1:], " ")
		}
	}

	shippingAddress := &network.Address{}
	shippingAddress.SetAddress(address)

	list, _, err := m.ShoppingCartService.FindShoppingCartListDetails(m.User.ID, address)
	if err != nil {
		return nil, err
	}
	units := make([]network.CheckoutOrdersUnit, 0)
	for _, goods := range list {
		for _, info := range goods.OrdersGoodsInfos {
			units = append(units, network.CheckoutOrdersUnit{
				ReferenceId: fmt.Sprintf("%d-%d", info.OrdersGoods.Goods.ID, info.OrdersGoods.Specification.ID),
				Description: m.Post.AdditionalInformation,
				Amount: network.Amount{
					CurrencyCode: "USD",
					Value:        strconv.FormatFloat(float64(info.OrdersGoods.Specification.MarketPrice)/100.0, 'f', 2, 64),
				},
				Shipping: &network.Shipping{
					Name:    name,
					Type:    "SHIPPING",
					Address: shippingAddress,
				},
			})
		}

	}

	/*paypalBillingAddress := network.PayPalAddress{}
	{
		var billingAddress *model.Address
		if m.Post.BillingAddressID > 0 {
			billingAddress = dao.GetByPrimaryKey(singleton.Orm(), &model.Address{}, m.Post.BillingAddressID).(*model.Address)
		} else {
			billingAddress = dao.GetBy(singleton.Orm(), &model.Address{}, map[string]any{"UserID": m.User.ID, "DefaultBilling": true}).(*model.Address)
			if billingAddress.IsEmpty() {
				billingAddress = address
			}
		}
		paypalBillingAddress.AddressLine1 = billingAddress.Detail
		if len(billingAddress.Company) > 0 {
			paypalBillingAddress.AddressLine2 = billingAddress.Company
		}
		paypalBillingAddress.AdminArea1 = billingAddress.CountyName + "." + billingAddress.ProvinceName
		paypalBillingAddress.AdminArea2 = billingAddress.CityName
		paypalBillingAddress.PostalCode = billingAddress.PostalCode
	}*/

	orders, err := network.CheckoutOrders(ctx, m.User.OID, &network.CheckoutOrdersRequest{
		Intent:        "CAPTURE",
		PurchaseUnits: units,
		//Payer:         payer,
		//PaymentSource: network.CheckoutOrdersPaymentSource{Card: card},
	})
	if err != nil {
		return nil, err
	}
	return result.NewData(orders), nil
}

var nameReg = regexp.MustCompile(`\s+`)
