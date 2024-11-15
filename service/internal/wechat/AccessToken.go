package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/entity/model"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var accessTokenMap = make(map[string]*AccessToken)

type AccessTokenService struct {
}

func (m AccessTokenService) GetAccessToken(wechatConfig *model.WechatConfig) string {

	if accessTokenMap[wechatConfig.AppID] != nil && (time.Now().Unix()-accessTokenMap[wechatConfig.AppID].Update) < accessTokenMap[wechatConfig.AppID].Expires_in {

		return accessTokenMap[wechatConfig.AppID].Access_token
	}

	//WxConfig := model.GetWxConfig(WxConfigID)

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + wechatConfig.AppID + "&secret=" + wechatConfig.AppSecret

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	d := make(map[string]interface{})

	err = json.Unmarshal(b, &d)
	if err != nil {
		log.Println(err)
		return ""
	}
	//fmt.Println(string(b))
	//fmt.Println(d)
	if d["access_token"] == nil {
		return ""
	}
	at := &AccessToken{}
	at.Access_token = d["access_token"].(string)
	at.Expires_in = int64(d["expires_in"].(float64))
	at.Update = time.Now().Unix()
	accessTokenMap[wechatConfig.AppID] = at
	return accessTokenMap[wechatConfig.AppID].Access_token
}

func (m AccessTokenService) MiniProgramInfo(Code, AppID, AppSecret string) (err error, OpenID, SessionKey string) {

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
