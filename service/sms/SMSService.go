package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/aliyun"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/library/result"

	"github.com/nbvghost/tool/encryption"

	"github.com/nbvghost/tool"
	"github.com/nbvghost/tool/collections"
)

type Service struct {
}

func (s Service) PhoneVerifyCode(ctx constrain.IContext, phone string, code string) (bool, error) {
	redisKey := key.NewRedisVerifyPhoneCodeKey(phone)
	has, err := ctx.Redis().SetIsMember(ctx, redisKey, code)
	if err != nil {
		return false, err
	}
	if !has {
		return false, errors.Errorf("验证码不正确")
	}
	_, err = ctx.Redis().SetRem(ctx, redisKey, code)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s Service) SendPhoneVerifyCode(ctx constrain.IContext, oid dao.PrimaryKey, phone string) error {
	contextValue := contexext.FromContext(ctx)

	desi := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	codes := desi[key.Random.Intn(10)] + desi[key.Random.Intn(10)] + desi[key.Random.Intn(10)] + desi[key.Random.Intn(10)] + desi[key.Random.Intn(10)]

	ip := util.GetIP(contextValue.Request)

	redisKey := key.NewRedisVerifyPhoneCodeKey(phone)

	memberLen, err := ctx.Redis().SetCard(ctx, redisKey)
	if err != nil {
		return err
	}
	if memberLen >= 5 {
		return fmt.Errorf("验证码请求太频繁，请稍后在试")
	}
	_, err = ctx.Redis().SetAdd(ctx, redisKey, codes)
	if err != nil {
		return err
	}
	err = ctx.Redis().Expire(ctx, redisKey, time.Minute*5)
	if err != nil {
		return err
	}
	sms := aliyun.NewSMS(oid)
	err = sms.Send("SMS_4785426", phone, map[string]any{"code": codes, "product": "手机号"})
	if err != nil {
		return err
	}
	_, err = ctx.Redis().Incr(ctx, key.NewRedisVerifyPhoneIPKey(ip))
	if err != nil {
		return err
	}
	return nil
}
func specialUrlEncode(value string) string {
	sdt := url.QueryEscape(value)
	sdt = strings.Replace(sdt, "+", "%20", -1)
	sdt = strings.Replace(sdt, "*", "%2A", -1)
	sdt = strings.Replace(sdt, "%7E", "~", -1)
	return sdt
}

