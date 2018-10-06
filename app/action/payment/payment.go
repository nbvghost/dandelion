package payment

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"

	"dandelion/app/util"

	"dandelion/app/service"

	"log"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type Controller struct {
	gweb.BaseController
	Orders service.OrdersService
	Wx     service.WxService
}

func (controller *Controller) Apply() {
	//controller.AddHandler(gweb.POSMethod("wxpay", controller.wxpayAction))
	controller.AddHandler(gweb.ALLMethod("notify", controller.notifyAction))
}
func (controller *Controller) notifyAction(context *gweb.Context) gweb.Result {

	resultMap := make(util.Map)

	b, err := ioutil.ReadAll(context.Request.Body)
	tool.Trace(err)

	defer context.Request.Body.Close()
	fmt.Println(string(b))

	err = xml.Unmarshal(b, &resultMap)
	tool.Trace(err)

	outXML := ``

	if controller.Wx.SignatureVerification(resultMap) {

		TotalFee, _ := strconv.ParseUint(resultMap["total_fee"], 10, 64)
		//OrderNo := result["out_trade_no"]
		//TimeEnd := result["time_end"]
		//attach := result["attach"]
		Success, Message := controller.Orders.OrderNotify(TotalFee, resultMap["out_trade_no"], resultMap["time_end"], resultMap["attach"])
		//self.Orders.OrderNotify(result)
		//Success, Message := controller.Orders.OrderNotify(resultMap)
		if Success == false {
			log.Println(Message)
			outXML = `<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[` + Message + `]]></return_msg></xml>`
		} else {
			outXML = `<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>`
		}
	} else {
		outXML = `<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[签名失败]]></return_msg></xml>`
	}

	return &gweb.XMLResult{Data: outXML}
}
