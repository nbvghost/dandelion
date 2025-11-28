package wechatpay

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"github.com/nbvghost/tool"
	"github.com/nbvghost/tool/collections"
	"github.com/nbvghost/tool/encryption"
	"github.com/shopspring/decimal"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type Service struct {
	OID     dao.PrimaryKey
	config  *model.WechatConfig
	Context constrain.IServiceContext
}

func (m *Service) Deliver(orders *model.Orders) error {
	return nil
}

func (m *Service) CloseOrder(OrderNo string) error {

	//WxConfig := service.MiniProgram()
	wxConfig := m.GetConfig()

	outMap := make(util.Map)
	outMap["appid"] = wxConfig.AppID
	outMap["mch_id"] = wxConfig.MchID
	outMap["nonce_str"] = tool.UUID()

	outMap["out_trade_no"] = OrderNo

	list := &collections.ListString{}
	for k, v := range outMap {
		list.Append(k + "=" + v)
	}
	list.SortL()

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + wxConfig.MchAPIv2Key)
	outMap["sign"] = sign

	b, err := xml.MarshalIndent(util.Map(outMap), "", "")
	if err != nil {
		return err
	}

	reader := strings.NewReader(string(b))
	response, err := http.Post("https://api.mch.weixin.qq.com/pay/closeorder", "text/xml", reader)
	if err != nil {
		return err
	}

	b, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	var inData = make(util.Map)
	err = xml.Unmarshal(b, &inData)
	if err != nil {
		return err
	}
	//fmt.Println(inData)

	if strings.EqualFold(inData["return_code"], "SUCCESS") && strings.EqualFold(inData["result_code"], "SUCCESS") {
		//Success = true
		//Message = "订单关闭成功"
		return nil
	} else {
		//loggerService := service.LoggerService{}
		//loggerService.Error("Appointment:"+strconv.Itoa(int(OrderNo)), inData["err_code"]+":"+inData["err_code_des"])
		if strings.EqualFold(inData["return_code"], "FAIL") {
			//Success = false
			//Message = inData["return_msg"]
			return errors.New(inData["return_msg"])
		} else {
			//Success = false
			//Message = inData["result_msg"]
			return errors.New(inData["return_msg"])
		}
		//return false, inData["err_code_des"]
	}
}

func NewClient(config *model.WechatConfig) (*core.Client, error) {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(config.MchID, config.MchCertificateSerialNumber, mchPrivateKey, config.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("new wechatpay pay client err:%s", err)
	}
	return client, nil
}

func (m *Service) GetConfig() *model.WechatConfig {
	if m.config == nil {
		m.config = &model.WechatConfig{}
		db.Orm().Model(model.WechatConfig{}).Where(`"OID"=?`, m.OID).Take(m.config)
		return m.config
	}
	return m.config
}
func (m *Service) OrderQuery(orders *model.Orders) (*serviceargument.OrderQueryResult, error) {
	client, err := NewClient(m.GetConfig())
	if err != nil {
		return nil, err
	}
	svc := jsapi.JsapiApiService{Client: client}

	resp, _, err := svc.QueryOrderByOutTradeNo(m.Context,
		jsapi.QueryOrderByOutTradeNoRequest{
			OutTradeNo: core.String(orders.OrderNo),
			Mchid:      core.String(m.GetConfig().MchID),
		},
	)
	if err != nil {
		return nil, err
	}

	payTime, err := time.ParseInLocation("2006-01-02T15:04:05-07:00", *resp.SuccessTime, time.Local)
	if err != nil {
		return nil, err
	}

	return &serviceargument.OrderQueryResult{
		State:            serviceargument.OrderQueryState(*resp.TradeState),
		PayerTotalAmount: *resp.Amount.PayerTotal,
		PayTime:          payTime,
		OutTradeNo:       *resp.OutTradeNo,
		TransactionID:    *resp.TransactionId,
		Attach:           *resp.Attach,
	}, nil
}

