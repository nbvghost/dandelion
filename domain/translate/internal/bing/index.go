package bing

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/tool/object"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Index struct {

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

func (m *Index) bingTranslateParams() error {
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

func (m *Index) BingTranslate(query, from, to string) (_ string, err error) {
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
		return m.BingTranslate(query, from, to)
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