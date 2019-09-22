package wxpay

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/app/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/gweb/tool/collections"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/tool"
)

/*const OpenToken = "JksTRaZE320kujOPZwpfQ6fHIdX3tV718ccg7es0EFY"
const OpenEncodingAesKey = "JTiYSL0cIZTV30Gx7jFfcvNgvGJGEZ4po2YCfceYLIk"
const OpenAppID = "wx0406ef9880e23fdc"
const OpenAppSecret = "04591700ed65e0ebfd95fd4efb948b73"*/

const AppID = "wx037d3b26b2ba34b2"
const AppSecret = "fe3faa4e6f8abd87fa4621cb5ed5f725"
const Token = "30e6e3b03bf7ec6d2ce56a50055e1cd1"
const EncodingAESKey = "egMWQnCkbuDd7u5GM7EJBnH8mISn5iwAorjRNnFx3dv"
const MchID = "1342120901"

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

var accessToken = &AccessToken{}
var ticket = &Ticket{}
var VerifyCache = &struct {
	//ComponentVerifyTicket             string

	Component_access_token            string
	Component_access_token_expires_in int64
	Component_access_token_update     int64

	Pre_auth_code            string
	Pre_auth_code_expires_in int64
	Pre_auth_code_update     int64
}{}

type PushInfo struct {
	AppId                 string `xml:AppId`
	CreateTime            int64  `xml:CreateTime`
	InfoType              string `xml:InfoType`
	ComponentVerifyTicket string `xml:ComponentVerifyTicket`
}

func init() {
	//VerifyCache.ComponentVerifyTicket = "ticket@@@TqYBpfMx2-PyjZuv3L2MKZYiAV2qN5Mf929O8ZlMvPIEOqjbpGATKCDW5VS54yJjUOrk3iLI4y5CQPScvESYQg"

	//Component_access_token := Api_component_token()

	//Api_create_preauthcode(Component_access_token)
}

/*
func Api_query_auth(authorization_code string, ComponentVerifyTicket string) (authorizer_appid, authorizer_access_token, authorizer_refresh_token, func_info string, expires_in int) {

	params := map[string]string{"component_appid": OpenAppID, "authorization_code": authorization_code}

	jd, err := json.Marshal(params)
	glog.Error(err)
	fmt.Println(string(jd))
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.BigEndian, jd)
	resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="+Api_component_token(ComponentVerifyTicket), "application/json", buf)
	glog.Error(err)
	b, err := ioutil.ReadAll(resp.Body)
	glog.Error(err)
	fmt.Println(string(b))

	m := make(map[string]interface{})

	err = json.Unmarshal(b, &m)
	glog.Error(err)

	if m["authorization_info"] != nil {
		authorization_info := m["authorization_info"].(map[string]interface{})

		authorizer_appid = authorization_info["authorizer_appid"].(string)
		authorizer_access_token = authorization_info["authorizer_access_token"].(string)
		expires_in, _ = strconv.Atoi(strconv.FormatFloat(authorization_info["expires_in"].(float64), 'f', 0, 64))
		authorizer_refresh_token = authorization_info["authorizer_refresh_token"].(string)
		//func_info = authorization_info["func_info"].([]interface{})
	}

	return
}
*/

/*func Api_component_token(ComponentVerifyTicket string) string {
	if time.Now().Unix()-VerifyCache.Component_access_token_update >= VerifyCache.Component_access_token_expires_in-10 || strings.EqualFold(VerifyCache.Component_access_token, "") {

		params := map[string]string{"component_appid": OpenAppID, "component_appsecret": OpenAppSecret, "component_verify_ticket": ComponentVerifyTicket}

		jd, err := json.Marshal(params)
		glog.Error(err)
		fmt.Println(string(jd))
		buf := bytes.NewBuffer(make([]byte, 0))
		binary.Write(buf, binary.BigEndian, jd)
		resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/component/api_component_token", "application/json", buf)
		glog.Error(err)
		b, err := ioutil.ReadAll(resp.Body)
		glog.Error(err)
		fmt.Println(string(b))

		var respData = &struct {
			Component_access_token string `json:"component_access_token"`
			Expires_in             int64  `json:"expires_in"`
		}{}

		err = json.Unmarshal(b, respData)
		glog.Error(err)

		VerifyCache.Component_access_token = respData.Component_access_token
		VerifyCache.Component_access_token_expires_in = respData.Expires_in
		VerifyCache.Component_access_token_update = time.Now().Unix()

		return VerifyCache.Component_access_token

	} else {
		return VerifyCache.Component_access_token
	}
}*/

