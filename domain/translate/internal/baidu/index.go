package baidu

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"github.com/nbvghost/tool/encryption"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Index struct {
	config *serviceargument.BaiduTranslateConfig
}

func New() *Index {
	return &Index{}
}

type Result struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`

	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (m *Index) translateBase(query string, from, to string) (*Result, error) {
	if m.config == nil {
		m.config = service.Configuration.GetBaiduTranslateConfiguration(0)
	}

	salt := fmt.Sprintf("%d", time.Now().Unix())
	sign := strings.ToLower(encryption.Md5ByString(fmt.Sprintf("%s%s%s%s", m.config.AppID, query, salt, m.config.AppKey)))
	postParams := url.Values{}
	//q := url.QueryEscape(query)
	postParams.Set("q", query)              //	string	是	请求翻译query	UTF-8编码
	postParams.Set("from", from)            //	string	是	翻译源语言	可设置为auto
	postParams.Set("to", to)                //	string	是	翻译目标语言	不可设置为auto
	postParams.Set("appid", m.config.AppID) //	string	是	APPID	可在管理控制台查看
	postParams.Set("salt", salt)            //	string	是	随机数	可为字母或数字的字符串
	postParams.Set("sign", sign)            //	string	是	签名	appid+q+salt+密钥的MD5值

	response, err := http.PostForm("https://fanyi-api.baidu.com/api/trans/vip/translate", postParams)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var r Result
	if err = json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	if len(r.ErrorCode) > 0 {
		return nil, fmt.Errorf(r.ErrorMsg)
	}

	return &r, nil
}
func (m *Index) Translate(query []string, from, to string) (map[int]string, error) {

	translateMap := make(map[int]string)
	smallArray := make([]string, 0)
	smallLen := 0

	for i := range query {
		if strings.Contains(query[i], "\n") {
			mTranslateMap, err := m.translateBase(query[i], from, to)
			if err != nil {
				return nil, err
			}
			kk := make([]string, 0)
			for _, s := range mTranslateMap.TransResult {
				kk = append(kk, s.Dst)
			}
			translateMap[i] = strings.Join(kk, "\n")
			continue
		}
		if smallLen+len(query[i]) > 1000 || i == len(query)-1 {
			smallArray = append(smallArray, query[i])
			mTranslateMap, err := m.translateBase(strings.Join(smallArray, "\n"), from, to)
			if err != nil {
				return nil, err
			}
			for i2, s := range mTranslateMap.TransResult {
				translateMap[i-(len(mTranslateMap.TransResult)-1)+i2] = s.Dst
			}
			smallArray = make([]string, 0)
			smallLen = 0
		} else {
			smallArray = append(smallArray, query[i])
			smallLen = smallLen + len(query[i])
		}
	}
	return translateMap, nil

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
