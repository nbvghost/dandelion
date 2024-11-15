package wechat

type MiniSecureKey struct {
	AppID     string
	AppSecret string
}
type MiniApp struct {
	MiniSecureKey
	MchID                      string //= "190000****"                                // 商户号
	MchCertificateSerialNumber string //= "3775B6A45ACD588826D15E583A95F5DD********"  // 商户证书序列号
	MchAPIv2Key                string //= "2ab9****************************"          // 商户APIv3密钥
	MchAPIv3Key                string //= "2ab9****************************"          // 商户APIv3密钥
}
type MiniWeb struct {
	MiniSecureKey
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
