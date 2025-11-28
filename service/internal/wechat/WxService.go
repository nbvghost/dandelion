package wechat

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/file"
	"github.com/nbvghost/dandelion/service/internal/user"
	"github.com/shopspring/decimal"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"

	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/library/db"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"

	"github.com/nbvghost/tool/encryption"

	"github.com/nbvghost/tool"
	"github.com/nbvghost/tool/collections"
)

type WxService struct {
	model.BaseDao
	//Admin service.AdminService
	//Goods goods.GoodsService
	User user.UserService
	//Orders       order.OrdersService
	Organization company.OrganizationService
	FileService  file.FileService

	Config             *model.WechatConfig
	OID                dao.PrimaryKey
	AccessTokenService AccessTokenService
}

var ticketMap = make(map[string]*Ticket) //&Ticket{}

var verifyCache = &struct {
	//ComponentVerifyTicket             string

	Component_access_token            string
	Component_access_token_expires_in int64
	Component_access_token_update     int64

	Pre_auth_code            string
	Pre_auth_code_expires_in int64
	Pre_auth_code_update     int64
}{}

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

	// 发送请求，以下载微信支付平台证书为例
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay5_1.shtml
	/*svc := certificates.CertificatesApiService{Client: client}
	resp, result, err := svc.DownloadCertificates(ctx)
	if err != nil {
		log.Fatalf("new wechatpay pay client err:%s", err)
		return nil, err
	}
	log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)*/
	return client, nil
}

/*
小程序
*/
/*func (service WxService) MiniProgramByAppId(db *gorm.DB, appId string) *model.WechatConfig {
	var wc model.WechatConfig
	db.Model(model.WechatConfig{}).Where(`"AppID"=?`, appId).Take(&wc)
	return &wc
}*/
func (m WxService) getConfig() *model.WechatConfig {
	if m.Config == nil {
		m.Config = &model.WechatConfig{}
		db.Orm().Model(model.WechatConfig{}).Where(`"OID"=?`, m.OID).Take(m.Config)
		return m.Config
	}
	return m.Config
}
func (m WxService) MiniProgram(db *gorm.DB) []dao.IEntity {
	//var wc []model.WechatConfig
	//db.Model(model.WechatConfig{}).Where(`"OID"=?`, OID).Take(&wc)
	return dao.Find(db, entity.WechatConfig).List()
}

type DeliveryInfo struct {
	DeliveryId   string `json:"delivery_id"`
	DeliveryName string `json:"delivery_name"`
}

