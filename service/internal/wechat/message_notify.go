package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/company"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
)

type MessageNotify struct {
	WxService           WxService
	OrganizationService company.OrganizationService
	AccessTokenService  AccessTokenService
}

// 新用户加入，绑定上下级关系
func (m MessageNotify) NewUserJoinNotify(NewUser *model.User, notifyUser *model.User) *result.ActionResult {

	as := &result.ActionResult{}

	/*userFormID := service.User.GetFromIDs(notifyUser.ID)
	if userFormID.ID == 0 {
		as.Code = result.Fail
		as.Message = "没有找到，用户的formid"
	} else {

		sendData := make(map[string]interface{})
		sendData["touser"] = notifyUser.OpenID

		weapp_template_msg_data := make(map[string]interface{})
		weapp_template_msg_data["page"] = "pages/user/user"
		weapp_template_msg_data["template_id"] = "YfEY2Xbju5-fm3Naww3EbVYQPUPIjorESo-KV-KXZvs"
		weapp_template_msg_data["form_id"] = userFormID.FormId

		data_data := make(map[string]interface{})
		data_data["keyword1"] = map[string]interface{}{"value": strconv.Itoa(int(NewUser.ID)), "color": "#173177"}
		if NewUser.Gender == 1 {
			data_data["keyword2"] = map[string]interface{}{"value": "男", "color": "#173177"}
		} else if NewUser.Gender == 2 {
			data_data["keyword2"] = map[string]interface{}{"value": "女", "color": "#173177"}
		} else {
			data_data["keyword2"] = map[string]interface{}{"value": "未知", "color": "#173177"}
		}

		data_data["keyword3"] = map[string]interface{}{"value": NewUser.CreatedAt.Format("2006-01-02 15:04:05"), "color": "#173177"}
		data_data["keyword4"] = map[string]interface{}{"value": NewUser.Name, "color": "#173177"}
		data_data["keyword5"] = map[string]interface{}{"value": NewUser.Name + "已经成为您的好友，他（她）下单您会获得奖励喔！", "color": "#173177"}

		weapp_template_msg_data["data"] = data_data

		sendData["weapp_template_msg"] = weapp_template_msg_data

		var errcode int
		as, errcode = service.SendUniformMessage(sendData)
		if as.Code == result.Success || errcode == 41028 {
			service.User.Delete(singleton.Orm(), &model.UserFormIds{}, userFormID.ID)
		}

	}*/

	return as
}

// 发货通知
func (m MessageNotify) OrderDeliveryNotify(Order *model.Orders, ogs []dao.IEntity, wxConfig *model.WechatConfig) *result.ActionResult {

	if Order.ID == 0 {
		return &result.ActionResult{Code: result.Fail, Message: "找不到订单", Data: nil}
	}

	notifyUser := dao.GetByPrimaryKey(db.Orm(), &model.User{}, Order.UserID).(*model.User)

	var as *result.ActionResult

	weapp_template_msg_data := make(map[string]interface{})
	weapp_template_msg_data["page"] = "pages/user/user"
	weapp_template_msg_data["template_id"] = "MHiJR_3T2W4LJVhwOVctO6Lr7fxC9rSCO924dwSoYrY"
	weapp_template_msg_data["form_id"] = Order.PrepayID
	weapp_template_msg_data["touser"] = notifyUser.OpenID

	data_data := make(map[string]interface{})
	data_data["keyword1"] = map[string]interface{}{"value": Order.ShipInfo.Name, "color": "#173177"}
	data_data["keyword2"] = map[string]interface{}{"value": Order.ShipInfo.No, "color": "#173177"}

	var Titles = ""
	for i := range ogs {
		value := ogs[i].(*model.OrdersGoods)
		//var goods model.Goods
		//json.Unmarshal([]byte(value.Goods), &goods)
		Titles += value.Goods.Title
	}
	if len(Titles) > 48 {
		Titles = Titles[:48] + "等"
	}

	data_data["keyword3"] = map[string]interface{}{"value": Titles, "color": "#173177"}
	data_data["keyword4"] = map[string]interface{}{"value": Order.DeliverTime.Format("2006-01-02 15:04:05"), "color": "#173177"}

	weapp_template_msg_data["data"] = data_data

	as = m.SendWXMessage(weapp_template_msg_data, wxConfig)

	return as
}

//收入提醒
/*
@slUser 收入的用户
*/
func (m MessageNotify) INComeNotify(slUser *model.User, itemName string, timeText string, typeText string) *result.ActionResult {
	//
	var as = &result.ActionResult{Code: result.Fail}

	/*if slUser.ID == 0 {
		return &result.ActionResult{Code: result.Fail, Message: "用户不存在", Data: nil}
	}

	//var notifyUser model.User
	//model.User.Get(singleton.Orm(), slUser.SuperiorID, &notifyUser)



	userFormID := service.User.GetFromIDs(slUser.ID)
	if userFormID.ID == 0 {
		as.Code = result.Fail
		as.Message = "没有找到，用户的formid"

	} else {

		sendData := make(map[string]interface{})
		sendData["touser"] = slUser.OpenID

		weapp_template_msg_data := make(map[string]interface{})
		weapp_template_msg_data["page"] = "pages/user/user"
		weapp_template_msg_data["template_id"] = "xV23xWZgdNViUiD1fk-1edKNY7QNJnv4SD6tY7pu8w4"
		weapp_template_msg_data["form_id"] = userFormID.FormId

		data_data := make(map[string]interface{})
		data_data["keyword1"] = map[string]interface{}{"value": itemName, "color": "#173177"}
		data_data["keyword2"] = map[string]interface{}{"value": timeText, "color": "#173177"}
		data_data["keyword3"] = map[string]interface{}{"value": typeText, "color": "#ff0000"}

		weapp_template_msg_data["data"] = data_data

		sendData["weapp_template_msg"] = weapp_template_msg_data

		var errcode int
		as, errcode = service.SendUniformMessage(sendData)
		if as.Code == result.Success || errcode == 41028 {
			service.User.Delete(singleton.Orm(), &model.UserFormIds{}, userFormID.ID)
		}

	}*/

	return as
}

