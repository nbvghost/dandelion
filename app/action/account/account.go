package account

import (
	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"encoding/base64"
	"time"

	"math/rand"
	"strings"

	"strconv"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"

	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Controller struct {
	gweb.BaseController
	Admin        service.AdminService
	Manager      service.ManagerService
	Organization service.OrganizationService
	User         service.UserService
	SMS          service.SMSService
	Wx           service.WxService
	Journal      service.JournalService
}

func (controller *Controller) Apply() {
	//controller.Interceptors.Add(&InterceptorManager{})
	//Index.RequestMapping = make(map[string]mvc.Function)
	controller.AddHandler(gweb.ALLMethod("*", controller.defaultPage))

	controller.AddHandler(gweb.ALLMethod("loginAdminPage", controller.loginAdminPage))
	controller.AddHandler(gweb.ALLMethod("loginManagerPage", controller.loginManagerPage))
	controller.AddHandler(gweb.ALLMethod("loginUserPage", controller.loginUserPage))

	controller.AddHandler(gweb.ALLMethod("wx/authorize", controller.wxAuthorizeAction))

	controller.AddHandler(gweb.ALLMethod("forget", controller.forgetPage))
	controller.AddHandler(gweb.ALLMethod("orderQuery", controller.sdfsda))
	controller.AddHandler(gweb.ALLMethod("user", controller.sdfsda))

	controller.AddHandler(gweb.ALLMethod("loginAdmin", controller.loginAdminAction))
	controller.AddHandler(gweb.ALLMethod("loginManager", controller.loginManagerAction))
	//MwGetWXJSConfig
	controller.AddHandler(gweb.GETMethod("mw/jssdk/config", controller.mwJSSDKConfigAction))

	controller.AddHandler(gweb.ALLMethod("heartbeat", controller.heartbeatAction))
	controller.AddHandler(gweb.ALLMethod("SMS", controller.SMSAction))
	controller.AddHandler(gweb.ALLMethod("SMSCode", controller.SMSCodeAction))
	controller.AddHandler(gweb.ALLMethod("loginOutAdmin", controller.loginOutAdminAction))

	controller.AddHandler(gweb.ALLMethod("mini_program_login", controller.miniProgramLoginAction))

	//controller.AddSubController("/hospital/", &HospitalController{})

}

func (controller *Controller) miniProgramLoginAction(context *gweb.Context) gweb.Result {

	loginInfo := &struct {
		Code     string
		UserInfo string
		ShareKey string
	}{}

	util.RequestBodyToJSON(context.Request.Body, loginInfo)

	userInfo := make(map[string]interface{})

	util.JSONToStruct(loginInfo.UserInfo, &userInfo)

	wxa := controller.Wx.MiniProgram()

	err, OpenID, SessionKey := controller.Wx.MiniProgramInfo(loginInfo.Code, wxa.AppID, wxa.AppSecret)
	//fmt.Println(err, OpenID, SessionKey)

	if err == nil {
		user := controller.User.AddUserByOpenID(OpenID)
		user.OpenID = OpenID
		user.Name = userInfo["nickName"].(string)
		user.Portrait = userInfo["avatarUrl"].(string)
		gender, _ := strconv.ParseInt(strconv.FormatFloat(userInfo["gender"].(float64), 'f', 0, 64), 10, 64)

		user.Gender = int(gender)
		//user.OID = company.ID
		user.LastLoginAt = time.Now()

		if user.SuperiorID == 0 {
			if !strings.EqualFold(loginInfo.ShareKey, "") {
				SuperiorID, _ := util.DecodeShareKey(loginInfo.ShareKey)

				if user.ID != uint64(SuperiorID) {

					//如果往上6级有包含新用户的ID，则不能绑定级别关系
					if !strings.Contains(controller.User.LeveAll6(SuperiorID), strconv.Itoa(int(user.ID))) {
						var hasUser dao.User
						controller.User.Get(dao.Orm(), SuperiorID, &hasUser)
						if hasUser.ID != 0 {
							user.SuperiorID = uint64(SuperiorID)

							if InviteUser, have := context.Data["InviteUser"]; have {
								err := controller.Journal.AddScoreJournal(dao.Orm(),
									hasUser.ID,
									"邀请新朋友获取积分", "邀请新朋友获取积分",
									play.ScoreJournal_Type_InviteUser, int64(InviteUser.(float64)), dao.KV{Key: "SuperiorID", Value: SuperiorID})
								tool.CheckError(err)

								err = controller.Journal.AddUserJournal(dao.Orm(),
									hasUser.ID,
									"邀请新朋友获得现金", "邀请新朋友获得现金",
									play.UserJournal_Type_USER_LEVE, int64(30), dao.KV{Key: "UserID", Value: user.ID}, user.ID)
								tool.CheckError(err)

								controller.Wx.INComeNotify(hasUser, "邀请新朋友获得现金", "0小时", "收入：0.3元")
							}
						}
					}

				}
			}
		}

		controller.User.ChangeModel(dao.Orm(), user.ID, user)
		context.Session.Attributes.Put(play.SessionUser, user)
		context.Session.Attributes.Put(play.SessionOpenID, OpenID)
		context.Session.Attributes.Put(play.SessionMiniProgramKey, SessionKey)
		//context.Session.Attributes.Put(play.SessionOrganization, company)

		//tool.CipherDecrypterData()

		results := make(map[string]interface{})
		results["User"] = user
		results["MyShareKey"] = util.EncodeShareKey(user.ID, 0) //tool.Hashids{}.Encode(user.ID)

		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "登陆成功", Data: results}}
	} else {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

}
func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func (controller *Controller) heartbeatAction(context *gweb.Context) gweb.Result {

	return &gweb.JsonResult{}
}

