package wechat

import (
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/company"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/user"
	"github.com/nbvghost/gweb/tool/encryption"

	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"github.com/nbvghost/dandelion/app/util"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"encoding/json"

	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
	"github.com/nbvghost/gweb/tool/collections"
)

type WxService struct {
	dao.BaseDao
	//Admin service.AdminService
	//Goods goods.GoodsService
	User user.UserService
	//Orders       order.OrdersService
	Organization company.OrganizationService
}

type TokenXML struct {
	AppId   string `xml:AppId`
	Encrypt string `xml:Encrypt`
}
type AccessToken struct {
	Access_token string
	Expires_in   int64
	Update       int64
}
type Ticket struct {
	Ticket     string
	Expires_in int64
	Update     int64
}

type PushInfo struct {
	AppId                 string `xml:AppId`
	CreateTime            int64  `xml:CreateTime`
	InfoType              string `xml:InfoType`
	ComponentVerifyTicket string `xml:ComponentVerifyTicket`
}
type WxOrderResult struct {
	Return_code  string `xml:"return_code"`
	Return_msg   string `xml:"return_msg"`
	Appid        string `xml:"appid"`
	Mch_id       string `xml:"mch_id"`
	Nonce_str    string `xml:"nonce_str"`
	Sign         string `xml:"sign"`
	Result_code  string `xml:"result_code"`
	Prepay_id    string `xml:"prepay_id"`
	Trade_type   string `xml:"trade_type"`
	Err_code_des string `xml:"err_code_des"`
}

