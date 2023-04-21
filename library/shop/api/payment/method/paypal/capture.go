package paypal

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/shop/api/payment/method/paypal/internal/network"
	"github.com/nbvghost/dandelion/service/configuration"
	"github.com/nbvghost/dandelion/service/order"
	"log"
)

type Capture struct {
	ConfigurationService configuration.ConfigurationService
	ShoppingCartService  order.ShoppingCartService
	User                 *model.User `mapping:""`
	Post                 struct {
		PaypalOrderID string `uri:"PaypalOrderID"`
	} `method:"post"`
}

func (m *Capture) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
func (m *Capture) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	capture, err := network.Capture(ctx, m.User.OID, &network.CaptureRequest{PaypalOrderID: m.Post.PaypalOrderID})
	if err != nil {
		return nil, err
	}
	log.Println(capture)
	return result.NewData(capture), nil
}
