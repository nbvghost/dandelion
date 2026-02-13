package account

import (
	"bytes"
	"encoding/json"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"io/ioutil"
	"net/http"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/pkg/errors"
)

type GetLoginUserPhone struct {
	User         *model.User         `mapping:""`
	WechatConfig *model.WechatConfig `mapping:""`
	Post         struct {
		iv            string
		encryptedData string
		Code          string
	} `method:"Post"`
}

func (g *GetLoginUserPhone) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
func (g *GetLoginUserPhone) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	accessToken := service.Wechat.AccessToken.GetAccessToken(g.WechatConfig)

	body, err := json.Marshal(map[string]any{"code": g.Post.Code})
	if err != nil {
		return nil, err
	}
	post, err := http.Post("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token="+accessToken, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer post.Body.Close()
	body, err = ioutil.ReadAll(post.Body)
	if err != nil {
		return nil, err
	}

	var rb ResultBody
	err = json.Unmarshal(body, &rb)
	if err != nil {
		return nil, err
	}
	if rb.Errcode != 0 {
		return nil, errors.New(rb.Errmsg)
	}

	//CountryCode: object.ParseInt(rb.PhoneInfo.CountryCode)
	err = dao.UpdateByPrimaryKey(db.GetDB(ctx), &model.User{}, ctx.UID(), &model.User{Phone: rb.PhoneInfo.PurePhoneNumber})
	if err != nil {
		return nil, err
	}
	r := result.NewData(map[string]any{"User": dao.GetByPrimaryKey(db.GetDB(ctx), &model.User{}, ctx.UID())})
	r.Message = "成功绑定手机"
	return r, nil
}

type ResultBody struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int    `json:"timestamp"`
			Appid     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}