func (m WxService) GetTraceWaybill(context constrain.IContext, ordersID dao.PrimaryKey, OrdersShipping *model.OrdersShipping) (string, error) {
	//trace_waybill
	////https://api.weixin.qq.com/cgi-bin/express/delivery/open_msg/trace_waybill?access_token=XXX

	orders := dao.GetByPrimaryKey(db.Orm(), &model.Orders{}, ordersID).(*model.Orders)
	ordersGoodsList := dao.Find(db.Orm(), &model.OrdersGoods{}).Where(`"OrdersID"=?`, orders.ID).List() //.(*model.OrdersGoods)

	u := dao.GetByPrimaryKey(db.Orm(), &model.User{}, orders.UserID).(*model.User)

	var address = &model.Address{}
	err := json.Unmarshal([]byte(orders.Address), address)
	if err != nil {
		return "", err
	}

	wxConfig := m.getConfig()

	ossUrl, err := oss.Url(context)
	if err != nil {
		return "", err
	}

	var r = struct {
		ErrCode      int    `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
		WaybillToken string `json:"waybill_token"`
	}{}

	{
		detailList := make([]map[string]any, 0)
		for i := range ordersGoodsList {
			item := ordersGoodsList[i].(*model.OrdersGoods)

			goods := item.Goods

			detailList = append(detailList, map[string]any{
				"goods_name":    goods.Title,
				"goods_img_url": "https:" + ossUrl + item.Image,
				//"goods_img_url": "https://oss.sites.ink/assets/default/goods/111/image/c12da74a2153dcfb5a9acc6478a12ee1.webp",
				"goods_desc": goods.Summary,
			})
		}

		marshal, err := json.Marshal(map[string]any{
			"openid":            u.OpenID,
			"receiver_phone":    address.Tel,
			"waybill_id":        OrdersShipping.No,
			"trans_id":          orders.TransactionID,
			"order_detail_path": fmt.Sprintf("/pages/order_info/order_info?ID=%d", orders.ID),
			"goods_info": map[string]any{
				"detail_list": detailList,
			},
		})
		if err != nil {
			return "", err
		}

		url := "https://api.weixin.qq.com/cgi-bin/express/delivery/open_msg/trace_waybill?access_token=" + m.AccessTokenService.GetAccessToken(wxConfig)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshal))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(body, &r)
		if err != nil {
			return "", err
		}
		if r.ErrCode != 0 {
			return "", errors.New(r.ErrMsg)
		}
	}
	return r.WaybillToken, nil
}
func (m WxService) GetDeliveryList(accessToken string) ([]DeliveryInfo, error) {
	url := "https://api.weixin.qq.com/cgi-bin/express/delivery/open_msg/get_delivery_list?access_token=" + accessToken

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var list struct {
		ErrCode      int            `json:"errcode"`
		DeliveryList []DeliveryInfo `json:"delivery_list"`
		Count        int            `json:"count"`
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}
	return list.DeliveryList, nil
}

func (m WxService) GetWXAConfig(prepay_id string, WxConfig *model.WechatConfig) (map[string]string, error) {

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

	//paySign := encryption.Md5ByString(list.Join("&") + "&key=" + WxConfig.MchAPIv2Key)
	outData["paySign"] = paySign
	return outData, nil
}
func (m WxService) SignatureVerification(dataMap util.Map, wxConfig MiniApp) bool {

	//appid := dataMap["appid"]
	//mch_id := dataMap["mch_id"]

	list := &collections.ListString{}
	for k, v := range dataMap {
		if !strings.EqualFold("sign", k) {
			list.Append(k + "=" + v)
		}

	}
	list.SortL()

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + wxConfig.MchAPIv2Key)

	//fmt.Println(list.Join("&") + "&key=" + WxConfig.PayKey)
	//fmt.Println(sign)

	if strings.EqualFold(dataMap["sign"], sign) {
		return true
	} else {
		return false
	}

}

func (m WxService) Order(ctx context.Context, OrderNo string, title, description string, detail, openid string, IP string, Money uint, attach string, wxConfig *model.WechatConfig) (Success result.ActionResultCode, Message string, wxResult *jsapi.PrepayWithRequestPaymentResponse) {
	client, err := NewClient(wxConfig)
	if err != nil {
		return result.Fail, err.Error(), nil
	}
	svc := jsapi.JsapiApiService{Client: client}
	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, _, err := svc.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(wxConfig.AppID),
			Mchid:       core.String(wxConfig.MchID),
			Description: core.String(title + "-" + description),
			OutTradeNo:  core.String(OrderNo),
			Attach:      core.String(attach),
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
		return result.Fail, err.Error(), nil
	}
	/*if !strings.EqualFold(detail, "") {
		mapData["detail"] = detail
	}*/
	//mapData["detail"] = `{ "goods_detail":[ { "goods_id":"iphone6s_16G", "wxpay_goods_id":"1001", "goods_name":"iPhone6s 16G", "quantity":1, "price":528800, "goods_category":"123456", "body":"苹果手机" }, { "goods_id":"iphone6s_32G", "wxpay_goods_id":"1002", "goods_name":"iPhone6s 32G", "quantity":1, "price":608800, "goods_category":"123789", "body":"苹果手机" } ] }`
	//mapData["nonce_str"] = tool.UUID()

	//mapData["openid"] = openid

	//mapData["spbill_create_ip"] = IP
	//mapData["total_fee"] = strconv.Itoa(int(Money))
	//mapData["trade_type"] = "JSAPI"
	//mapData["sign_type"] = "MD5"

	/*list := &collections.ListString{}
	for k, v := range mapData {
		list.Append(k + "=" + v)
	}
	list.SortL()
	sign := encryption.Md5ByString(list.Join("&") + "&key=" + wxConfig.PayKey)
	mapData["sign"] = sign
	xmlb, _ := xml.Marshal(&mapData)
	strReader := strings.NewReader(string(xmlb))*/

	/*if !strings.EqualFold(wxResult.Return_code, "SUCCESS") {
		//return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: resultXML.Return_msg, Data: nil}}
		return result.Fail, wxResult.Return_msg, nil
	}*/

	/*if !strings.EqualFold(wxResult.Result_code, "SUCCESS") {
		//return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: resultXML.Err_code_des, Data: nil}}
		return result.Fail, wxResult.Err_code_des, nil
	}*/

	return result.Success, "下单成功", resp
}
func (m WxService) MPOrder(ctx context.Context, OrderNo string, title, description string, ogs []model.OrdersGoods, openid string, IP string, Money uint, attach string, wxConfig *model.WechatConfig) (Success result.ActionResultCode, Message string, wxResult *jsapi.PrepayWithRequestPaymentResponse) {

	CostGoodsPrice := int64(0)

	goods_detail := make([]map[string]interface{}, 0)
	for _, value := range ogs {
		goodsObj := make(map[string]interface{})
		goodsObj["goods_id"] = value.OrdersGoodsNo

		/*var goods model.Goods
		json.Unmarshal([]byte(value.Goods), &goods)

		var specification model.Specification
		json.Unmarshal([]byte(value.Specification), &specification)*/

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

	//WxConfig := service.MiniProgram()

	return m.Order(ctx, OrderNo, title, description, string(detailB), openid, IP, Money, attach, wxConfig)
}

// func (self WxService) GetWxConfig(DB *gorm.DB, CompanyID uint) *WxConfig {
// 	content_item := &WxConfig{}
// 	err := DB.Where("CompanyID=?", CompanyID).First(content_item).Error
// 	log.Println(err)

// 	if content_item.ID == 0 {
// 		err = DB.Create(content_item).Error
// 		log.Println(err)
// 		return content_item
// 	} else {
// 		return content_item
// 	}
// }
/*func (entity WxService) ChangeWxConfig(DB *gorm.DB, ID uint, Value model.WxConfig) error {

	//content_item := b.GetWxConfig(DB, CompanyID)
	//content_item.V = Value
	return DB.Model(&model.WxConfig{}).Where("ID=?", ID).Updates(Value).Error
}*/

/*func (entity WxService) MiniProgramByAppIDAndMchID(AppID, MchID string) model.WxConfig {
	var wx model.WxConfig
	err := singleton.Orm().Model(&model.WxConfig{}).Where("AppID=? and MchID=?", AppID, MchID).First(&wx).Error
	log.Println(err)
	return wx
}*/
/*func (entity WxService) GetWxConfig(ID uint) model.WxConfig {
	var wx model.WxConfig
	err := singleton.Orm().Model(&model.WxConfig{}).Where("ID=?", ID).First(&wx).Error
	log.Println(err)
	return wx
}*/
func (m WxService) MWQRCodeTemp(OID uint, UserID uint, qrtype, params string, wxConfig *model.WechatConfig) *result.ActionResult {

	//user := context.Session.Attributes.Get(play.SessionUser).(*model.User)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//context.Request.ParseForm()
	//Page := context.Request.FormValue("Page")
	//Page := context.Request.URL.Query().Get("Page")
	//MyShareKey := tool.Hashids{}.Encode(user.ID)

	access_token := m.AccessTokenService.GetAccessToken(wxConfig)

	postData := make(map[string]interface{})
	//results := make(map[string]interface{})

	postData["expire_seconds"] = 2592000
	postData["action_name"] = "QR_STR_SCENE"
	postData["action_info"] = map[string]interface{}{"scene": map[string]interface{}{"scene_str": strconv.Itoa(int(UserID)) + "|" + qrtype + "|" + params}}

	body := strings.NewReader(util.StructToJSON(postData))
	//postData := url.Values{}
	//postData.Add("scene","sdfsd")
	resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="+access_token, "application/json", body)
	if err != nil {
		return &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}
	}
	//fmt.Println(string(b))
	defer resp.Body.Close()
	path, err := m.FileService.WriteTempFile(b, "image/png")
	return &result.ActionResult{Code: result.Success, Message: "", Data: path}

}

// 订单查询
func (m WxService) OrderQuery(ctx context.Context, OrderNo string, wxConfig *model.WechatConfig) (*payments.Transaction, error) {
	client, err := NewClient(wxConfig)
	if err != nil {
		return nil, err
	}
	svc := jsapi.JsapiApiService{Client: client}

	resp, _, err := svc.QueryOrderByOutTradeNo(ctx,
		jsapi.QueryOrderByOutTradeNoRequest{
			OutTradeNo: core.String(OrderNo),
			Mchid:      core.String(wxConfig.MchID),
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTransfersInfo 查询提现接口
func (m WxService) GetTransfersInfo(transfers *model.Transfers, wxConfig *model.WechatConfig) (*transferbatch.TransferBatchGet, error) {
	mchPrivateKey, err := utils.LoadPrivateKey(wxConfig.PrivateKey)
	if err != nil {
		log.Printf("load merchant private key error:%s", err)
		return nil, err
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(wxConfig.MchID, wxConfig.MchCertificateSerialNumber, mchPrivateKey, wxConfig.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechatpay pay client err:%s", err)
		return nil, err
	}
	svc := transferbatch.TransferBatchApiService{Client: client}
	resp, apiResult, err := svc.GetTransferBatchByOutNo(ctx,
		transferbatch.GetTransferBatchByOutNoRequest{
			OutBatchNo:      core.String(transfers.OrderNo),
			NeedQueryDetail: core.Bool(false),
			Offset:          core.Int64(0),
			Limit:           core.Int64(100),
			DetailStatus:    core.String("ALL"),
		},
	)
	if err != nil {
		// 处理错误
		log.Printf("call GetTransferBatchByOutNo err:%s", err)
		return nil, err
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", apiResult.Response.StatusCode, resp)
		if apiResult.Response.StatusCode != 200 {
			return nil, errors.New("转账查询请求错误")
		}
		return resp.TransferBatch, nil
	}
}

// Transfers 提现
func (m WxService) Transfers(transfers model.Transfers, transferDetailInputs []transferbatch.TransferDetailInput, wxConfig *model.WechatConfig) error {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKey(wxConfig.PrivateKey)
	if err != nil {
		log.Printf("load merchant private key error:%s", err)
		return err
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(wxConfig.MchID, wxConfig.MchCertificateSerialNumber, mchPrivateKey, wxConfig.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechatpay pay client err:%s", err)
		return err
	}

	svc := transferbatch.TransferBatchApiService{Client: client}
	resp, apiResult, err := svc.InitiateBatchTransfer(ctx,
		transferbatch.InitiateBatchTransferRequest{
			Appid:       core.String(wxConfig.AppID),
			OutBatchNo:  core.String(transfers.OrderNo),
			BatchName:   core.String("提现"),
			BatchRemark: core.String(transfers.Desc),
			TotalAmount: core.Int64(int64(transfers.Amount)),
			TotalNum:    core.Int64(int64(len(transferDetailInputs))),
			/*TransferDetailList: []transferbatch.TransferDetailInput{transferbatch.TransferDetailInput{
				OutDetailNo:    core.String(fmt.Sprintf("%d", userJournal.ID)),
				TransferAmount: core.Int64(int64(transfers.Amount)),
				TransferRemark: core.String("用户余额提现"),
				Openid:         core.String(transfers.OpenId),
				UserName:       core.String(transfers.ReUserName),
			}},*/
			TransferDetailList: transferDetailInputs,
		},
	)

	if err != nil {
		if apiError, ok := err.(*core.APIError); ok {
			return errors.New(apiError.Message)
		}
		// 处理错误
		log.Printf("call GetTransferBatchByNo err:%s", err)
		return err
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", apiResult.Response.StatusCode, resp.String())
		if apiResult.Response.StatusCode != 200 {
			return errors.New("提现请求错误")
		}
		return nil
	}
}

// 关闭订单
func (m WxService) CloseOrder(OrderNo string, OID dao.PrimaryKey, wxConfig *model.WechatConfig) (Success bool, Message string) {

	//WxConfig := service.MiniProgram()

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
	log.Println(err)

	reader := strings.NewReader(string(b))
	response, err := http.Post("https://api.mch.weixin.qq.com/pay/closeorder", "text/xml", reader)
	log.Println(err)

	b, err = ioutil.ReadAll(response.Body)
	log.Println(err)

	fmt.Println(string(b))

	var inData = make(util.Map)
	err = xml.Unmarshal(b, &inData)
	log.Println(err)
	//fmt.Println(inData)

	if strings.EqualFold(inData["return_code"], "SUCCESS") && strings.EqualFold(inData["result_code"], "SUCCESS") {
		Success = true
		Message = "订单关闭成功"
		return
	} else {
		//loggerService := service.LoggerService{}
		//loggerService.Error("Appointment:"+strconv.Itoa(int(OrderNo)), inData["err_code"]+":"+inData["err_code_des"])
		if strings.EqualFold(inData["return_code"], "FAIL") {
			Success = false
			Message = inData["return_msg"]
			return
		} else {
			Success = false
			Message = inData["result_msg"]
			return
		}
		//return false, inData["err_code_des"]
	}

}

// Refund 退款-订单内的所有的商品/订单内某个商品
func (m WxService) Refund(ctx context.Context, order *model.Orders, ordersGoods *model.OrdersGoods, reason string) error {

	client, err := NewClient(m.getConfig())
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
			NotifyUrl:    core.String(m.getConfig().RefundNotifyUrl),
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
			NotifyUrl:    core.String(m.getConfig().RefundNotifyUrl),
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

	resp, _, err := svc.Create(ctx, createRequest)
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
	//outMap["refund_desc"] = Desc

	/*if Type == 0 {
		outMap["refund_account"] = "REFUND_SOURCE_UNSETTLED_FUNDS" //0
	} else {
		outMap["refund_account"] = "REFUND_SOURCE_RECHARGE_FUNDS" //1
	}*/

	/*list := &collections.ListString{}
	for k, v := range outMap {
		list.Append(k + "=" + v)
	}
	list.SortL()

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + wxConfig.MchAPIv2Key)
	outMap["sign"] = sign

	b, err := xml.MarshalIndent(util.Map(outMap), "", "")
	log.Println(err)

	// Load client cert
	cert, err := tls.LoadX509KeyPair("cert/wxpay/apiclient_cert.pem", "cert/wxpay/apiclient_key.pem")
	if err != nil {
		log.Fatal(err)
	}*/

	// Load CA cert
	//caCert, err := ioutil.ReadFile("cert/wxpay/rootca.pem")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	/*tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		//RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	reader := strings.NewReader(string(b))
	response, err := client.Post("https://api.mch.weixin.qq.com/secapi/pay/refund", "text/xml", reader)
	log.Println(err)
	if err != nil {
		Success = false
		Message = err.Error()
		return
	}

	b, err = ioutil.ReadAll(response.Body)
	log.Println(err)
	if err != nil {
		Success = false
		Message = err.Error()
		return
	}

	//fmt.Println(string(b))

	var inData = make(util.Map)
	err = xml.Unmarshal(b, &inData)
	log.Println(err)*/
	//fmt.Println(inData)

	/*if strings.EqualFold(inData["return_code"], "SUCCESS") && strings.EqualFold(inData["result_code"], "SUCCESS") {
		Success = true
		Message = "退款申请成功"
		return
	} else {
		//loggerService := service.LoggerService{}
		//loggerService.Error("Appointment:"+strconv.Itoa(int(OrderNo)), inData["err_code"]+":"+inData["err_code_des"])

		if strings.EqualFold(inData["return_code"], "FAIL") {
			Success = false
			Message = inData["return_msg"]
			return
		} else {
			//fmt.Println(inData["err_code"])
			//fmt.Println(inData["err_code_des"])

			err_code := "ORDERNOTEXIST,USER_ACCOUNT_ABNORMAL"

			if strings.Contains(err_code, inData["err_code"]) {
				Success = false
				Message = inData["err_code_des"]
				return

				//return true, inData["err_code_des"]
			}
			Success = false
			Message = inData["err_code_des"] + ":" + inData["err_code"]
			return
		}
		//return false, inData["err_code_des"]
	}*/
}
func (m WxService) Decrypt(encryptedData, session_key, iv_text string) (bool, string) {
	bkey, err := base64.StdEncoding.DecodeString(session_key)

	//aesKey := Base64.decodeBase64(encodingAesKey + "=");
	block, err := aes.NewCipher(bkey) //选择加密算法
	if err != nil {
		return false, ""
	}
	iv, err := base64.StdEncoding.DecodeString(iv_text)

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)

	blockModel := cipher.NewCBCDecrypter(block, iv)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)

	length := len(plantText)
	unpadding := int(plantText[length-1])
	return true, string(plantText[:(length - unpadding)])
}

func (m WxService) GetSHA1(token, timestamp, nonce, encrypt string) string {

	array := []string{timestamp, nonce, encrypt, token}
	sb := ""
	// 字符串排序
	sort.Strings(array)
	//fmt.Println(array)
	for i := 0; i < len(array); i++ {
		sb = sb + array[i]
	}
	// SHA1签名生成
	h := sha1.New()
	io.WriteString(h, sb)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (m WxService) MwGetTicket(WxConfig *model.WechatConfig) string {

	if ticketMap[WxConfig.AppID] != nil && (time.Now().Unix()-ticketMap[WxConfig.AppID].Update) < ticketMap[WxConfig.AppID].Expires_in {

		return ticketMap[WxConfig.AppID].Ticket
	}

	url := "http://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + m.AccessTokenService.GetAccessToken(WxConfig)

	resp, err := http.Get(url)
	log.Println(err)

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	d := make(map[string]interface{})

	err = json.Unmarshal(b, &d)
	log.Println(err)
	//fmt.Println(string(b))
	//fmt.Println(d)
	if d["ticket"] != nil && d["expires_in"] != nil {
		ticketMap[WxConfig.AppID] = &Ticket{}
		ticketMap[WxConfig.AppID].Ticket = d["ticket"].(string)
		ticketMap[WxConfig.AppID].Expires_in = int64(d["expires_in"].(float64))
		ticketMap[WxConfig.AppID].Update = time.Now().Unix()

		return ticketMap[WxConfig.AppID].Ticket
	} else {
		return ""
	}

}
func (m WxService) MwGetWXJSConfig(url string, OID dao.PrimaryKey) map[string]interface{} {

	wxConfig := m.getConfig()

	appId := wxConfig.AppID
	timestamp := time.Now().Unix()
	nonceStr := tool.UUID()
	//chooseWXPay
	list := &collections.ListString{}
	list.Append("noncestr=" + nonceStr)
	list.Append("jsapi_ticket=" + m.MwGetTicket(wxConfig))
	list.Append("timestamp=" + strconv.FormatInt(timestamp, 10))

	_url := strings.Split(url, "#")[0]
	list.Append("url=" + _url)
	list.SortL()
	signature := util.SignSha1(list.Join("&"))

	results := make(map[string]interface{})
	results["appId"] = appId
	results["timestamp"] = timestamp
	results["nonceStr"] = nonceStr
	results["signature"] = signature

	return results
}
func (m WxService) GetWechatConfig(tx *gorm.DB, OID dao.PrimaryKey) model.WechatConfig {
	var contentConfig model.WechatConfig
	tx.Model(&model.WechatConfig{}).Where(map[string]interface{}{"OID": OID}).First(&contentConfig)
	return contentConfig
}
func (m WxService) InitWechatConfig(tx *gorm.DB, OID dao.PrimaryKey) error {
	item := m.GetWechatConfig(tx, OID)
	if !item.IsZero() {
		return nil //fmt.Errorf("已经存在Wechat配制文件")
	}

	wechatConfig := &model.WechatConfig{
		OID:   OID,
		AppID: strings.ToLower(encryption.Md5ByString(fmt.Sprintf("%d", OID))),
	}
	return dao.Create(tx, wechatConfig)
}

//var GlobalWXConfig = model.WxConfig{CompanyID: -1, AppID: "wx037d3b26b2ba34b2", AppSecret: "fe3faa4e6f8abd87fa4621cb5ed5f725", Token: "30e6e3b03bf7ec6d2ce56a50055e1cd1", EncodingAESKey: "egMWQnCkbuDd7u5GM7EJBnH8mISn5iwAorjRNnFx3dv", MchID: "1342120901"}
