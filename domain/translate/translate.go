package translate

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/nbvghost/dandelion/service"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/html"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/tool/object"
	"github.com/pkg/errors"

	alimt20181012 "github.com/alibabacloud-go/alimt-20181012/v2/client"
)

type Html struct {
	LanguageCode    map[string]string
	AccessKeyID     string
	AccessKeySecret string
	sync.RWMutex
}

func (m *Html) html(node *html.Node) ([]byte, error) {
	var buf bytes.Buffer
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		err := html.Render(&buf, c)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil

}

type nodeID string

func newNodeID(seqID int) nodeID {
	return nodeID(fmt.Sprintf("node%d", seqID))
}

var varRegexp = regexp.MustCompile(`\$\{[\S\x20]+}`)

type translateInfo struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func (m *Html) TranslateHtml(context constrain.IContext, docBytes []byte) ([]byte, error) {
	contextValue := contexext.FromContext(context)
	var err error
	var node *html.Node
	node, err = html.Parse(bytes.NewBuffer(docBytes))
	if err != nil {
		return nil, err
	}

	//todo 当前暂定en为内容的默认语言,后面走数据库
	if strings.EqualFold(contextValue.Lang, "en") {
		return docBytes, nil
	}

	if _, ok := m.LanguageCode[contextValue.Lang]; !ok {
		//当前语言不在翻译列表中
		return docBytes, nil
	}

	var seqID int
	setMap := map[string]nodeID{}
	willTranslateTexts := map[string]translateInfo{}
	//var willTranslateLocker sync.RWMutex

	extractFunc := func(text string) string {
		texts := strings.Split(text, "\n")
		for i := range texts {
			if varRegexp.MatchString(texts[i]) {
				continue
			}
			v := strings.TrimSpace(texts[i])
			var nid nodeID
			var ok bool
			if nid, ok = setMap[v]; !ok {
				seqID++
				nid = newNodeID(seqID)
				setMap[v] = nid
				willTranslateTexts[string(nid)] = translateInfo{
					Src: v,
				}
			}
			texts[i] = fmt.Sprintf("{{.%s.Dst}}", string(nid))
		}
		return strings.Join(texts, "\n")
	}

	//提取要翻译的文字
	var f func(*html.Node)
	f = func(n *html.Node) {
		var noTranslate bool
		for _, v := range n.Attr {
			if strings.EqualFold(v.Key, "no-translate") {
				noTranslate = true
				break
			}
		}
		if noTranslate {
			return
		}

		if n.Type == html.TextNode && !strings.EqualFold(n.Parent.Data, "style") && !strings.EqualFold(n.Parent.Data, "script") {
			text := strings.TrimSpace(n.Data)
			if len(text) > 0 {
				n.Data = extractFunc(text)
			}
		}
		if strings.EqualFold(n.Data, "input") {
			for i, v := range n.Attr {
				if strings.EqualFold(v.Key, "placeholder") {
					n.Attr[i].Val = extractFunc(n.Attr[i].Val)
					break
				}
			}
		}
		if strings.EqualFold(n.Data, "html") && n.Type == html.ElementNode {
			if len(n.Attr) > 0 {
				var hasLangAttr bool
				for i, v := range n.Attr {
					if strings.EqualFold(v.Key, "lang") {
						n.Attr[i].Val = contextValue.Lang
						hasLangAttr = true
						break
					}
				}
				if !hasLangAttr {
					n.Attr = append(n.Attr, html.Attribute{
						Key: "lang",
						Val: contextValue.Lang,
					})
				}
			} else {
				n.Attr = append(n.Attr, html.Attribute{
					Key: "lang",
					Val: contextValue.Lang,
				})
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	tx := db.Orm().Begin(&sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var translateList []string

	//TODO 缓存在cache server 里，或者放在redis
	var translateModelList []model.Translate
	tx.Model(model.Translate{}).Where(`"TextType"=? and "LangType"=?`, "en", contextValue.Lang).Find(&translateModelList)

	for k, v := range willTranslateTexts {
		var has bool
		//todo 优化二分查找或者其它方法
		for _, e := range translateModelList {
			if strings.EqualFold(v.Src, e.Text) {
				v.Dst = e.LangText
				willTranslateTexts[k] = v
				has = true
				break
			}
		}
		if !has {
			translateList = append(translateList, v.Src)
		}
	}

	if len(translateList) > 0 {
		translateMap, err := m.Translate(translateList, "en", contextValue.Lang)
		if err != nil {
			context.Logger().Error("translate error", zap.Error(err))
			return nil, err
		}

		for index := range translateList {
			v := translateList[index]
			translateText := translateMap[index]
			for k, vv := range willTranslateTexts {
				if vv.Src == v {
					vv.Dst = translateText
					willTranslateTexts[k] = vv
					break
				}
			}
			if len(translateText) == 0 {
				translateText = v
			}
			if err = tx.Model(model.Translate{}).Create(&model.Translate{
				Text:     v,
				TextType: "en",
				LangType: contextValue.Lang,
				LangText: translateText,
			}).Error; err != nil {
				context.Logger().Error("write db.Translate error", zap.Error(err))
				return nil, err
			}
		}
	}

	/*//var translateText string
	for _, v := range translateList {
		if len(v) > 0 {

			var translateText string
			translateText, err = m.Translate(v, "en", contextValue.Lang)
			if err != nil {
				context.Logger().Error("translate error", zap.Error(err))
				return nil, err
			}

			//willTranslateLocker.Lock()
			for k, vv := range willTranslateTexts {
				if vv.Src == v {
					vv.Dst = translateText
					willTranslateTexts[k] = vv
					break
				}
			}
			//willTranslateLocker.Unlock()

			if len(translateText) == 0 {
				translateText = v
			}

			if err = tx.Model(model.Translate{}).Create(&model.Translate{
				Text:     v,
				TextType: "en",
				LangType: contextValue.Lang,
				LangText: translateText,
			}).Error; err != nil {
				context.Logger().Error("write db.Translate error", zap.Error(err))
				return nil, err
			}
		}
	}*/

	docBytes, err = m.html(node)
	if err != nil {
		return nil, err
	}

	var t *template.Template
	t, err = template.New("default").Parse(string(docBytes))
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(nil)
	err = t.Execute(buffer, willTranslateTexts)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

var bingTranslateHttpParamsReg = regexp.MustCompile(`var\sparams_RichTranslateHelper.*?;`)
var bingTranslateHttpParamsValueReg = regexp.MustCompile(`\[.?]`)
var bingTranslateHttpIGReg = regexp.MustCompile(`,"ig":"(.*?)",`)

// [555,"sdfds"]
var bingParams = struct {
	IG      string
	Token   string
	key     int64
	Timeout int64
	Client  *http.Client
	Cookies []*http.Cookie
}{}

func (m *Html) bingTranslateParams() error {
	if bingParams.Client == nil {
		bingParams.Client = &http.Client{}
	}
	if time.Unix((bingParams.key+bingParams.Timeout)/1000, 0).After(time.Now()) {
		return nil
	}
	response, err := http.Get("https://cn.bing.com/translator")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bingParams.Cookies = response.Cookies()
	doc, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if bingTranslateHttpParamsReg.Match(doc) {
		params := bingTranslateHttpParamsReg.Find(doc)
		valueJson := bingTranslateHttpParamsValueReg.Find(params)
		var values []any
		if err = json.Unmarshal(valueJson, &values); err != nil {
			return err
		}
		if len(values) > 0 {
			bingParams.key = int64(object.ParseInt(values[0]))
		}
		if len(values) > 1 && values[1] != nil {
			v, ok := values[1].(string)
			if ok {
				bingParams.Token = v
			}
		}
		if len(values) > 2 && values[2] != nil {
			bingParams.Timeout = int64(object.ParseInt(values[2]))
		}

		var ig string
		igParams := bingTranslateHttpIGReg.FindSubmatch(doc)
		if len(igParams) >= 2 {
			ig = string(igParams[1])
		}
		bingParams.IG = ig
	}
	if len(bingParams.IG) == 0 || bingParams.key == 0 || len(bingParams.Token) == 0 {
		return fmt.Errorf("获取参数错误")
	}
	return nil
}

type bingResult struct {
	DetectedLanguage struct {
		Language string  `json:"language"`
		Score    float64 `json:"score"`
	} `json:"detectedLanguage"`
	Translations []struct {
		Text    string `json:"text"`
		To      string `json:"to"`
		SentLen struct {
			SrcSentLen   []int `json:"srcSentLen"`
			TransSentLen []int `json:"transSentLen"`
		} `json:"sentLen"`
	} `json:"translations"`
}

func (m *Html) bingTranslate(query, from, to string) (_ string, err error) {
	if err = m.bingTranslateParams(); err != nil {
		return "", err
	}

	var translate model.Translate

	/*tx := db.Orm().Begin(&sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var translate model.Translate
	tx.Model(model.Translate{}).Where(`"Text"=? and "LangType"=?`, query, to).First(&translate)
	if !translate.IsZero() {
		return translate.LangText, nil
	}*/

	postParams := url.Values{}
	//postParams.Set("fromLang", "en")
	postParams.Set("fromLang", "auto-detect")
	postParams.Set("text", query)
	postParams.Set("to", to)
	postParams.Set("token", bingParams.Token)
	postParams.Set("key", fmt.Sprintf("%d", bingParams.key))

	bingUrl := fmt.Sprintf("https://cn.bing.com/ttranslatev3?isVertical=1&&IG=%s&IID=translator", bingParams.IG)
	//c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	req, err := http.NewRequest("POST", bingUrl, strings.NewReader(postParams.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("origin", "https: //cn.bing.com")
	req.Header.Set("referer", "https: //cn.bing.com/translator")
	//req.Header.Set("cookie", "MUID=22F939AB92B66021282728BA93646187; MUIDB=22F939AB92B66021282728BA93646187; _EDGE_V=1; SRCHD=AF=NOFORM; SRCHUID=V=2&GUID=2082497D61D044048F6FB0F98A4391A6&dmnchg=1; MUIDV=NU=1; imgv=flts=20220221; ABDEF=V=13&ABDV=11&MRNB=1648784738445&MRB=0; _FP=hta=off; _UR=QS=0&TQS=0; NAP=V=1.9&E=1a64&C=fA90D5o3AOwzhmRxrZO1Y2jp-2ZO2i8m6Mthb844ZnunCn_Msfqbfw&W=1; ZHCHATSTRONGATTRACT=TRUE; SUID=M; _EDGE_S=SID=1FD97DDB65F96D0817C56CA5642B6CA2; _SS=SID=1FD97DDB65F96D0817C56CA5642B6CA2; SRCHUSR=DOB=20211218&T=1649259507000&TPC=1649234246000; ipv6=hit=1649263108481&t=4; ZHCHATWEAKATTRACT=TRUE; btstkn=7dbKx%252BXuDz1K1UuPmde1f8Vfvm4ztAJzKmECLUhgCjPa11JXX1MNrl7D9A58LsC51lKU4IaxSj7K1Q0RZ9AcK2eLd23C5eXBQWAZSpOX7Ug%253D; _tarLang=default=my; _TTSS_OUT=hist=WyJsemgiLCJ6aC1IYW5zIiwibXkiXQ==; _TTSS_IN=hist=WyJ6aC1IYW5zIiwiZXMiLCJlbiIsImF1dG8tZGV0ZWN0Il0=; _HPVN=CS=eyJQbiI6eyJDbiI6NDksIlN0IjoyLCJRcyI6MCwiUHJvZCI6IlAifSwiU2MiOnsiQ24iOjQ5LCJTdCI6MCwiUXMiOjAsIlByb2QiOiJIIn0sIlF6Ijp7IkNuIjo0OSwiU3QiOjEsIlFzIjowLCJQcm9kIjoiVCJ9LCJBcCI6dHJ1ZSwiTXV0ZSI6dHJ1ZSwiTGFkIjoiMjAyMi0wNC0wNlQwMDowMDowMFoiLCJJb3RkIjowLCJHd2IiOjAsIkRmdCI6bnVsbCwiTXZzIjowLCJGbHQiOjAsIkltcCI6MjgyfQ==; SNRHOP=I=&TS=; SRCHHPGUSR=SRCHLANG=zh-Hans&BRW=XW&BRH=M&CW=2012&CH=948&SW=2048&SH=1080&DPR=1&UTC=480&DM=0&WTS=63784856307&HV=1649261953&BZA=0")

	for i := range bingParams.Cookies {
		req.AddCookie(bingParams.Cookies[i])
	}

	response, err := bingParams.Client.Do(req)
	//response, err := bingParams.Client.PostForm(fmt.Sprintf("https://cn.bing.com/ttranslatev3?isVertical=1&&IG=%s&IID=translator", bingParams.IG), postParams)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var obj = struct {
		StatusCode   int    `json:"statusCode"`
		ErrorMessage string `json:"errorMessage"`
		ShowCaptcha  bool   `json:"ShowCaptcha"`
	}{}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return "", err
	}

	if obj.StatusCode > 0 {
		return "", errors.New(obj.ErrorMessage)
	}
	if obj.ShowCaptcha {
		time.Sleep(time.Second)
		return m.bingTranslate(query, from, to)
	}

	var rs []bingResult
	if err = json.Unmarshal(b, &rs); err != nil {
		return "", err
	}
	if len(rs) == 0 {
		return "", fmt.Errorf("translate eror")
	}
	r := rs[0]

	if len(r.Translations) == 0 {
		return "", errors.Errorf("translate eror")
	}

	translate.Text = query
	translate.LangType = to
	translate.LangText = r.Translations[0].Text
	/*if err = tx.Model(model.Translate{}).Create(&translate).Error; err != nil {
		return "", err
	}*/
	return translate.LangText, nil

}

type _Result struct {
	TranslatedText string `json:"translatedText"`
	Error          string `json:"error"`
}

var _aliyunClient *alimt20181012.Client

func (m *Html) GetClient() (*alimt20181012.Client, error) {
	if _aliyunClient == nil {
		// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
		// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
		config := &openapi.Config{
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
			AccessKeyId: tea.String(m.AccessKeyID),
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
			AccessKeySecret: tea.String(m.AccessKeySecret),
		}
		// Endpoint 请参考 https://api.aliyun.com/product/alimt
		config.Endpoint = tea.String("mt.aliyuncs.com")

		_result, _err := alimt20181012.NewClient(config)
		if _err != nil {
			return nil, _err
		}
		_aliyunClient = _result
		return _aliyunClient, nil
	}
	return _aliyunClient, nil
}
func (m *Html) Translate(query []string, from, to string) (map[int]string, error) {
	translateMap := make(map[int]string)
	samllMap := make(map[string]string)
	samllLen := 0
	for i := range query {
		if len(query[i]) > 1000 {
			translateText, err := m.translateBase(query[i], from, to)
			if err != nil {
				return nil, err
			}
			translateMap[i] = translateText
		} else {
			if samllLen+len(query[i]) > 8000 || len(samllMap) >= 50 {
				mTranslateMap, err := m.translateBatchBase(samllMap, from, to)
				if err != nil {
					return nil, err
				}
				for i2, s := range mTranslateMap {
					translateMap[i2] = s
				}
				samllMap = make(map[string]string)
				samllLen = 0
			} else {
				samllMap[object.ParseString(i)] = query[i]
				samllLen = samllLen + len(query[i])
			}
		}
	}
	return translateMap, nil
}
func (m *Html) translateBase(query string, from, to string) (string, error) {
	client, err := m.GetClient()
	if err != nil {
		return "", err
	}
	getBatchTranslateRequest := &alimt20181012.TranslateGeneralRequest{
		FormatType:     tea.String("text"),
		Scene:          tea.String("general"),
		SourceLanguage: tea.String(from),
		SourceText:     tea.String(query),
		TargetLanguage: tea.String(to),
	}
	runtime := &util.RuntimeOptions{}
	res, err := client.TranslateGeneralWithOptions(getBatchTranslateRequest, runtime)
	if err != nil {
		return "", err
	}
	if *res.StatusCode != 200 {
		return "", errors.New("网络错误")
	}
	if *res.Body.Code != 200 {
		return "", errors.New(tea.StringValue(res.Body.Message))
	}
	return tea.StringValue(res.Body.Data.Translated), nil
}
func (m *Html) translateBatchBase(query map[string]string, from, to string) (map[int]string, error) {
	outArr := make(map[int]string)
	client, err := m.GetClient()
	if err != nil {
		return outArr, err
	}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return outArr, err
	}
	getBatchTranslateRequest := &alimt20181012.GetBatchTranslateRequest{
		ApiType:        tea.String("translate_standard"),
		FormatType:     tea.String("text"),
		Scene:          tea.String("general"),
		SourceLanguage: tea.String(from),
		SourceText:     tea.String(string(queryJson)),
		TargetLanguage: tea.String(to),
	}
	runtime := &util.RuntimeOptions{}
	res, err := client.GetBatchTranslateWithOptions(getBatchTranslateRequest, runtime)
	if err != nil {
		return outArr, err
	}
	if *res.StatusCode != 200 {
		return outArr, errors.New("网络错误")
	}
	if *res.Body.Code != 200 {
		return outArr, errors.New(tea.StringValue(res.Body.Message))
	}

	trans := res.Body.TranslatedList
	for i := range trans {
		item := trans[i]
		index := object.ParseInt(item["index"])
		translated := object.ParseString(item["translated"])
		outArr[index] = translated
	}
	return outArr, nil
}

/*
	func (m *Html) Translate(query, from, to string) (string, error) {
		//###
		//POST http://translate.app.usokay.com/translate
		//Content-Type: application/json
		//
		//{
		//	"q": "name",
		//	"source": "en",
		//	"target": "zh",
		//	"format": "text",
		//	"alternatives": 0,
		//	"api_key": "ba07e09c-6e8c-4c1f-b3e0-88091934d51f"
		//}

		postParams := make(map[string]any)
		//q := url.QueryEscape(query)
		postParams["q"] = query
		postParams["source"] = from
		postParams["target"] = to
		postParams["format"] = "text"
		postParams["alternatives"] = 0
		postParams["api_key"] = m.ApiKey

		postParamsBytes, err := json.Marshal(&postParams)
		if err != nil {
			return "", err
		}
		response, err := http.Post("http://translate.app.usokay.com/translate", "application/json", bytes.NewReader(postParamsBytes))
		if err != nil {
			return "", err
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()

		var r _Result
		if err = json.Unmarshal(body, &r); err != nil {
			return "", err
		}
		if response.StatusCode != 200 {
			return "", errors.New(r.Error)
		}
		return r.TranslatedText, nil
	}
*/
func (m *Html) _TranslateBaidu(query, from, to string) (list []translateInfo, err error) {
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

func NewTranslateHtml(languageCode map[string]string) *Html {
	apiKey := service.Configuration.GetAliyunConfiguration(0)
	return &Html{LanguageCode: languageCode, AccessKeyID: apiKey.AccessKeyID, AccessKeySecret: apiKey.AccessKeySecret}
}
