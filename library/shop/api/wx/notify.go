package wx

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service"
	"log"
	"net/http"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

type Notify struct {
	Get struct {
		OID dao.PrimaryKey `uri:"OID"`
	} `method:"Get"`
	Post struct {
		OID dao.PrimaryKey `uri:"OID"`
	} `method:"Post"`
}

func (m *Notify) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return m.handle(context, m.Get.OID)
}
func (m *Notify) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	return m.handle(context, m.Post.OID)
}

func (m *Notify) handle(context constrain.IContext, OID dao.PrimaryKey) (r constrain.IResult, err error) {
	wxConfig := service.Payment.NewWechat(context, OID).GetConfig() //service.Wechat.Wx.MiniProgramByOID(db.Orm(), OID)

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
	message, err := service.Order.Orders.OrderPaySuccess(uint(*content.Amount.PayerTotal), *content.OutTradeNo, *content.TransactionId, payTime, model.OrdersType(*content.Attach))
	if err != nil {
		return result.NewJsonResult(map[string]any{"code": "FAIL", "message": message}).WithStatusCode(http.StatusBadRequest), nil
	} else {
		return result.NewJsonResult(map[string]any{"code": "SUCCESS", "message": ""}), nil
	}
}