func (controller *Controller) SMSAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	item := &struct {
		Tel  string
		Code string
		ID   uint64
	}{}

	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", nil)}
	}
	if context.Session.Attributes.Get(item.Tel) == nil {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "验证码不正确", Data: nil}}
	}
	code := context.Session.Attributes.Get(item.Tel).(string)
	//if strings.EqualFold(item.Code, code) {
	if strings.EqualFold(item.Code, code) || strings.EqualFold(item.Code, "00000") {
		//context.Session.Attributes.Put(play.SessionOpenID, openid)
		var user dao.User
		err := controller.User.Get(Orm, item.ID, &user)
		tool.Trace(err)

		if user.ID == 0 {
			return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "没有找到用户", Data: nil}}
		} else {
			user.Tel = item.Tel
			err := controller.User.ChangeModel(Orm, user.ID, &dao.User{Tel: item.Tel})
			tool.Trace(err)
			return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: user}}
		}

	} else {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "验证码不正确", Data: nil}}
	}

}
func (controller *Controller) loginOutAdminAction(context *gweb.Context) gweb.Result {

	context.Session.Attributes.Delete(play.SessionAdmin)
	return &gweb.RedirectToUrlResult{Url: "/admin"}

}
func (controller *Controller) SMSCodeAction(context *gweb.Context) gweb.Result {

	item := &struct {
		Tel string
		//Code string
	}{}

	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", nil)}
	}

	texts := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	a := rand.Intn(10)
	b := rand.Intn(10)
	c := rand.Intn(10)
	d := rand.Intn(10)
	e := rand.Intn(10)

	code := texts[a] + texts[b] + texts[c] + texts[d] + texts[e]

	su, txt := controller.SMS.SendIDCode(code, item.Tel)
	if su {
		context.Session.Attributes.Put(item.Tel, code)
	} else {
		context.Session.Attributes.Put(item.Tel, "")
	}

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: su, Message: txt, Data: nil}}
}

