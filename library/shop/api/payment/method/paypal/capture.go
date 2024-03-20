package paypal

import (
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/shop/api/payment/method/paypal/internal/network"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/tool/object"
	"time"
)

type Capture struct {
	User *model.User `mapping:""`
	Post struct {
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
	if len(capture.PurchaseUnits) == 0 {
		return nil, errors.New("payment failed,invalid order")
	}
	mOrder := service.Order.Orders.GetOrdersByOrderNo(capture.PurchaseUnits[0].ReferenceId)
	if mOrder.IsZero() {
		return nil, errors.New("unable to confirm order, confirmation order failed")
	}
	if mOrder.PayMoney != object.ParseUint(object.ParseFloat(capture.PurchaseUnits[0].Payments.Captures[0].Amount.Value)*100) {
		return nil, errors.New("the order could not be confirmed, and the payment amount did not match the order amount")
	}

	/*var address model.Address
	err=util.JSONToStruct(mOrder.Address,&address)
	if err != nil {
		return nil, err
	}*/

	changeOrder := &model.Orders{
		IsPay:     model.OrdersIsPayPayed,
		PayMethod: model.OrdersPayMethodPaypal,
		Status:    model.OrdersStatusPay,
		//todo paypal如果改了收货地址，这里要改一下， Address:   capture.PurchaseUnits[0].Shipping.,
		PayTime: time.Now(),
	}
	err = dao.UpdateByPrimaryKey(db.Orm(), &model.Orders{}, mOrder.ID, changeOrder)
	if err != nil {
		return nil, err
	}
	return result.NewData(capture), nil
}