/*func Api_create_preauthcode(component_access_token string) string {
	if time.Now().Unix()-VerifyCache.Pre_auth_code_update >= VerifyCache.Pre_auth_code_expires_in-10 || strings.EqualFold(VerifyCache.Pre_auth_code, "") {

		params := map[string]string{"component_appid": OpenAppID}
		jd, err := json.Marshal(params)
		glog.Error(err)
		fmt.Println(string(jd))
		buf := bytes.NewBuffer(make([]byte, 0))
		binary.Write(buf, binary.BigEndian, jd)
		resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token="+component_access_token, "application/json", buf)
		glog.Error(err)
		b, err := ioutil.ReadAll(resp.Body)
		glog.Error(err)
		fmt.Println(string(b))

		var respData = &struct {
			Pre_auth_code string `json:"pre_auth_code"`
			Expires_in    int64  `json:"expires_in"`
		}{}

		err = json.Unmarshal(b, respData)
		glog.Error(err)

		VerifyCache.Pre_auth_code = respData.Pre_auth_code
		VerifyCache.Pre_auth_code_expires_in = respData.Expires_in
		VerifyCache.Pre_auth_code_update = time.Now().Unix()
		fmt.Println(respData)

		return VerifyCache.Pre_auth_code
	} else {
		return VerifyCache.Pre_auth_code
	}
}*/
func GetAccessToken() string {

	if (time.Now().Unix() - accessToken.Update) < accessToken.Expires_in {

		return accessToken.Access_token
	}

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + AppID + "&secret=" + AppSecret

	resp, err := http.Get(url)
	glog.Error(err)

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	d := make(map[string]interface{})

	err = json.Unmarshal(b, &d)
	glog.Error(err)
	fmt.Println(string(b))
	fmt.Println(d)

	accessToken.Access_token = d["access_token"].(string)
	accessToken.Expires_in = int64(d["expires_in"].(float64))
	accessToken.Update = time.Now().Unix()

	return accessToken.Access_token
}

func GetTicket() string {

	if (time.Now().Unix() - ticket.Update) < ticket.Expires_in {

		return ticket.Ticket
	}

	url := "http://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + GetAccessToken()

	resp, err := http.Get(url)
	glog.Error(err)

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	d := make(map[string]interface{})

	err = json.Unmarshal(b, &d)
	glog.Error(err)
	fmt.Println(string(b))
	fmt.Println(d)

	ticket.Ticket = d["ticket"].(string)
	ticket.Expires_in = int64(d["expires_in"].(float64))
	ticket.Update = time.Now().Unix()

	return ticket.Ticket
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

func OrderJS(OrderNo string, ShopName string, DetailJSON string, Host string, openid string, IP string, TotalMoney uint64) {

	//https://api.mch.weixin.qq.com/pay/unifiedorder
	postXml := `<xml>
   <appid>` + AppID + `</appid>
   <body>` + ShopName + "-服务/产品" + `</body>
   <mch_id>` + MchID + `</mch_id>
   <detail><![CDATA[` + DetailJSON + `]]></detail>
   <nonce_str>` + tool.UUID() + `</nonce_str>
   <notify_url>` + Host + `/wx/notify` + `</notify_url>
   <openid>` + openid + `</openid>
   <out_trade_no>` + OrderNo + `</out_trade_no>
   <spbill_create_ip>` + IP + `</spbill_create_ip>
   <total_fee>` + strconv.FormatUint(TotalMoney, 10) + `</total_fee>
   <trade_type>JSAPI</trade_type>
   <sign>0CB01533B8C1EF103065174F50BCA001</sign>
</xml>`
	fmt.Println(postXml)

}
func GetWXJSConfig(url string) (appId string, timestamp int64, nonceStr string, signature string) {
	appId = AppID
	timestamp = time.Now().Unix()
	nonceStr = tool.UUID()
	//chooseWXPay
	list := &collections.ListString{}
	list.Append("noncestr=" + nonceStr)
	list.Append("jsapi_ticket=" + GetTicket())
	list.Append("timestamp=" + strconv.FormatInt(timestamp, 10))

	_url := strings.Split(url, "#")[0]
	list.Append("url=" + _url)
	list.SortL()

	signature = util.SignSha1(list.Join("&"))

	return
}
