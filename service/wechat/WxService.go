package wechat

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"

	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/gweb"

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
}

var accessTokenMap = make(map[string]*AccessToken)
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
		log.Fatal("load merchant private key error")
		return nil, err
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(config.MchID, config.MchCertificateSerialNumber, mchPrivateKey, config.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
		return nil, err
	}

	// 发送请求，以下载微信支付平台证书为例
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay5_1.shtml
	/*svc := certificates.CertificatesApiService{Client: client}
	resp, result, err := svc.DownloadCertificates(ctx)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
		return nil, err
	}
	log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)*/
	return client, nil
}

/*
小程序
*/
func (service WxService) MiniProgramByAppId(db *gorm.DB, appId string) *model.WechatConfig {
	var wc model.WechatConfig
	db.Model(model.WechatConfig{}).Where(`"AppID"=?`, appId).Take(&wc)
	return &wc
}
func (service WxService) MiniProgramByOID(db *gorm.DB, OID types.PrimaryKey) *model.WechatConfig {
	var wc model.WechatConfig
	db.Model(model.WechatConfig{}).Where(`"OID"=?`, OID).Take(&wc)
	return &wc
}
func (service WxService) MiniProgram(db *gorm.DB) []types.IEntity {
	//var wc []model.WechatConfig
	//db.Model(model.WechatConfig{}).Where(`"OID"=?`, OID).Take(&wc)
	return dao.Find(db, entity.WechatConfig).List()
}
func (service WxService) GetAccessToken(WxConfig *model.WechatConfig) string {

	if accessTokenMap[WxConfig.AppID] != nil && (time.Now().Unix()-accessTokenMap[WxConfig.AppID].Update) < accessTokenMap[WxConfig.AppID].Expires_in {

		return accessTokenMap[WxConfig.AppID].Access_token
	}

	//WxConfig := model.GetWxConfig(WxConfigID)

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + WxConfig.AppID + "&secret=" + WxConfig.AppSecret

	resp, err := http.Get(url)
	log.Println(err)

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	d := make(map[string]interface{})

	err = json.Unmarshal(b, &d)
	log.Println(err)
	//fmt.Println(string(b))
	//fmt.Println(d)
	if d["access_token"] == nil {
		return ""
	}
	at := &AccessToken{}
	at.Access_token = d["access_token"].(string)
	at.Expires_in = int64(d["expires_in"].(float64))
	at.Update = time.Now().Unix()
	accessTokenMap[WxConfig.AppID] = at
	return accessTokenMap[WxConfig.AppID].Access_token
}