type WXDetail struct {
	Goods_detail []WXGoodsDetail `json:"goods_detail"`
}
type WXGoodsDetail struct {
	Goods_id   string `json:"goods_id"`
	Goods_name string `json:"goods_name"`
	Quantity   string `json:"quantity"`
	Price      string `json:"price"`
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

func (entity WxService) GetAccessToken(WxConfig dao.WxConfig) string {

	if accessTokenMap[WxConfig.AppID] != nil && (time.Now().Unix()-accessTokenMap[WxConfig.AppID].Update) < accessTokenMap[WxConfig.AppID].Expires_in {

		return accessTokenMap[WxConfig.AppID].Access_token
	}

	//WxConfig := entity.GetWxConfig(WxConfigID)

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + WxConfig.AppID + "&secret=" + WxConfig.AppSecret

	resp, err := http.Get(url)
	glog.Error(err)

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	d := make(map[string]interface{})

	err = json.Unmarshal(b, &d)
	glog.Error(err)
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

func (entity WxService) GetWXAConfig(prepay_id string, WxConfig dao.WxConfig) (outData map[string]string) {
	//WxConfig := entity.MiniProgram()
	outData = make(map[string]string)
	outData["appId"] = WxConfig.AppID
	outData["timeStamp"] = strconv.Itoa(int(time.Now().Unix()))
	outData["nonceStr"] = tool.UUID()
	outData["package"] = "prepay_id=" + prepay_id
	outData["signType"] = "MD5"

	list := &collections.ListString{}
	for k, v := range outData {
		list.Append(k + "=" + v)
	}

	list.SortL()

	paySign := encryption.Md5ByString(list.Join("&") + "&key=" + WxConfig.PayKey)
	outData["paySign"] = paySign
	return outData
}
func (entity WxService) SignatureVerification(dataMap util.Map) bool {

	//appid := dataMap["appid"]
	//mch_id := dataMap["mch_id"]
	WxConfig := entity.MiniProgram()

	list := &collections.ListString{}
	for k, v := range dataMap {
		if !strings.EqualFold("sign", k) {
			list.Append(k + "=" + v)
		}

	}
	list.SortL()

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + WxConfig.PayKey)

	//fmt.Println(list.Join("&") + "&key=" + WxConfig.PayKey)
	//fmt.Println(sign)

	if strings.EqualFold(dataMap["sign"], sign) {
		return true
	} else {
		return false
	}

}

func (entity WxService) MiniProgramInfo(Code, AppID, AppSecret string) (err error, OpenID, SessionKey string) {

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

//新用户加入，绑定上下级关系
func (entity WxService) NewUserJoinNotify(NewUser dao.User, notifyUser dao.User) *result.ActionResult {

	as := &result.ActionResult{}

	userFormID := entity.User.GetFromIDs(notifyUser.ID)
	if userFormID.ID == 0 {
		as.Code = result.ActionFail
		as.Message = "没有找到，用户的formid"
	} else {

		sendData := make(map[string]interface{})
		sendData["touser"] = notifyUser.OpenID

		weapp_template_msg_data := make(map[string]interface{})
		weapp_template_msg_data["page"] = "pages/user/user"
		weapp_template_msg_data["template_id"] = "YfEY2Xbju5-fm3Naww3EbVYQPUPIjorESo-KV-KXZvs"
		weapp_template_msg_data["form_id"] = userFormID.FormId

		data_data := make(map[string]interface{})
		data_data["keyword1"] = map[string]interface{}{"value": strconv.Itoa(int(NewUser.ID)), "color": "#173177"}
		if NewUser.Gender == 1 {
			data_data["keyword2"] = map[string]interface{}{"value": "男", "color": "#173177"}
		} else if NewUser.Gender == 2 {
			data_data["keyword2"] = map[string]interface{}{"value": "女", "color": "#173177"}
		} else {
			data_data["keyword2"] = map[string]interface{}{"value": "未知", "color": "#173177"}
		}

		data_data["keyword3"] = map[string]interface{}{"value": NewUser.CreatedAt.Format("2006-01-02 15:04:05"), "color": "#173177"}
		data_data["keyword4"] = map[string]interface{}{"value": NewUser.Name, "color": "#173177"}
		data_data["keyword5"] = map[string]interface{}{"value": NewUser.Name + "已经成为您的好友，他（她）下单您会获得奖励喔！", "color": "#173177"}

		weapp_template_msg_data["data"] = data_data

		sendData["weapp_template_msg"] = weapp_template_msg_data

		var errcode int
		as, errcode = entity.SendUniformMessage(sendData)
		if as.Code == result.ActionOK || errcode == 41028 {
			entity.User.Delete(dao.Orm(), &dao.UserFormIds{}, userFormID.ID)
		}

	}

	return as
}

//发货通知
func (entity WxService) OrderDeliveryNotify(Order dao.Orders, ogs []dao.OrdersGoods) *result.ActionResult {

	if Order.ID == 0 {
		return &result.ActionResult{Code: result.ActionFail, Message: "找不到订单", Data: nil}
	}

	var notifyUser dao.User
	entity.User.Get(dao.Orm(), Order.UserID, &notifyUser)

	var as *result.ActionResult

	weapp_template_msg_data := make(map[string]interface{})
	weapp_template_msg_data["page"] = "pages/user/user"
	weapp_template_msg_data["template_id"] = "MHiJR_3T2W4LJVhwOVctO6Lr7fxC9rSCO924dwSoYrY"
	weapp_template_msg_data["form_id"] = Order.PrepayID
	weapp_template_msg_data["touser"] = notifyUser.OpenID

	data_data := make(map[string]interface{})
	data_data["keyword1"] = map[string]interface{}{"value": Order.ShipName, "color": "#173177"}
	data_data["keyword2"] = map[string]interface{}{"value": Order.ShipNo, "color": "#173177"}

	var Titles = ""
	for _, value := range ogs {
		var goods dao.Goods
		json.Unmarshal([]byte(value.Goods), &goods)
		Titles += goods.Title
	}
	if len(Titles) > 48 {
		Titles = Titles[:48] + "等"
	}

	data_data["keyword3"] = map[string]interface{}{"value": Titles, "color": "#173177"}
	data_data["keyword4"] = map[string]interface{}{"value": Order.DeliverTime.Format("2006-01-02 15:04:05"), "color": "#173177"}

	weapp_template_msg_data["data"] = data_data

	as = entity.SendWXMessage(weapp_template_msg_data)

	return as
}

//收入提醒
/*
@slUser 收入的用户
*/
func (entity WxService) INComeNotify(slUser dao.User, itemName string, timeText string, typeText string) *result.ActionResult {
	//

	if slUser.ID == 0 {
		return &result.ActionResult{Code: result.ActionFail, Message: "用户不存在", Data: nil}
	}

	//var notifyUser dao.User
	//entity.User.Get(dao.Orm(), slUser.SuperiorID, &notifyUser)

	var as = &result.ActionResult{Code: result.ActionFail}

	userFormID := entity.User.GetFromIDs(slUser.ID)
	if userFormID.ID == 0 {
		as.Code = result.ActionFail
		as.Message = "没有找到，用户的formid"

	} else {

		sendData := make(map[string]interface{})
		sendData["touser"] = slUser.OpenID

		weapp_template_msg_data := make(map[string]interface{})
		weapp_template_msg_data["page"] = "pages/user/user"
		weapp_template_msg_data["template_id"] = "xV23xWZgdNViUiD1fk-1edKNY7QNJnv4SD6tY7pu8w4"
		weapp_template_msg_data["form_id"] = userFormID.FormId

		data_data := make(map[string]interface{})
		data_data["keyword1"] = map[string]interface{}{"value": itemName, "color": "#173177"}
		data_data["keyword2"] = map[string]interface{}{"value": timeText, "color": "#173177"}
		data_data["keyword3"] = map[string]interface{}{"value": typeText, "color": "#ff0000"}

		weapp_template_msg_data["data"] = data_data

		sendData["weapp_template_msg"] = weapp_template_msg_data

		var errcode int
		as, errcode = entity.SendUniformMessage(sendData)
		if as.Code == result.ActionOK || errcode == 41028 {
			entity.User.Delete(dao.Orm(), &dao.UserFormIds{}, userFormID.ID)
		}

	}

	return as
}

//新订单
func (entity WxService) NewOrderNotify(Order dao.Orders, ogs []dao.OrdersGoods) *result.ActionResult {

	if Order.ID == 0 {
		return &result.ActionResult{Code: result.ActionFail, Message: "找不到订单", Data: nil}
	}

	var notifyUser dao.User
	entity.User.Get(dao.Orm(), Order.UserID, &notifyUser)

	var as *result.ActionResult

	weapp_template_msg_data := make(map[string]interface{})
	weapp_template_msg_data["page"] = "pages/user/user"
	weapp_template_msg_data["template_id"] = "bah5ch6kSTi4dvbYzlZ80m7usPIe7PWZEW7Csk_HOy0"
	weapp_template_msg_data["form_id"] = Order.PrepayID
	weapp_template_msg_data["touser"] = notifyUser.OpenID

	data_data := make(map[string]interface{})
	data_data["keyword1"] = map[string]interface{}{"value": notifyUser.Name, "color": "#173177"}
	data_data["keyword2"] = map[string]interface{}{"value": Order.OrderNo, "color": "#173177"}

	var address dao.Address
	json.Unmarshal([]byte(Order.Address), &address)
	addressText := address.Name + "/" + address.ProvinceName + address.CityName + address.CountyName + address.Detail + address.PostalCode + "/" + address.Tel

	data_data["keyword3"] = map[string]interface{}{"value": addressText, "color": "#173177"}

	data_data["keyword4"] = map[string]interface{}{"value": Order.PayTime.Format("2006-01-02 15:04:05"), "color": "#173177"}

	var org dao.Organization
	entity.Organization.Get(dao.Orm(), Order.OID, &org)
	data_data["keyword5"] = map[string]interface{}{"value": org.Name, "color": "#173177"}

	data_data["keyword6"] = map[string]interface{}{"value": strconv.Itoa(int(Order.PayMoney/100)) + "元", "color": "#173177"}

	var Titles = ""
	for _, value := range ogs {
		var goods dao.Goods
		json.Unmarshal([]byte(value.Goods), &goods)
		Titles += goods.Title
	}
	if len(Titles) > 48 {
		Titles = Titles[:48] + "等"
	}
	data_data["keyword7"] = map[string]interface{}{"value": Titles, "color": "#173177"}
	data_data["keyword8"] = map[string]interface{}{"value": "如有疑问，请联系客服", "color": "#173177"}

	weapp_template_msg_data["data"] = data_data

	as = entity.SendWXMessage(weapp_template_msg_data)

	return as
}
func (entity WxService) SendUniformMessage(sendData map[string]interface{}) (*result.ActionResult, int) {

	//gzh := entity.MiniWeb()
	xcx := entity.MiniProgram()

	b, err := json.Marshal(sendData)
	glog.Error(err)

	access_token := entity.GetAccessToken(xcx)
	strReader := strings.NewReader(string(b))
	respones, err := http.Post("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/uniform_send?access_token="+access_token, "application/json", strReader)
	glog.Error(err)
	if err != nil {
		return &result.ActionResult{Code: result.ActionFail, Message: err.Error(), Data: nil}, -1
	}
	defer respones.Body.Close()
	body, err := ioutil.ReadAll(respones.Body)
	glog.Error(err)
	mapData := make(map[string]interface{})
	fmt.Println(string(body))
	err = json.Unmarshal(body, &mapData)
	glog.Error(err)
	if mapData["errcode"] != nil {
		if mapData["errcode"].(float64) == 0 {
			return &result.ActionResult{Code: result.ActionOK, Message: "发送成功", Data: nil}, 0
		}
	}
	return &result.ActionResult{Code: result.ActionFail, Message: mapData["errmsg"].(string), Data: nil}, int(mapData["errcode"].(float64))

}
func (entity WxService) SendWXMessage(sendData map[string]interface{}) *result.ActionResult {
	b, err := json.Marshal(sendData)
	glog.Error(err)

	WxConfig := entity.MiniProgram()

	access_token := entity.GetAccessToken(WxConfig)
	strReader := strings.NewReader(string(b))
	respones, err := http.Post("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token="+access_token, "application/json", strReader)
	glog.Error(err)
	if err != nil {
		return &result.ActionResult{Code: result.ActionFail, Message: err.Error(), Data: nil}
	}
	defer respones.Body.Close()
	body, err := ioutil.ReadAll(respones.Body)
	glog.Error(err)
	mapData := make(map[string]interface{})
	fmt.Println(string(body))
	err = json.Unmarshal(body, &mapData)
	glog.Error(err)
	if mapData["errcode"] != nil {
		if mapData["errcode"].(float64) == 0 {
			return &result.ActionResult{Code: result.ActionOK, Message: "发送成功", Data: nil}
		} else {
			//mapData["errcode"].(float64)
			return &result.ActionResult{Code: result.ActionFail, Message: mapData["errmsg"].(string), Data: nil}
		}
	}
	return &result.ActionResult{Code: result.ActionFail, Message: mapData["errmsg"].(string), Data: nil}

}
func (entity WxService) Order(OrderNo string, title, description string, detail, openid string, IP string, Money uint64, attach string, WxConfig dao.WxConfig) (Success result.ActionResultCode, Message string, wxResult WxOrderResult) {

	//wxConfig := entity.GetWxConfig(WxConfigID)

	mapData := make(util.Map)
	mapData["appid"] = WxConfig.AppID
	mapData["attach"] = attach
	mapData["body"] = title + "-" + description
	mapData["mch_id"] = WxConfig.MchID

	if !strings.EqualFold(detail, "") {
		mapData["detail"] = detail
	}
	//mapData["detail"] = `{ "goods_detail":[ { "goods_id":"iphone6s_16G", "wxpay_goods_id":"1001", "goods_name":"iPhone6s 16G", "quantity":1, "price":528800, "goods_category":"123456", "body":"苹果手机" }, { "goods_id":"iphone6s_32G", "wxpay_goods_id":"1002", "goods_name":"iPhone6s 32G", "quantity":1, "price":608800, "goods_category":"123789", "body":"苹果手机" } ] }`
	mapData["nonce_str"] = tool.UUID()
	mapData["notify_url"] = "https://" + conf.Config.Domain + "/payment/notify"
	mapData["openid"] = openid
	mapData["out_trade_no"] = OrderNo
	mapData["spbill_create_ip"] = IP
	mapData["total_fee"] = strconv.Itoa(int(Money))
	mapData["trade_type"] = "JSAPI"
	mapData["sign_type"] = "MD5"

	list := &collections.ListString{}

	//xml := `<xml>`
	for k, v := range mapData {
		list.Append(k + "=" + v)
		//xml = xml + "<" + k + ">" + v + "</" + k + ">"
	}

	list.SortL()

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + WxConfig.PayKey)
	//fmt.Println(list.Join("&") + "&key=" + self.MiniProgram().PayKey)

	mapData["sign"] = sign

	xmlb, _ := xml.Marshal(&mapData)

	//fmt.Println(sign)
	strReader := strings.NewReader(string(xmlb))

	respones, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "text/xml", strReader)
	glog.Error(err)
	if err != nil {
		return result.ActionFail, err.Error(), WxOrderResult{}
	}

	b, err := ioutil.ReadAll(respones.Body)
	glog.Error(err)
	if err != nil {
		return result.ActionFail, err.Error(), WxOrderResult{}
	}
	//fmt.Println(err)
	//fmt.Println(string(b))

	err = xml.Unmarshal(b, &wxResult)
	if err != nil {
		return result.ActionFail, "支付网关返回结果出错", WxOrderResult{}
	}

	if !strings.EqualFold(wxResult.Return_code, "SUCCESS") {
		//return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionFail, Message: resultXML.Return_msg, Data: nil}}
		return result.ActionFail, wxResult.Return_msg, WxOrderResult{}
	}

	if !strings.EqualFold(wxResult.Result_code, "SUCCESS") {
		//return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionFail, Message: resultXML.Err_code_des, Data: nil}}
		return result.ActionFail, wxResult.Err_code_des, WxOrderResult{}
	}

	return result.ActionOK, "下单成功", wxResult
}
func (entity WxService) MPOrder(OrderNo string, title, description string, ogs []dao.OrdersGoods, openid string, IP string, Money uint64, attach string) (Success result.ActionResultCode, Message string, result WxOrderResult) {

	CostGoodsPrice := uint64(0)

	goods_detail := make([]map[string]interface{}, 0)
	for _, value := range ogs {
		goodsObj := make(map[string]interface{})
		goodsObj["goods_id"] = value.OrdersGoodsNo

		var goods dao.Goods
		json.Unmarshal([]byte(value.Goods), &goods)

		var specification dao.Specification
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

	WxConfig := entity.MiniProgram()

	return entity.Order(OrderNo, title, description, string(detailB), openid, IP, Money, attach, WxConfig)
}

// func (self WxService) GetWxConfig(DB *gorm.DB, CompanyID uint64) *WxConfig {
// 	item := &WxConfig{}
// 	err := DB.Where("CompanyID=?", CompanyID).First(item).Error
// 	glog.Error(err)

// 	if item.ID == 0 {
// 		err = DB.Create(item).Error
// 		glog.Error(err)
// 		return item
// 	} else {
// 		return item
// 	}
// }
func (entity WxService) ChangeWxConfig(DB *gorm.DB, ID uint64, Value dao.WxConfig) error {

	//item := b.GetWxConfig(DB, CompanyID)
	//item.V = Value
	return DB.Model(&dao.WxConfig{}).Where("ID=?", ID).Updates(Value).Error
}

/*func (entity WxService) MiniProgramByAppIDAndMchID(AppID, MchID string) dao.WxConfig {
	var wx dao.WxConfig
	err := dao.Orm().Model(&dao.WxConfig{}).Where("AppID=? and MchID=?", AppID, MchID).First(&wx).Error
	glog.Error(err)
	return wx
}*/
/*func (entity WxService) GetWxConfig(ID uint64) dao.WxConfig {
	var wx dao.WxConfig
	err := dao.Orm().Model(&dao.WxConfig{}).Where("ID=?", ID).First(&wx).Error
	glog.Error(err)
	return wx
}*/
func (entity WxService) MWQRCodeTemp(OID uint64, UserID uint64, qrtype, params string) *result.ActionResult {

	//user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	//context.Request.ParseForm()
	//Page := context.Request.FormValue("Page")
	//Page := context.Request.URL.Query().Get("Page")
	//MyShareKey := tool.Hashids{}.Encode(user.ID)
	wxconfig := entity.MiniWeb()

	access_token := entity.GetAccessToken(wxconfig)

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
		return &result.ActionResult{Code: result.ActionFail, Message: err.Error(), Data: nil}
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &result.ActionResult{Code: result.ActionFail, Message: err.Error(), Data: nil}
	}
	//fmt.Println(string(b))
	defer resp.Body.Close()
	path := tool.WriteTempFile(b, "image/png")
	return &result.ActionResult{Code: result.ActionOK, Message: "", Data: path}

}