func (s Service) SendAliyunSms(ParamMap map[string]interface{}, TemplateCode string, tel string, accessKeyId, accessSecret string) (bool, string) {

	//accessKeyId:=""
	//accessSecret:=""

	params := url.Values{}
	params.Add("SignatureMethod", "HMAC-SHA1")
	params.Add("SignatureNonce", tool.UUID())
	params.Add("AccessKeyId", accessKeyId)
	params.Add("SignatureVersion", "1.0")

	te := time.Now().UTC()
	fmt.Println(te.Format("2006-01-02 15:04:05"))
	params.Add("Timestamp", te.Format("2006-01-02 15:04:05"))
	params.Add("Format", "json")

	// 2. 业务API参数
	params.Add("Action", "SendSms")
	params.Add("Version", "2017-05-25")
	params.Add("RegionId", "cn-hangzhou")
	params.Add("PhoneNumbers", tel)
	params.Add("SignName", "美蒂欧官网")

	ParamMapBtytes, _ := json.Marshal(ParamMap)

	params.Add("TemplateParam", string(ParamMapBtytes))
	params.Add("TemplateCode", TemplateCode)
	//params.Add("OutId", "123")

	list := &collections.ListString{}
	for k, v := range params {
		list.Append(specialUrlEncode(k) + "=" + specialUrlEncode(v[0]))
	}
	list.SortL()

	queryString := list.Join("&")

	fmt.Println(queryString)

	stringToSign := "GET&" + specialUrlEncode("/") + "&" + specialUrlEncode(queryString)

	key := []byte(accessSecret + "&")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(stringToSign))
	//fmt.Printf("%x\n", mac.Sum(nil))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	signature := specialUrlEncode(sign)
	//fmt.Printf("%x\n", sign)

	//sign := tool.Md5ByString(accessSecret +"&"+ stringToSign)
	//fmt.Println(Secret + list.Join("") + Secret)

	//params.Add("Signature", sign)
	//queryString

	fmt.Println(signature, url.QueryEscape(queryString))

	resp, err := http.Get("http://dysmsapi.aliyuncs.com?Signature=" + signature + "&" + queryString)
	log.Println(err)
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(b))

	//map[error_response:map[code:15 msg:Remote service error sub_code:isv.BUSINESS_LIMIT_CONTROL sub_msg:触发分钟级流控Permits:1 request_id:z2b5hq31ty7v]]
	//map[alibaba_aliqin_fc_sms_num_send_response:map[result:map[err_code:0 model:465323720574031803^0 msg:OK success:true] request_id:zioman9yylzz]]
	result := make(map[string]interface{})
	err = json.Unmarshal(b, &result)
	log.Println(err)
	//fmt.Println(result)
	if result["error_response"] != nil {
		var dd = result["error_response"].(map[string]interface{})["sub_msg"].(string)
		return false, dd
	} else {
		if result["alibaba_aliqin_fc_sms_num_send_response"] != nil {
			_resu := result["alibaba_aliqin_fc_sms_num_send_response"].(map[string]interface{})["result"].(map[string]interface{})["success"].(bool)
			if _resu {
				return true, ""
			} else {
				return false, "短信发送过快"
			}
		} else {
			return true, ""
		}

	}

}
func (s Service) SendIDCode(Code string, tel string) (result.ActionResultCode, string) {

	params := url.Values{}
	params.Add("app_key", "24807838")
	params.Add("format", "json")
	params.Add("method", "alibaba.aliqin.fc.sms.num.send")
	//params.Add("partner_id", "apidoc")

	params.Add("sign_method", "md5")
	params.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Add("v", "2.0")
	//params.Add("extend", "123456")
	params.Add("rec_num", tel)
	params.Add("sms_free_sign_name", "登录验证")
	params.Add("sms_param", `{"code":"`+Code+`","product":"网上预约服务"}`)
	params.Add("sms_template_code", "SMS_5068557")
	params.Add("sms_type", "normal")

	list := &collections.ListString{}
	for k, v := range params {
		list.Append(k + v[0])
	}
	list.SortL()

	Secret := "1759a0c33fea083eed5c1a5df4cc496e"

	sign := encryption.Md5ByString(Secret + list.Join("") + Secret)
	//fmt.Println(Secret + list.Join("") + Secret)

	params.Add("sign", sign)

	resp, err := http.PostForm("http://gw.api.taobao.com/router/rest", params)
	log.Println(err)
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	//map[error_response:map[code:15 msg:Remote service error sub_code:isv.BUSINESS_LIMIT_CONTROL sub_msg:触发分钟级流控Permits:1 request_id:z2b5hq31ty7v]]
	//map[alibaba_aliqin_fc_sms_num_send_response:map[result:map[err_code:0 model:465323720574031803^0 msg:OK success:true] request_id:zioman9yylzz]]
	Result := make(map[string]interface{})
	err = json.Unmarshal(b, &Result)
	log.Println(err)
	//fmt.Println(result)
	if Result["error_response"] != nil {
		var dd = Result["error_response"].(map[string]interface{})["sub_msg"].(string)
		return result.Fail, dd
	} else {
		if Result["alibaba_aliqin_fc_sms_num_send_response"] != nil {
			_resu := Result["alibaba_aliqin_fc_sms_num_send_response"].(map[string]interface{})["result"].(map[string]interface{})["success"].(bool)
			if _resu {
				return result.Success, ""
			} else {
				return result.Fail, "短信发送过快"
			}
		} else {
			return result.Success, ""
		}

	}
}