func (service WxService) GetWXAConfig(prepay_id string, WxConfig *model.WechatConfig) (map[string]string, error) {

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
func (service WxService) SignatureVerification(dataMap util.Map, wxConfig MiniApp) bool {

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

func (service WxService) MiniProgramInfo(Code, AppID, AppSecret string) (err error, OpenID, SessionKey string) {

	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + AppID + "&secret=" + AppSecret + "&js_code=" + Code + "&grant_type=authorization_code")
	if err == nil {
		b, _ := ioutil.ReadAll(resp.Body)

		readData := make(map[string]interface{})

		fmt.Println(string(b))
		json.Unmarshal(b, &readData)

		if readData["openid"] != nil && readData["session_key"] != nil {

			OpenID := readData["openid"].(string)
			SessionKey := readData["session_key"].(string)

			return nil, OpenID, SessionKey
		} else {
			if readData["errmsg"] != nil {
				return errors.New("登陆失败:" + readData["errmsg"].(string)), "", ""
			} else {
				return errors.New("登陆失败"), "", ""
			}
		}

	} else {
		return errors.New("登陆失败:" + err.Error()), "", ""
	}

}

func (service WxService) Order(ctx context.Context, OrderNo string, title, description string, detail, openid string, IP string, Money uint, attach string, wxConfig *model.WechatConfig) (Success result.ActionResultCode, Message string, wxResult *jsapi.PrepayWithRequestPaymentResponse) {
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
func (service WxService) MPOrder(ctx context.Context, OrderNo string, title, description string, ogs []model.OrdersGoods, openid string, IP string, Money uint, attach string, wxConfig *model.WechatConfig) (Success result.ActionResultCode, Message string, wxResult *jsapi.PrepayWithRequestPaymentResponse) {

	CostGoodsPrice := uint(0)

	goods_detail := make([]map[string]interface{}, 0)
	for _, value := range ogs {
		goodsObj := make(map[string]interface{})
		goodsObj["goods_id"] = value.OrdersGoodsNo

		var goods model.Goods
		json.Unmarshal([]byte(value.Goods), &goods)

		var specification model.Specification
		json.Unmarshal([]byte(value.Specification), &specification)

		goodsObj["goods_name"] = goods.Title + "-" + specification.Label
		goodsObj["quantity"] = value.Quantity
		goodsObj["price"] = value.SellPrice
		goods_detail = append(goods_detail, goodsObj)

		CostGoodsPrice = CostGoodsPrice + value.CostPrice
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

	return service.Order(ctx, OrderNo, title, description, string(detailB), openid, IP, Money, attach, wxConfig)
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
func (service WxService) MWQRCodeTemp(OID uint, UserID uint, qrtype, params string, wxConfig *model.WechatConfig) *result.ActionResult {

	//user := context.Session.Attributes.Get(play.SessionUser).(*model.User)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//context.Request.ParseForm()
	//Page := context.Request.FormValue("Page")
	//Page := context.Request.URL.Query().Get("Page")
	//MyShareKey := tool.Hashids{}.Encode(user.ID)

	access_token := service.GetAccessToken(wxConfig)

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
	path := gweb.WriteTempFile(b, "image/png")
	return &result.ActionResult{Code: result.Success, Message: "", Data: path}

}

//订单查询
func (service WxService) OrderQuery(ctx context.Context, OrderNo string, wxConfig *model.WechatConfig) (*payments.Transaction, error) {
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

//查询提现接口
func (service WxService) GetTransfersInfo(transfers model.Transfers, wxConfig *model.WechatConfig) (Success bool) {

	//WxConfig := service.MiniProgram()

	outMap := make(util.Map)
	outMap["nonce_str"] = tool.UUID()
	outMap["partner_trade_no"] = transfers.OrderNo
	outMap["mch_id"] = wxConfig.MchID
	outMap["appid"] = wxConfig.AppID

	list := &collections.ListString{}
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
		return false
	}

	// Load CA cert
	/*caCert, err := ioutil.ReadFile("cert/wxpay/rootca.pem")
	if err != nil {
		log.Fatal(err)
	}*/
	//caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		//RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	reader := strings.NewReader(string(b))
	response, err := client.Post("https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo", "text/xml", reader)
	log.Println(err)
	if err != nil {
		return false
	}

	b, err = ioutil.ReadAll(response.Body)
	log.Println(err)
	if err != nil {
		return false
	}

	//fmt.Println(string(b))

	var inData = make(util.Map)
	err = xml.Unmarshal(b, &inData)
	log.Println(err)
	if err != nil {
		return false
	}

	//fmt.Println(inData)

	if strings.EqualFold(inData["return_code"], "SUCCESS") && strings.EqualFold(inData["result_code"], "SUCCESS") {
		Success = true
		return
	} else {
		//loggerService := service.LoggerService{}
		//loggerService.Error("Appointment:"+strconv.Itoa(int(OrderNo)), inData["err_code"]+":"+inData["err_code_des"])

		if strings.EqualFold(inData["return_code"], "FAIL") {
			Success = false
			return
		} else {
			//fmt.Println(inData["err_code"])
			//fmt.Println(inData["err_code_des"])
			Success = false
			return
		}
		//return false, inData["err_code_des"]
	}

}

//提现
func (service WxService) Transfers(transfers model.Transfers, wxConfig *model.WechatConfig) (Success bool, Message string) {
	//WxConfig := service.MiniProgram()

	outMap := make(util.Map)
	outMap["mch_appid"] = wxConfig.AppID
	outMap["mchid"] = wxConfig.MchID
	outMap["nonce_str"] = tool.UUID()

	outMap["partner_trade_no"] = transfers.OrderNo
	outMap["openid"] = transfers.OpenId
	outMap["check_name"] = "FORCE_CHECK"
	outMap["re_user_name"] = transfers.ReUserName
	outMap["amount"] = strconv.Itoa(int(transfers.Amount))
	outMap["desc"] = transfers.Desc
	outMap["spbill_create_ip"] = transfers.IP

	list := &collections.ListString{}
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
	}

	// Load CA cert
	/*caCert, err := ioutil.ReadFile("cert/wxpay/rootca.pem")
	if err != nil {
		log.Fatal(err)
	}*/
	//caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		//RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	reader := strings.NewReader(string(b))
	response, err := client.Post("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", "text/xml", reader)
	log.Println(err)

	b, err = ioutil.ReadAll(response.Body)
	log.Println(err)

	//fmt.Println(string(b))

	var inData = make(util.Map)
	err = xml.Unmarshal(b, &inData)
	log.Println(err)

	//fmt.Println(inData)

	if strings.EqualFold(inData["return_code"], "SUCCESS") && strings.EqualFold(inData["result_code"], "SUCCESS") {
		Success = true
		Message = "提现申请成功"
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
			Success = false
			Message = inData["err_code"] + ":" + inData["err_code_des"]
			return
		}
		//return false, inData["err_code_des"]
	}
}

//关闭订单
func (service WxService) CloseOrder(OrderNo string, OID types.PrimaryKey, wxConfig *model.WechatConfig) (Success bool, Message string) {

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

//退款-订单内的所有的商品/订单内某个商品
func (service WxService) Refund(ctx context.Context, order *model.Orders, ordersGoods *model.OrdersGoods, reason string, wxConfig *model.WechatConfig) (*refunddomestic.Refund, error) {

	client, err := NewClient(wxConfig)
	if err != nil {
		return nil, err
	}
	svc := refunddomestic.RefundsApiService{Client: client}

	var createRequest refunddomestic.CreateRequest

	if ordersGoods == nil {
		//outMap["out_refund_no"] = order.OrderNo
		//outMap["out_trade_no"] = order.OrderNo
		//outMap["refund_fee"] = strconv.Itoa(int(order.PayMoney))
		//outMap["total_fee"] = strconv.Itoa(int(order.PayMoney))
		createRequest = refunddomestic.CreateRequest{
			OutTradeNo:   core.String(order.OrderNo),
			OutRefundNo:  core.String(order.OrderNo),
			Reason:       core.String(reason),
			NotifyUrl:    core.String(wxConfig.RefundNotifyUrl),
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
		//op := model.Orders.GetOrdersPackageByOrderNo(order.OrdersPackageNo)
		//op := Orders.GetOrdersByOrderNo(order.OrdersPackageNo)
		//outMap["out_refund_no"] = order.OrderNo
		//outMap["out_trade_no"] = order.OrdersPackageNo
		//outMap["refund_fee"] = strconv.Itoa(int(order.PayMoney))
		//outMap["total_fee"] = strconv.Itoa(int(ordersPackage.TotalPayMoney))
		createRequest = refunddomestic.CreateRequest{
			OutTradeNo:   core.String(order.OrderNo),
			OutRefundNo:  core.String(order.OrdersPackageNo),
			Reason:       core.String(reason),
			NotifyUrl:    core.String(wxConfig.RefundNotifyUrl),
			FundsAccount: refunddomestic.REQFUNDSACCOUNT_AVAILABLE.Ptr(),
			Amount: &refunddomestic.AmountReq{
				Refund:   core.Int64(int64(ordersGoods.SellPrice)),
				Total:    core.Int64(int64(order.PayMoney)),
				From:     nil,
				Currency: core.String("CNY"),
			},
			GoodsDetail: nil,
		}
	}

	resp, _, err := svc.Create(ctx, createRequest)
	if err != nil {
		return nil, err
	}
	return resp, nil
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
func (service WxService) Decrypt(encryptedData, session_key, iv_text string) (bool, string) {
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

func (service WxService) getSHA1(token, timestamp, nonce, encrypt string) string {

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

func (service WxService) MwGetTicket(WxConfig *model.WechatConfig) string {

	if ticketMap[WxConfig.AppID] != nil && (time.Now().Unix()-ticketMap[WxConfig.AppID].Update) < ticketMap[WxConfig.AppID].Expires_in {

		return ticketMap[WxConfig.AppID].Ticket
	}

	url := "http://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + service.GetAccessToken(WxConfig)

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
func (service WxService) MwGetWXJSConfig(url string, OID types.PrimaryKey) map[string]interface{} {

	wxConfig := service.MiniProgramByOID(singleton.Orm(), OID)

	appId := wxConfig.AppID
	timestamp := time.Now().Unix()
	nonceStr := tool.UUID()
	//chooseWXPay
	list := &collections.ListString{}
	list.Append("noncestr=" + nonceStr)
	list.Append("jsapi_ticket=" + service.MwGetTicket(wxConfig))
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

//var GlobalWXConfig = model.WxConfig{CompanyID: -1, AppID: "wx037d3b26b2ba34b2", AppSecret: "fe3faa4e6f8abd87fa4621cb5ed5f725", Token: "30e6e3b03bf7ec6d2ce56a50055e1cd1", EncodingAESKey: "egMWQnCkbuDd7u5GM7EJBnH8mISn5iwAorjRNnFx3dv", MchID: "1342120901"}