/*func (self WxService) WX() WxConfig {

	return WxConfig{AppID: "wx037d3b26b2ba34b2", AppSecret: "c930d5b6a337c6bad9b41556cdcb94d2", Token: "", EncodingAESKey: "", MchID: "1253136001", PayKey: "6af34073b83d6f8a4f35289b92226f20"}
}*/
/*
小程序
*/
func (entity WxService) MiniProgram() dao.WxConfig {

	var wx dao.WxConfig
	wx.AppID = "wx10a413fd12780596"
	wx.AppSecret = "72eb78717d093645510895a89dedaa79"
	wx.MchID = "1253136001"
	wx.PayKey = "6af34073b83d6f8a4f35289b92226f20"
	//err := dao.Orm().Model(&dao.WxConfig{}).Where("OID=? and Type=?", OID, "miniprogram").First(&wx).Error
	//glog.Error(err)
	return wx
}

/*
公众号
*/
func (entity WxService) MiniWeb() dao.WxConfig {

	var wx dao.WxConfig
	wx.AppID = "wx037d3b26b2ba34b2"
	wx.AppSecret = "090938339bac641c9c336e58b118375d"
	//var wx dao.WxConfig
	//err := dao.Orm().Model(&dao.WxConfig{}).Where("OID=? and Type=?", OID, "miniweb").First(&wx).Error
	//glog.Error(err)
	return wx
}

