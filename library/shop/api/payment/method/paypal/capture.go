package paypal

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/library/shop/api/payment/method/paypal/internal"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/tool/object"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
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
	capture, err := internal.Capture(ctx, m.User.OID, &internal.CaptureRequest{PaypalOrderID: m.Post.PaypalOrderID})
	if err != nil {
		return nil, err
	}
	if len(capture.PurchaseUnits) == 0 {
		return nil, errors.New("payment failed,invalid order")
	}
	mOrder := repository.OrdersDao.GetOrdersByOrderNo(capture.PurchaseUnits[0].ReferenceId)
	if mOrder.IsZero() {
		return nil, errors.New("unable to confirm order, confirmation order failed")
	}

	if mOrder.PayMoney != object.ParseUint(object.ParseFloat(capture.PurchaseUnits[0].Payments.Captures[0].Amount.Value)*100) {
		return nil, errors.New("the order could not be confirmed, and the payment amount did not match the order amount")
	}

	changeOrder := &model.Orders{
		IsPay:     model.OrdersIsPayPayed,
		PayMethod: model.OrdersPayMethodPaypal,
		Status:    model.OrdersStatusPay,
		//todo paypal如果改了收货地址，这里要改一下， Address:   capture.PurchaseUnits[0].Shipping.,
		PayTime:       time.Now(),
		TransactionID: capture.PurchaseUnits[0].Payments.Captures[0].Id,
	}
	err = dao.UpdateByPrimaryKey(db.Orm(), &model.Orders{}, mOrder.ID, changeOrder)
	if err != nil {
		return nil, err
	}

	{
		botMessage := strings.Builder{}
		botMessage.WriteString(fmt.Sprintf("网站[%s]有新的订单支付成功，请注意查收。\n", ctx.AppName()))
		err = service.Wechat.SendText(botMessage.String())
		if err != nil {
			return nil, err
		}
	}

	return result.NewData(capture), nil
}