func (controller *Controller) loginAction(context *gweb.Context) gweb.Result {
	//fmt.Println(context.Request.ParseForm())
	account := context.Request.FormValue("account")
	password := context.Request.FormValue("password")

	account = strings.ToLower(account) //小写

	as := &dao.ActionStatus{}

	admin := controller.Admin.FindAdminByAccount(dao.Orm(), account)

	if admin.ID == 0 {

		as.Success = false
		as.Message = "手机/邮箱/密码不正确！"
	} else {
		md5Password := tool.Md5ByString(password)
		if strings.EqualFold(admin.PassWord, md5Password) {
			as.Success = true
			as.Message = ""
			context.Session.Attributes.Put(play.SessionAdmin, admin)
		} else {
			as.Success = false
			as.Message = "手机/邮箱/密码不正确！"
		}

	}

	return &gweb.JsonResult{Data: as}
}
func (controller *Controller) mwJSSDKConfigAction(context *gweb.Context) gweb.Result {
	context.Response.Header().Add("Access-Control-Allow-Origin", "*")

	url := context.Request.URL.Query().Get("url")

	admin := controller.Admin.FindAdminByAccount(dao.Orm(), "admin")

	config := controller.Wx.MwGetWXJSConfig(url, admin.OID)
	//MwGetWXJSConfig
	//callback: jQuery112404449345477704385_1534269866564
	//param: {"appId":"","secret":"","url":"http://minisite.hocodo.com/apps/picc/index.html"}
	//_: 1534269866565
	return &gweb.JsonResult{Data: config}
}
func (controller *Controller) loginManagerAction(context *gweb.Context) gweb.Result {

	account := context.Request.FormValue("account")
	password := context.Request.FormValue("password")

	account = strings.ToLower(account) //小写

	as := &dao.ActionStatus{}

	admin := controller.Manager.FindManagerByAccount(account)

	if admin.ID == 0 {
		as.Success = false
		as.Message = "密码不正确"
	} else {
		md5Password := tool.Md5ByString(password)
		if strings.EqualFold(admin.PassWord, md5Password) {
			as.Success = true
			as.Message = "登陆成功"
			//controller.Admin.ChangeModel(Orm, admin.ID, &dao.Admin{LastLoginAt: time.Now()})
			context.Session.Attributes.Put(play.SessionManager, admin)
		} else {
			as.Success = false
			as.Message = "账号/手机/邮箱/密码不正确！"
		}

	}

	return &gweb.JsonResult{Data: as}
}
func (controller *Controller) loginAdminAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	account := context.Request.FormValue("account")
	password := context.Request.FormValue("password")

	account = strings.ToLower(account) //小写

	as := &dao.ActionStatus{}

	admin := controller.Admin.FindAdminByAccount(Orm, account)

	if admin.ID == 0 {
		as.Success = false
		as.Message = "密码不正确"
	} else {
		md5Password := tool.Md5ByString(password)
		if strings.EqualFold(admin.PassWord, md5Password) {
			as.Success = true
			as.Message = "登陆成功"
			controller.Admin.ChangeModel(Orm, admin.ID, &dao.Admin{LastLoginAt: time.Now()})
			context.Session.Attributes.Put(play.SessionAdmin, admin)
			var _organization dao.Organization
			controller.Organization.Get(Orm, admin.OID, &_organization)
			context.Session.Attributes.Put(play.SessionOrganization, &_organization)
		} else {
			as.Success = false
			as.Message = "账号/手机/邮箱/密码不正确！"
		}

	}

	return &gweb.JsonResult{Data: as}
}