func (m *Service) Order(OrderNo string, title, description string, detail, openid string, IP string, Money uint, ordersType model.OrdersType) (*serviceargument.OrderResult, error) {

	wxConfig := m.GetConfig()

	client, err := NewClient(wxConfig)
	if err != nil {
		return nil, err
	}
	svc := jsapi.JsapiApiService{Client: client}
	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, _, err := svc.PrepayWithRequestPayment(m.Context,
		jsapi.PrepayRequest{
			Appid:       core.String(wxConfig.AppID),
			Mchid:       core.String(wxConfig.MchID),
			Description: core.String(title + "-" + description),
			OutTradeNo:  core.String(OrderNo),
			Attach:      core.String(string(ordersType)),
			NotifyUrl:   core.String(wxConfig.OrderNotifyUrl),
			Amount: &jsapi.Amount{
				Total:    core.Int64(int64(Money)),
				Currency: core.String("CNY"),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(openid),
			},
			SceneInfo: &jsapi.SceneInfo{
				PayerClientIp: core.String(IP),
			},
		},
	)

	if err != nil {
		return nil, err
	}
	return &serviceargument.OrderResult{
		PrepayId: *resp.PrepayId,
	}, nil
}

func (m *Service) Refund(order *model.Orders, ordersGoods *model.OrdersGoods, reason string) error {

	client, err := NewClient(m.GetConfig())
	if err != nil {
		return err
	}
	svc := refunddomestic.RefundsApiService{Client: client}

	var createRequest refunddomestic.CreateRequest

	if ordersGoods == nil {
		//退全部
		//outMap["out_refund_no"] = order.OrderNo
		//outMap["out_trade_no"] = order.OrderNo
		//outMap["refund_fee"] = strconv.Itoa(int(order.PayMoney))
		//outMap["total_fee"] = strconv.Itoa(int(order.PayMoney))

		createRequest = refunddomestic.CreateRequest{
			OutTradeNo:   core.String(order.OrderNo),
			OutRefundNo:  core.String(order.OrderNo),
			Reason:       core.String(reason),
			NotifyUrl:    core.String(m.GetConfig().RefundNotifyUrl),
			FundsAccount: refunddomestic.REQFUNDSACCOUNT_AVAILABLE.Ptr(),
			Amount: &refunddomestic.AmountReq{
				Refund:   core.Int64(int64(order.PayMoney)),
				Total:    core.Int64(int64(order.PayMoney)),
				From:     nil,
				Currency: core.String("CNY"),
			},
			GoodsDetail: nil,
		}
	} else {
		//部分退款
		//op := model.Orders.GetOrdersPackageByOrderNo(order.OrdersPackageNo)
		//op := Orders.GetOrdersByOrderNo(order.OrdersPackageNo)
		//outMap["out_refund_no"] = order.OrderNo
		//outMap["out_trade_no"] = order.OrdersPackageNo
		//outMap["refund_fee"] = strconv.Itoa(int(order.PayMoney))
		//outMap["total_fee"] = strconv.Itoa(int(ordersPackage.TotalPayMoney))

		goods := ordersGoods.Goods
		specification := ordersGoods.Specification

		refundAmount := ordersGoods.SellPrice.Mul(decimal.NewFromUint64(uint64(ordersGoods.Quantity)))

		createRequest = refunddomestic.CreateRequest{
			OutTradeNo:   core.String(order.OrderNo),
			OutRefundNo:  core.String(ordersGoods.OrdersGoodsNo),
			Reason:       core.String(reason),
			NotifyUrl:    core.String(m.GetConfig().RefundNotifyUrl),
			FundsAccount: refunddomestic.REQFUNDSACCOUNT_AVAILABLE.Ptr(),
			Amount: &refunddomestic.AmountReq{
				Refund:   core.Int64(int64(refundAmount.IntPart())),
				Total:    core.Int64(int64(order.PayMoney)),
				From:     nil,
				Currency: core.String("CNY"),
			},
			GoodsDetail: []refunddomestic.GoodsDetail{
				{
					MerchantGoodsId: core.String(fmt.Sprintf("%d", goods.ID)),
					//WechatpayGoodsId: core.String(fmt.Sprintf("%d",goods.ID)),
					GoodsName:      core.String(goods.Title + "/" + specification.Label),
					UnitPrice:      core.Int64(int64(specification.MarketPrice.IntPart())),
					RefundAmount:   core.Int64(int64(refundAmount.IntPart())),
					RefundQuantity: core.Int64(int64(ordersGoods.Quantity)),
				},
			},
		}
	}

	resp, _, err := svc.Create(m.Context, createRequest)
	if err != nil {
		return err
	}

	switch *resp.Status {
	case refunddomestic.STATUS_SUCCESS:
		log.Println(fmt.Sprintf("%s:%s\n", *resp.OutTradeNo, "退款成功"))
	case refunddomestic.STATUS_CLOSED:
		log.Println(fmt.Sprintf("%s:%s\n", *resp.OutTradeNo, "退款关闭"))
	case refunddomestic.STATUS_PROCESSING:
		log.Println(fmt.Sprintf("%s:%s\n", *resp.OutTradeNo, "退款处理中"))
	case refunddomestic.STATUS_ABNORMAL:
		return errors.New("退款异常,请联系客服处理")
	default:
		return errors.New("未知操作,请联系客服处理")
	}

	//SUCCESS: 退款成功
	//CLOSED: 退款关闭
	//PROCESSING: 退款处理中
	//ABNORMAL: 退款异常

	return nil
}
func (m *Service) GetWXAConfig(prepay_id string) (map[string]string, error) {
	WxConfig := m.GetConfig()
	rsaCryptoUtilCertificate, err := utils.LoadPrivateKey(WxConfig.PrivateKey)
	if err != nil {
		log.Fatal("load merchant private key error")
		return nil, err
	}

	outData := make(map[string]string)
	outData["appId"] = WxConfig.AppID
	outData["timeStamp"] = strconv.Itoa(int(time.Now().Unix()))
	outData["nonceStr"] = tool.UUID()
	outData["package"] = "prepay_id=" + prepay_id
	outData["signType"] = "RSA"

	list := &collections.ListString{}
	list.Append(outData["appId"])
	list.Append(outData["timeStamp"])
	list.Append(outData["nonceStr"])
	list.Append(outData["package"])

	paySign, err := utils.SignSHA256WithRSA(fmt.Sprintf("%s\n%s\n%s\n%s\n", outData["appId"], outData["timeStamp"], outData["nonceStr"], outData["package"]), rsaCryptoUtilCertificate)
	if err != nil {
		return nil, err
	}
	outData["paySign"] = paySign
	return outData, nil
}
func (m *Service) MPOrder(OrderNo string, title, description string, ogs []model.OrdersGoods, openid string, IP string, Money uint, ordersType model.OrdersType) (*serviceargument.OrderResult, error) {

	CostGoodsPrice := int64(0)

	goods_detail := make([]map[string]interface{}, 0)
	for _, value := range ogs {
		goodsObj := make(map[string]interface{})
		goodsObj["goods_id"] = value.OrdersGoodsNo
		goodsObj["goods_name"] = value.Goods.Title + "-" + value.Specification.Label
		goodsObj["quantity"] = value.Quantity
		goodsObj["price"] = value.SellPrice
		goods_detail = append(goods_detail, goodsObj)

		CostGoodsPrice = CostGoodsPrice + value.CostPrice.IntPart()
	}

	detail := make(map[string]interface{})
	detail["cost_price"] = CostGoodsPrice
	//detail["receipt_id"] = CostGoodsPrice
	detail["goods_detail"] = goods_detail

	golgaldetail := make(map[string]interface{})
	golgaldetail["version"] = 1.0
	//golgaldetail["goods_tag"] = 1.0
	golgaldetail["detail"] = detail

	detailB, _ := json.Marshal(&golgaldetail)

	return m.Order(OrderNo, title, description, string(detailB), openid, IP, Money, ordersType)
}
