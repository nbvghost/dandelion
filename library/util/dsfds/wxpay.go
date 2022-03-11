package wxpay

/*const OpenToken = "JksTRaZE320kujOPZwpfQ6fHIdX3tV718ccg7es0EFY"
const OpenEncodingAesKey = "JTiYSL0cIZTV30Gx7jFfcvNgvGJGEZ4po2YCfceYLIk"
const OpenAppID = "wx0406ef9880e23fdc"
const OpenAppSecret = "04591700ed65e0ebfd95fd4efb948b73"*/

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
/*func GetAccessToken() string {

	if (time.Now().Unix() - accessToken.Update) < accessToken.Expires_in {

		return accessToken.Access_token
	}

	WxConfig := server.WxConfigService{}.GlobalWXAConfig()

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + WxConfig.AppID + "&secret=" + WxConfig.AppSecret

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
}*/

/*func GetTicket() string {

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
*/

/**/

/*

func OrderQuery(OrderNo uint) (return_code, result_code, trade_state, time_end string, total_fee uint) {
	WxConfig := server.WxConfigService{}.GlobalWXAConfig()

	outMap := make(map[string]string)
	outMap["appid"] = WxConfig.AppID
	outMap["mch_id"] = WxConfig.MchID
	outMap["nonce_str"] = tool.UUID()
	outMap["out_trade_no"] = strconv.Itoa(int(OrderNo))

	list := &tool.List{}
	for k, v := range outMap {
		list.Append(k + "=" + v)
	}
	list.SortL()

	sign := util.Md5ByString(list.Join("&") + "&key=" + WxConfig.PayKey)
	outMap["sign"] = sign

	b, err := xml.MarshalIndent(util.Map(outMap), "", "")
	glog.Trace(err)
	//fmt.Println(string(b))

	reader := strings.NewReader(string(b))
	response, err := http.Post("https://api.mch.weixin.qq.com/pay/orderquery", "text/xml", reader)
	glog.Trace(err)

	b, err = ioutil.ReadAll(response.Body)
	glog.Trace(err)
	//fmt.Println(string(b))

	inMap := make(util.Map)
	err = xml.Unmarshal(b, &inMap)
	glog.Trace(err)
	fmt.Println(inMap)

	return_code = inMap["return_code"]
	result_code = inMap["result_code"]
	trade_state = inMap["trade_state"]
	time_end = inMap["time_end"]
	total_fee, _ = strconv.ParseUint(inMap["total_fee"], 10, 64)

	return
}


func GetWXJSConfig(url string) (appId string, timestamp int64, nonceStr string, signature string) {
	WxConfig := server.WxConfigService{}.GlobalWXAConfig()

	appId = WxConfig.AppID
	timestamp = time.Now().Unix()
	nonceStr = tool.UUID()
	//chooseWXPay
	list := &tool.List{}
	list.Append("noncestr=" + nonceStr)
	list.Append("jsapi_ticket=" + GetTicket())
	list.Append("timestamp=" + strconv.FormatInt(timestamp, 10))

	_url := strings.Split(url, "#")[0]
	list.Append("url=" + _url)
	list.SortL()

	signature = util.SignSha1(list.Join("&"))

	return
}*/