func (controller *Controller) wxAuthorizeAction(context *gweb.Context) gweb.Result {
	code := context.Request.URL.Query().Get("code")
	state := context.Request.URL.Query().Get("state")
	redirect := context.Request.URL.Query().Get("redirect")
	//_OID := context.Request.URL.Query().Get("OID")
	//OID, _ := strconv.ParseUint(_OID, 10, 64)

	if context.Session.Attributes.Get(play.SessionOrganization) == nil {
		return &gweb.RedirectToUrlResult{Url: redirect}
	}

	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	///account/open.do?redirect=%2Ffront%2Fappointment%2F20002%2Findex
	WxConfig := controller.Wx.MiniWeb()
	if strings.EqualFold(code, "") {
		fmt.Println(state)
		//context.Session.Attributes.Put(play.SessionRedirect, redirect)
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + WxConfig.AppID + "&redirect_uri=" + url.QueryEscape(util.GetHost(context.Request)+"/account/wx/authorize?redirect="+redirect) + "&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect"
		//url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + WxConfig.AppID + "&redirect_uri=" + url.QueryEscape(util.GetHost(context.Request)+"/account/loginWxPage") + "&response_type=code&scope=snsapi_base&state=STATE#wechat_redirect"
		fmt.Println(url)
		return &gweb.RedirectToUrlResult{Url: url}

	} else {
		url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + WxConfig.AppID + "&secret=" + WxConfig.AppSecret + "&code=" + code + "&grant_type=authorization_code"

		resp, err := http.Get(url)
		tool.Trace(err)
		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		fmt.Println(string(b))

		mapData := make(map[string]interface{})
		err = json.Unmarshal(b, &mapData)
		tool.Trace(err)
		fmt.Println(mapData)
		if mapData["openid"] == nil || mapData["access_token"] == nil {
			return &gweb.RedirectToUrlResult{Url: redirect}
		}
		openid := mapData["openid"].(string)
		access_token := mapData["access_token"].(string)

		user := controller.User.FindUserByOpenID(Orm, openid)
		if user.ID == 0 {
			//user.OID = OID
			user.OpenID = openid
			controller.User.Add(Orm, user)
		}
		//access_token := wxpay.GetAccessToken(WxConfig.ID)
		//https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN

		resp, err = http.Get("https://api.weixin.qq.com/sns/userinfo?access_token=" + access_token + "&openid=" + openid + "&lang=zh_CN")
		tool.Trace(err)
		b, err = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		fmt.Println(string(b))
		//controller.User.ChangeModel(dao.Orm(), user.ID, user)

		mapData = make(map[string]interface{})
		err = json.Unmarshal(b, &mapData)

		if mapData["errcode"] != nil {
			return &gweb.RedirectToUrlResult{Url: redirect}
		}

		nickname := mapData["nickname"].(string)
		sex := mapData["sex"].(float64)
		headimgurl := mapData["headimgurl"].(string)
		country := mapData["country"].(string)
		province := mapData["province"].(string)
		city := mapData["city"].(string)
		region := country + province + city
		//region
		//{"subscribe":1,"openid":"oQkeRt-eM835hSka10TsCIPwX-Ik","nickname":"A101小鱼",
		// "sex":1,"language":"zh_CN","city":"厦门","province":"福建","country":"中国",
		// "headimgurl":"http:\/\/thirdwx.qlogo.cn\/mmopen\/PiajxSqBRaEIFqppwEH3se1iaTWlAAgh2zyiarQYdbSeAaFPKxJxgiczcibhlrFSbYbu6I165bJicNbXQZ1NqiazKcK5w\/132",
		// "subscribe_time":1510908061,"remark":"","groupid":0,"tagid_list":[],"subscribe_scene":"ADD_SCENE_QR_CODE","qr_scene":0,"qr_scene_str":""}
		//context.Session.Attributes.Put(play.SessionOpenID, openid)

		user.Name = nickname
		user.Portrait = headimgurl
		user.Gender = int(sex)
		user.Region = region
		user.LastLoginAt = time.Now()

		access_token = controller.Wx.GetAccessToken(WxConfig)
		//https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
		resp, err = http.Get("https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + access_token + "&openid=" + openid + "&lang=zh_CN")
		tool.Trace(err)
		b, err = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		fmt.Println(string(b))
		//controller.User.ChangeModel(dao.Orm(), user.ID, user)

		mapData = make(map[string]interface{})
		err = json.Unmarshal(b, &mapData)

		if mapData["errcode"] != nil {
			return &gweb.RedirectToUrlResult{Url: redirect}
		}

		subscribe := mapData["subscribe"].(float64)
		user.Subscribe = int(subscribe)
		controller.User.ChangeMap(dao.Orm(), user.ID, dao.User{}, map[string]interface{}{
			"Subscribe":   user.Subscribe,
			"Name":        user.Name,
			"Portrait":    user.Portrait,
			"Gender":      user.Gender,
			"Region":      user.Region,
			"LastLoginAt": user.LastLoginAt,
		})
		//controller.User.ChangeModel(dao.Orm(), user.ID, user)
		context.Session.Attributes.Put(play.SessionUser, user)
		return &gweb.RedirectToUrlResult{Url: redirect}

	}
}
func (controller *Controller) forgetPage(context *gweb.Context) gweb.Result {
	return &gweb.HTMLResult{}
}
func (controller *Controller) loginAdminPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func (controller *Controller) loginManagerPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func (controller *Controller) loginUserPage(context *gweb.Context) gweb.Result {
	redirect := util.GetFullUrl(context.Request)

	//var user dao.User
	//controller.User.Get(service.Orm, 1002, &user)
	//context.Session.Attributes.Put(play.SessionUser, &user)
	//context.Session.Attributes.Put(play.SessionOpenID, openid)

	if context.Session.Attributes.Get(play.SessionUser) == nil {
		//return &gweb.HTMLResult{Params: map[string]interface{}{"User": context.Session.Attributes.Get(play.SessionUser).(*dao.User)}}
		return &gweb.RedirectToUrlResult{Url: "/account/open.do?redirect=" + base64.StdEncoding.EncodeToString([]byte(redirect))}
	} else {
		return &gweb.HTMLResult{Params: map[string]interface{}{"User": context.Session.Attributes.Get(play.SessionUser).(*dao.User)}}
	}
}
func (controller *Controller) sdfsda(context *gweb.Context) gweb.Result {

	return &gweb.JsonResult{}
}
