package wx

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/wechat"
	"log"
)

type Message struct {
	WxService     wechat.WxService
	OrdersService order.OrdersService
	Get           struct {
		OID       dao.PrimaryKey `uri:"OID"`
		Signature string         `form:"signature"`
		Timestamp string         `form:"timestamp"`
		Nonce     string         `form:"nonce"`
		EchoStr   string         `form:"echostr"`
	} `method:"Get"`
}

func (m *Message) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	log.Println("message", m.Get.OID, m.Get.Signature, m.Get.Timestamp, m.Get.Nonce, m.Get.EchoStr)
	return result.NewTextResult(m.Get.EchoStr), nil
}