//订单查询
func (entity WxService) OrderQuery(OrderNo string) (Success bool, Result util.Map) {
	var inData = make(util.Map)
	WxConfig := entity.MiniProgram()

	outMap := make(util.Map)
	outMap["mch_id"] = WxConfig.MchID
	outMap["appid"] = WxConfig.AppID
	outMap["nonce_str"] = tool.UUID()
	outMap["out_trade_no"] = OrderNo
	outMap["sign_type"] = "MD5"

	list := &collections.ListString{}
	for k, v := range outMap {
		list.Append(k + "=" + v)
	}
	list.SortL()

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + WxConfig.PayKey)
	outMap["sign"] = sign

	b, err := xml.MarshalIndent(util.Map(outMap), "", "")
	glog.Trace(err)
	//fmt.Println(string(b))

	client := &http.Client{}
	reader := strings.NewReader(string(b))
	response, err := client.Post("https://api.mch.weixin.qq.com/pay/orderquery", "text/xml", reader)
	if err != nil {
		return false, inData
	}
	glog.Trace(err)

	b, err = ioutil.ReadAll(response.Body)
	glog.Trace(err)

	//fmt.Println(string(b))

	err = xml.Unmarshal(b, &inData)
	glog.Trace(err)

	//fmt.Println(inData)

	if strings.EqualFold(inData["return_code"], "SUCCESS") && strings.EqualFold(inData["result_code"], "SUCCESS") && strings.EqualFold(inData["trade_state"], "SUCCESS") {
		Success = true
		Result = inData
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

//查询提现接口
func (entity WxService) GetTransfersInfo(transfers dao.Transfers) (Success bool) {

	WxConfig := entity.MiniProgram()

	outMap := make(util.Map)
	outMap["nonce_str"] = tool.UUID()
	outMap["partner_trade_no"] = transfers.OrderNo
	outMap["mch_id"] = WxConfig.MchID
	outMap["appid"] = WxConfig.AppID

	list := &collections.ListString{}
	for k, v := range outMap {
		list.Append(k + "=" + v)
	}
	list.SortL()

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + WxConfig.PayKey)
	outMap["sign"] = sign

	b, err := xml.MarshalIndent(util.Map(outMap), "", "")
	glog.Trace(err)
	//fmt.Println(string(b))
	//certs, err := tls.LoadX509KeyPair("cert/wxpay/apiclient_cert.pem", "cert/wxpay/apiclient_key.pem")

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
	glog.Trace(err)
	if err != nil {
		return false
	}

	b, err = ioutil.ReadAll(response.Body)
	glog.Trace(err)
	if err != nil {
		return false
	}

	//fmt.Println(string(b))

	var inData = make(util.Map)
	err = xml.Unmarshal(b, &inData)
	glog.Trace(err)
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
func (entity WxService) Transfers(transfers dao.Transfers) (Success bool, Message string) {
	WxConfig := entity.MiniProgram()

	outMap := make(util.Map)
	outMap["mch_appid"] = WxConfig.AppID
	outMap["mchid"] = WxConfig.MchID
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

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + WxConfig.PayKey)
	outMap["sign"] = sign

	b, err := xml.MarshalIndent(util.Map(outMap), "", "")
	glog.Trace(err)
	//fmt.Println(string(b))
	//certs, err := tls.LoadX509KeyPair("cert/wxpay/apiclient_cert.pem", "cert/wxpay/apiclient_key.pem")

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
	glog.Trace(err)

	b, err = ioutil.ReadAll(response.Body)
	glog.Trace(err)

	//fmt.Println(string(b))

	var inData = make(util.Map)
	err = xml.Unmarshal(b, &inData)
	glog.Trace(err)

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
func (entity WxService) CloseOrder(OrderNo string, OID uint64) (Success bool, Message string) {

	WxConfig := entity.MiniProgram()

	outMap := make(util.Map)
	outMap["appid"] = WxConfig.AppID
	outMap["mch_id"] = WxConfig.MchID
	outMap["nonce_str"] = tool.UUID()

	outMap["out_trade_no"] = OrderNo

	list := &collections.ListString{}
	for k, v := range outMap {
		list.Append(k + "=" + v)
	}
	list.SortL()

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + WxConfig.PayKey)
	outMap["sign"] = sign

	b, err := xml.MarshalIndent(util.Map(outMap), "", "")
	glog.Trace(err)

	reader := strings.NewReader(string(b))
	response, err := http.Post("https://api.mch.weixin.qq.com/pay/closeorder", "text/xml", reader)
	glog.Trace(err)

	b, err = ioutil.ReadAll(response.Body)
	glog.Trace(err)

	fmt.Println(string(b))

	var inData = make(util.Map)
	err = xml.Unmarshal(b, &inData)
	glog.Trace(err)
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

//退款
func (entity WxService) Refund(order dao.Orders, ordersPackage dao.OrdersPackage, PayMoney, RefundMoney uint64, Desc string, Type uint64) (Success bool, Message string) {
	WxConfig := entity.MiniProgram()

	//Orders := OrdersService{}
	//op := Orders.GetOrdersPackageByOrderNo(order.OrdersPackageNo)
	//op := Orders.GetOrdersByOrderNo(order.OrdersPackageNo)

	outMap := make(util.Map)
	outMap["appid"] = WxConfig.AppID
	outMap["mch_id"] = WxConfig.MchID
	outMap["nonce_str"] = tool.UUID()

	if strings.EqualFold(order.OrdersPackageNo, "") {
		outMap["out_refund_no"] = order.OrderNo
		outMap["out_trade_no"] = order.OrderNo
		outMap["refund_fee"] = strconv.Itoa(int(order.PayMoney))
		outMap["total_fee"] = strconv.Itoa(int(order.PayMoney))
	} else {

		//op := entity.Orders.GetOrdersPackageByOrderNo(order.OrdersPackageNo)
		//op := Orders.GetOrdersByOrderNo(order.OrdersPackageNo)
		outMap["out_refund_no"] = order.OrderNo
		outMap["out_trade_no"] = order.OrdersPackageNo
		outMap["refund_fee"] = strconv.Itoa(int(order.PayMoney))
		outMap["total_fee"] = strconv.Itoa(int(ordersPackage.TotalPayMoney))
	}

	outMap["refund_desc"] = Desc

	if Type == 0 {
		outMap["refund_account"] = "REFUND_SOURCE_UNSETTLED_FUNDS" //0
	} else {
		outMap["refund_account"] = "REFUND_SOURCE_RECHARGE_FUNDS" //1
	}

	list := &collections.ListString{}
	for k, v := range outMap {
		list.Append(k + "=" + v)
	}
	list.SortL()

	sign := encryption.Md5ByString(list.Join("&") + "&key=" + WxConfig.PayKey)
	outMap["sign"] = sign

	b, err := xml.MarshalIndent(util.Map(outMap), "", "")
	glog.Trace(err)
	//fmt.Println(string(b))
	//certs, err := tls.LoadX509KeyPair("cert/wxpay/apiclient_cert.pem", "cert/wxpay/apiclient_key.pem")

	// Load client cert
	cert, err := tls.LoadX509KeyPair("cert/wxpay/apiclient_cert.pem", "cert/wxpay/apiclient_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile("cert/wxpay/rootca.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		//RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	reader := strings.NewReader(string(b))
	response, err := client.Post("https://api.mch.weixin.qq.com/secapi/pay/refund", "text/xml", reader)
	glog.Trace(err)
	if err != nil {
		Success = false
		Message = err.Error()
		return
	}

	b, err = ioutil.ReadAll(response.Body)
	glog.Trace(err)
	if err != nil {
		Success = false
		Message = err.Error()
		return
	}

	//fmt.Println(string(b))

	var inData = make(util.Map)
	err = xml.Unmarshal(b, &inData)
	glog.Trace(err)
	//fmt.Println(inData)

	if strings.EqualFold(inData["return_code"], "SUCCESS") && strings.EqualFold(inData["result_code"], "SUCCESS") {
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
	}
}
func (entity WxService) Decrypt(encryptedData, session_key, iv_text string) (bool, string) {
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

func (entity WxService) getSHA1(token, timestamp, nonce, encrypt string) string {

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

func (entity WxService) MwGetTicket(WxConfig dao.WxConfig) string {

	if ticketMap[WxConfig.AppID] != nil && (time.Now().Unix()-ticketMap[WxConfig.AppID].Update) < ticketMap[WxConfig.AppID].Expires_in {

		return ticketMap[WxConfig.AppID].Ticket
	}

	url := "http://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + entity.GetAccessToken(WxConfig)

	resp, err := http.Get(url)
	glog.Error(err)

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	d := make(map[string]interface{})

	err = json.Unmarshal(b, &d)
	glog.Error(err)
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
func (entity WxService) MwGetWXJSConfig(url string, OID uint64) map[string]interface{} {

	wxConfig := entity.MiniWeb()

	appId := wxConfig.AppID
	timestamp := time.Now().Unix()
	nonceStr := tool.UUID()
	//chooseWXPay
	list := &collections.ListString{}
	list.Append("noncestr=" + nonceStr)
	list.Append("jsapi_ticket=" + entity.MwGetTicket(wxConfig))
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

//var GlobalWXConfig = dao.WxConfig{CompanyID: -1, AppID: "wx037d3b26b2ba34b2", AppSecret: "fe3faa4e6f8abd87fa4621cb5ed5f725", Token: "30e6e3b03bf7ec6d2ce56a50055e1cd1", EncodingAESKey: "egMWQnCkbuDd7u5GM7EJBnH8mISn5iwAorjRNnFx3dv", MchID: "1342120901"}
