package wx

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/gpa/types"
)

type Token struct {
	Get struct {
		CompanyID types.PrimaryKey `uri:"CompanyID"`
	} `method:"Get"`
}

func (m *Token) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	/* if strings.EqualFold(context.Request.Method, "GET") {

		signature := context.Request.URL.Query().Get("signature")
		timestamp := context.Request.URL.Query().Get("timestamp")
		nonce := context.Request.URL.Query().Get("nonce")
		echostr := context.Request.URL.Query().Get("echostr")
		fmt.Println(signature, timestamp, nonce, echostr)

		CompanyID, _ := strconv.ParseUint(context.PathParams["CompanyID"], 10, 64)
		wxConfig := service.WxConfig.GetWxConfig(service.Orm, CompanyID)

		list := &tool.List{}
		list.Append(wxConfig.Token)
		list.Append(timestamp)
		list.Append(nonce)
		list.SortL()

		_signature := util.SignSha1(list.Join(""))
		if strings.EqualFold(signature, _signature) {
			return &gweb.TextResult{Data: echostr}
		} else {
			return &gweb.TextResult{Data: "ERROR"}
		}

	} */

	/*//encrypt_type=aes&msg_signature=0da847abf28d7fbd1e5f131c5e2c5045f0b7618d&nonce=289158306&signature=8cbe9a765ba9c6fd383c9abcd68cde7234f438cf&timestamp=1510642666
	encrypt_type := context.Request.URL.Query().Get("encrypt_type")
	msg_signature := context.Request.URL.Query().Get("msg_signature")
	nonce := context.Request.URL.Query().Get("nonce")
	//signature:=context.Request.URL.Query().Get("signature")
	timestamp := context.Request.URL.Query().Get("timestamp")
	fmt.Println(context.Request.URL.Query().Encode())

	b, err := ioutil.ReadAll(context.Request.Body)
	log.Println(err)
	fmt.Println("dsfs", string(b))

	if strings.EqualFold(encrypt_type, "aes") {
		tokenXML := &wxpay.TokenXML{}

		xml.Unmarshal(b, tokenXML)

		sdfd, content := wxpay.DecryptMsg(msg_signature, timestamp, nonce, tokenXML.Encrypt)
		if sdfd {
			pushInfo := &wxpay.PushInfo{}
			xml.Unmarshal([]byte(content), pushInfo)
			//fmt.Println(pushInfo)
			//wxpay.VerifyCache.ComponentVerifyTicket = pushInfo.ComponentVerifyTicket
			//fmt.Println(wxpay.VerifyCache)

			err := service.Configuration.ChangeConfiguration(service.Orm, play.ConfigurationKey_component_verify_ticket, pushInfo.ComponentVerifyTicket)
			log.Println(err)

		}

	}*/
	return &result.TextResult{Data: "success"}, nil
	/*signature := context.Request.URL.Query().Get("signature")
	echostr := context.Request.URL.Query().Get("echostr")

	nonce := context.Request.URL.Query().Get("nonce")
	timestamp := context.Request.URL.Query().Get("timestamp")*/

	/*list := &tools.List{}
	list.Append(data.WXConfig.Token)
	list.Append(timestamp)
	list.Append(nonce)
	list.SortL()

	if strings.EqualFold(util.SignSha1(list.Join("")), signature) {

	} else {

	}
	return &web.TextResult{""}*/
}
