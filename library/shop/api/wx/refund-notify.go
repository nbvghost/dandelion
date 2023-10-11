package wx

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/wechat"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"log"
	"net/http"
)

type RefundNotify struct {
	WxService     wechat.WxService
	OrdersService order.OrdersService
	Get           struct {
		OID dao.PrimaryKey `uri:"OID"`
	} `method:"Get"`
	Post struct {
		OID dao.PrimaryKey `uri:"OID"`
	} `method:"Post"`
}

func (m *RefundNotify) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return m.handle(context, m.Get.OID)
}
func (m *RefundNotify) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	return m.handle(context, m.Post.OID)
}

func (m *RefundNotify) handle(context constrain.IContext, OID dao.PrimaryKey) (r constrain.IResult, err error) {
	wxConfig := m.WxService.MiniProgramByOID(db.Orm(), OID)

	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(wxConfig.MchID)

	content := new(refunddomestic.Refund)

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

	return result.NewJsonResult(map[string]any{"code": "FAIL", "message": ""}).WithStatusCode(http.StatusBadRequest), nil

}
