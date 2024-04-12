package wx

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"log"
	"net/http"
	"strings"
)

type RefundNotify struct {
	Get struct {
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
	wxConfig := service.Wechat.Wx.MiniProgramByOID(db.Orm(), OID)

	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(wxConfig.MchID)

	content := new(serviceargument.RefundNotifyData)

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

	if strings.EqualFold(content.RefundStatus, "SUCCESS") || strings.EqualFold(content.RefundStatus, "CLOSED") {
		orders := repository.OrdersDao.GetOrdersByOrderNo(content.OutTradeNo)
		if orders.IsZero() {
			return result.NewJsonResult(map[string]any{"code": "FAIL", "message": "订单不存在"}).WithStatusCode(http.StatusBadRequest), nil
		}

		//var ordersGoods *model.OrdersGoods
		if !strings.EqualFold(content.OutTradeNo, content.OutRefundNo) {
			//ordersGoods = m.OrdersService.GetOrdersGoodsByOrdersGoodsNo(db.Orm(), content.OutTradeNo)
		}

		err = service.Order.Orders.OrdersRefundSuccess(&orders)
		if err != nil {
			return result.NewJsonResult(map[string]any{"code": "FAIL", "message": err.Error()}).WithStatusCode(http.StatusBadRequest), nil
		}
		return result.NewJsonResult(map[string]any{"code": "SUCCESS", "message": ""}), nil
	}

	return result.NewJsonResult(map[string]any{"code": "FAIL", "message": ""}).WithStatusCode(http.StatusBadRequest), nil

}
