package baidu

type Index struct {

}
type translateInfo struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func (m *Index) _TranslateBaidu(query, from, to string) (list []translateInfo, err error) {
	//var translate model.Translate

	/*tx := db.Orm().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	tx.Model(model.Translate{}).Where(`"Text"=? and "LangType"=?`, query, to).First(&translate)
	if !translate.IsZero() {
		return translate.LangText, nil
	}*/

	//------------

	/*salt := fmt.Sprintf("%d", time.Now().Unix())
	sign := strings.ToLower(encryption.Md5ByString(fmt.Sprintf("%s%s%s%s", appid, query, salt, securityKey)))
	postParams := url.Values{}
	//q := url.QueryEscape(query)
	postParams.Set("q", query)            //	string	是	请求翻译query	UTF-8编码
	postParams.Set("from", from)          //	string	是	翻译源语言	可设置为auto
	postParams.Set("to", m.baiduCode[to]) //	string	是	翻译目标语言	不可设置为auto
	postParams.Set("appid", appid)        //	string	是	APPID	可在管理控制台查看
	postParams.Set("salt", salt)          //	string	是	随机数	可为字母或数字的字符串
	postParams.Set("sign", sign)          //	string	是	签名	appid+q+salt+密钥的MD5值

	response, err := http.PostForm(tranUrl, postParams)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var r result
	if err = json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	if len(r.ErrorCode) > 0 {
		return nil, fmt.Errorf(r.ErrorMsg)
	}*/

	//------------

	//translate.Text = query
	//translate.LangType = to

	/*var translateText string
	var isGap bool
	for _, v := range r.TransResult {
		if strings.EqualFold(v.Src, "###") {
			translateText = translateText + "\n###\n"
			isGap = true
		} else {
			if translateText == "" {
				translateText = v.Dst
			} else {
				if isGap {
					translateText = translateText + v.Dst
				} else {
					translateText = translateText + "\n" + v.Dst
				}
				isGap = false
			}
		}
	}*/
	//translate.LangText = r.TransResult[0].Dst

	/*if err = tx.Model(model.Translate{}).Create(&translate).Error; err != nil {
		return "", err
	}*/
	return nil, nil
}