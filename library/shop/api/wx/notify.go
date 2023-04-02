package wx

import (
	"log"
	"net/http"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/wechat"
	"github.com/nbvghost/gpa/types"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

type Notify struct {
	WxService     wechat.WxService
	OrdersService order.OrdersService
	Get           struct {
		OID types.PrimaryKey `uri:"OID"`
	} `method:"Get"`
	Post struct {
		OID types.PrimaryKey `uri:"OID"`
	} `method:"Post"`
}

func (m *Notify) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return m.handle(context, m.Get.OID)
}
func (m *Notify) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	return m.handle(context, m.Post.OID)
}

func (m *Notify) handle(context constrain.IContext, OID types.PrimaryKey) (r constrain.IResult, err error) {
	wxConfig := m.WxService.MiniProgramByOID(singleton.Orm(), OID)

	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(wxConfig.MchID)

	content := new(payments.Transaction)

	handler, err := notify.NewRSANotifyHandler(wxConfig.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		log.Println(err)
		return result.NewJsonResult(map[string]any{"code": "FAIL", "message": err.Error()}).WithStatusCode(http.StatusBadRequest), nil
	}

	// 4. 使用client进行接口调用
	contextValue := contexext.FromContext(context)
	request, err := handler.ParseNotifyRequest(context, contextValue.Request, content)
	if err != nil {
		log.Println(err)
		return result.NewJsonResult(map[string]any{"code": "FAIL", "message": err.Error()}).WithStatusCode(http.StatusBadRequest), nil
	}
	log.Println(request.Resource.Plaintext)
	log.Println(content)

	payTime, err := time.ParseInLocation("2006-01-02T15:04:05-07:00", *content.SuccessTime, time.Local)
	if err != nil {
		log.Println(err)
		return result.NewJsonResult(map[string]any{"code": "FAIL", "message": err.Error()}).WithStatusCode(http.StatusBadRequest), nil
	}
	message, err := m.OrdersService.OrderNotify(uint(*content.Amount.PayerTotal), *content.OutTradeNo, payTime, *content.Attach)
	if err != nil {
		return result.NewJsonResult(map[string]any{"code": "FAIL", "message": message}).WithStatusCode(http.StatusBadRequest), nil
	} else {
		return result.NewJsonResult(map[string]any{"code": "SUCCESS", "message": ""}), nil
	}
}