// 新订单
func (m MessageNotify) NewOrderNotify(Order model.Orders, ogs []model.OrdersGoods, wxConfig *model.WechatConfig) *result.ActionResult {

	if Order.ID == 0 {
		return &result.ActionResult{Code: result.Fail, Message: "找不到订单", Data: nil}
	}

	notifyUser := dao.GetByPrimaryKey(db.Orm(), &model.User{}, Order.UserID).(*model.User)

	var as *result.ActionResult

	weapp_template_msg_data := make(map[string]interface{})
	weapp_template_msg_data["page"] = "pages/user/user"
	weapp_template_msg_data["template_id"] = "bah5ch6kSTi4dvbYzlZ80m7usPIe7PWZEW7Csk_HOy0"
	weapp_template_msg_data["form_id"] = Order.PrepayID
	weapp_template_msg_data["touser"] = notifyUser.OpenID

	data_data := make(map[string]interface{})
	data_data["keyword1"] = map[string]interface{}{"value": notifyUser.Name, "color": "#173177"}
	data_data["keyword2"] = map[string]interface{}{"value": Order.OrderNo, "color": "#173177"}

	var address model.Address
	json.Unmarshal([]byte(Order.Address), &address)
	addressText := address.Name + "/" + address.ProvinceName + address.CityName + address.CountyName + address.Detail + address.PostalCode + "/" + address.Tel

	data_data["keyword3"] = map[string]interface{}{"value": addressText, "color": "#173177"}

	data_data["keyword4"] = map[string]interface{}{"value": Order.PayTime.Format("2006-01-02 15:04:05"), "color": "#173177"}

	//var org model.Organization
	//service.OrganizationService.Get(singleton.Orm(), Order.OID, &org)
	org := dao.GetByPrimaryKey(db.Orm(), &model.Organization{}, Order.OID).(*model.Organization)

	data_data["keyword5"] = map[string]interface{}{"value": org.Name, "color": "#173177"}

	data_data["keyword6"] = map[string]interface{}{"value": strconv.Itoa(int(Order.PayMoney/100)) + "元", "color": "#173177"}

	var Titles = ""
	for _, value := range ogs {
		//var goods model.Goods
		//json.Unmarshal([]byte(value.Goods), &goods)
		Titles += value.Goods.Title
	}
	if len(Titles) > 48 {
		Titles = Titles[:48] + "等"
	}
	data_data["keyword7"] = map[string]interface{}{"value": Titles, "color": "#173177"}
	data_data["keyword8"] = map[string]interface{}{"value": "如有疑问，请联系客服", "color": "#173177"}

	weapp_template_msg_data["data"] = data_data

	as = m.SendWXMessage(weapp_template_msg_data, wxConfig)

	return as
}
func (m MessageNotify) SendUniformMessage(sendData map[string]interface{}, wxConfig *model.WechatConfig) (*result.ActionResult, int) {

	//gzh := model.MiniWeb()
	//xcx := service.MiniProgram()

	b, err := json.Marshal(sendData)
	log.Println(err)

	access_token := m.AccessTokenService.GetAccessToken(wxConfig)
	strReader := strings.NewReader(string(b))
	respones, err := http.Post("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/uniform_send?access_token="+access_token, "application/json", strReader)
	log.Println(err)
	if err != nil {
		return &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}, -1
	}
	defer respones.Body.Close()
	body, err := ioutil.ReadAll(respones.Body)
	log.Println(err)
	mapData := make(map[string]interface{})
	fmt.Println(string(body))
	err = json.Unmarshal(body, &mapData)
	log.Println(err)
	if mapData["errcode"] != nil {
		if mapData["errcode"].(float64) == 0 {
			return &result.ActionResult{Code: result.Success, Message: "发送成功", Data: nil}, 0
		}
	}
	return &result.ActionResult{Code: result.Fail, Message: mapData["errmsg"].(string), Data: nil}, int(mapData["errcode"].(float64))

}
func (m MessageNotify) SendWXMessage(sendData map[string]interface{}, wxConfig *model.WechatConfig) *result.ActionResult {
	b, err := json.Marshal(sendData)
	log.Println(err)

	//WxConfig := service.MiniProgram()

	access_token := m.AccessTokenService.GetAccessToken(wxConfig)
	strReader := strings.NewReader(string(b))
	respones, err := http.Post("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token="+access_token, "application/json", strReader)
	log.Println(err)
	if err != nil {
		return &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}
	}
	defer respones.Body.Close()
	body, err := ioutil.ReadAll(respones.Body)
	log.Println(err)
	mapData := make(map[string]interface{})
	fmt.Println(string(body))
	err = json.Unmarshal(body, &mapData)
	log.Println(err)
	if mapData["errcode"] != nil {
		if mapData["errcode"].(float64) == 0 {
			return &result.ActionResult{Code: result.Success, Message: "发送成功", Data: nil}
		} else {
			//mapData["errcode"].(float64)
			return &result.ActionResult{Code: result.Fail, Message: mapData["errmsg"].(string), Data: nil}
		}
	}
	return &result.ActionResult{Code: result.Fail, Message: mapData["errmsg"].(string), Data: nil}

}
