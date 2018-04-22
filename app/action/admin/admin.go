package admin

import (
	"net/http"
	"net/url"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"

	"encoding/json"

	"strconv"

	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/service/dao"
)

type InterceptorAdmin struct {
}

func (this InterceptorAdmin) Execute(context *gweb.Context) bool {

	//util.Trace(context.Session,"context.Session")
	if context.Session.Attributes.Get(play.SessionAdmin) == nil {
		//http.SetCookie(context.Response, &http.Cookie{Name: "UID", MaxAge:-1, Path: "/"})
		//fmt.Println(context.Request.URL.Path)
		//fmt.Println(context.Request.URL.Query().Encode())
		redirect := ""
		if len(context.Request.URL.Query().Encode()) == 0 {
			redirect = context.Request.URL.Path
		} else {
			redirect = context.Request.URL.Path + "?" + context.Request.URL.Query().Encode()
		}
		//fmt.Println(url.QueryEscape(redirect))
		http.Redirect(context.Response, context.Request, "/account/loginAdminPage?redirect="+url.QueryEscape(redirect), http.StatusFound)
		return false
	} else {
		return true
	}
}

type InterceptorData struct {
}

func (this InterceptorData) Execute(context *gweb.Context) bool {

	return true
}

type Controller struct {
	gweb.BaseController
}

func (i *Controller) Apply() {
	i.Interceptors.Add(&InterceptorAdmin{})
	i.Interceptors.Add(&InterceptorData{})
	i.AddHandler(gweb.ALLMethod("", defaultPage))
	i.AddHandler(gweb.ALLMethod("*", mainPage))
	i.AddHandler(gweb.ALLMethod("index", mainPage))
	i.AddHandler(gweb.ALLMethod("shop", shopAction))
	i.AddHandler(gweb.ALLMethod("appointmentAction", appointmentAction))
	i.AddHandler(gweb.ALLMethod("articlePage", articlePage))
	i.AddHandler(gweb.ALLMethod("add_article", addArticlePage))
	i.AddHandler(gweb.ALLMethod("classifyAction", classifyAction))
	//i.AddHandler("wxauthorizer",  wxauthorizerAction})

}

/*func wxauthorizerAction(context *gweb.Context) gweb.Result {

	config := service.Configuration.GetConfiguration(service.Orm, play.ConfigurationKey_component_verify_ticket)

	auth_code := context.Request.URL.Query().Get("auth_code")
	//expires_in := context.Request.URL.Query().Get("expires_in")

	if strings.EqualFold(auth_code, "") == false {
		authorizer_appid, authorizer_access_token, authorizer_refresh_token, func_info, expires_in := wxpay.Api_query_auth(auth_code, config.V)

		admin := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)

		item := service.WxConfig.GetWxConfig(service.Orm, admin.CompanyID)
		err := service.WxConfig.ChangeWxConfig(service.Orm, item.ID, dao.WxConfig{
			AuthorizerAppID:        authorizer_appid,
			AuthorizerAccessToken:  authorizer_access_token,
			AuthorizerExpiresIn:    expires_in,
			AuthorizerRefreshToken: authorizer_refresh_token,
			AuthorizerFuncInfo:     func_info,
		})
		tool.CheckError(err)
	}

	Component_access_token := wxpay.Api_component_token(config.V)

	PreAuthCode := wxpay.Api_create_preauthcode(Component_access_token)

	return &gweb.HTMLResult{Params: map[string]interface{}{"AppID": wxpay.OpenAppID, "PreAuthCode": PreAuthCode}}
}*/
func appointmentAction(context *gweb.Context) gweb.Result {
	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	case play.ActionKey_get:
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
		d := service.Appointment.GetAppointment(ID)
		result.Data = &dao.ActionStatus{false, "", d}
	case play.ActionKey_list:
		shop := context.Session.Attributes.Get(play.SessionShop).(*dao.Company)
		Index, _ := strconv.Atoi(context.Request.URL.Query().Get("Index"))
		List, Total := service.Appointment.FindAppointmentOfPaging(Index, shop.ID)
		result.Data = &dao.Pager{List, Index, Total, play.Paging}
	case play.ActionKey_save:
		context.Request.ParseForm()
		shop := context.Session.Attributes.Get(play.SessionShop).(*dao.Company)
		jsonText := context.Request.Form.Get("json")
		appointment := &dao.Appointment{}
		err := json.Unmarshal([]byte(jsonText), appointment)
		appointment.CompanyID = shop.ID
		tool.CheckError(err)
		if err == nil {
			result.Data = service.Appointment.SaveAppointment(appointment)
		} else {
			result.Data = &dao.ActionStatus{false, err.Error(), nil}
		}
	}
	return result
}
func shopAction(context *gweb.Context) gweb.Result {
	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	case play.ActionKey_get:
		shop := context.Session.Attributes.Get(play.SessionShop).(*dao.Company)
		result.Data = &dao.ActionStatus{true, "ok", shop}
	case play.ActionKey_change:
		context.Request.ParseForm()

		Photos := context.Request.Form.Get("Photos")
		Categories := context.Request.Form.Get("Categories")
		Province := context.Request.Form.Get("Province")
		City := context.Request.Form.Get("City")
		District := context.Request.Form.Get("District")
		Name := context.Request.Form.Get("Name")
		Address := context.Request.Form.Get("Address")
		Telephone := context.Request.Form.Get("Telephone")
		Special := context.Request.Form.Get("Special")
		Opentime := context.Request.Form.Get("Opentime")
		Avgprice := context.Request.Form.Get("Avgprice")
		Introduction := context.Request.Form.Get("Introduction")
		Recommend := context.Request.Form.Get("Recommend")

		_shop := context.Session.Attributes.Get(play.SessionShop).(*dao.Company)

		su := service.Company.ChangeCompany(_shop.ID, Photos, Categories, Province, City, District, Name, Address, Telephone, Special, Opentime, Avgprice, Introduction, Recommend)

		if su {
			_shop := context.Session.Attributes.Get(play.SessionShop).(*dao.Company)
			_admin := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)

			shop := service.Company.GetCompany(_shop.ID)
			context.Session.Attributes.Put(play.SessionShop, shop)

			admin := service.Admin.GetAdmin(service.Orm, _admin.ID)
			context.Session.Attributes.Put(play.SessionAdmin, admin)
		}

		result.Data = (&dao.ActionStatus{}).Smart(su, "修改成功", "修改失败")

	}
	return result
}
func classifyAction(context *gweb.Context) gweb.Result {
	result := &gweb.JsonResult{}
	action := context.Request.URL.Query().Get("action")
	switch action {
	case play.ActionKey_list:
		admin := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)
		data := service.Classify.FindByShopID(service.Orm, admin.CompanyID)
		result.Data = &dao.ActionStatus{true, "SUCCESS", data}
	case play.ActionKey_change:
		context.Request.ParseForm()
		ID, _ := strconv.ParseUint(context.Request.Form.Get("ID"), 10, 64)
		Label := context.Request.Form.Get("Label")
		err := service.Classify.ChangeModel(service.Orm, ID, &dao.Classify{Label: Label})
		result.Data = (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)
	case play.ActionKey_del:
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
		err := service.Classify.Delete(service.Orm, dao.Classify{}, ID)
		result.Data = (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)
	case play.ActionKey_add:
		context.Request.ParseForm()
		admin := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)
		Label := context.Request.Form.Get("Label")

		classify := &dao.Classify{}
		classify.CompanyID = admin.CompanyID
		classify.Label = Label

		err := service.Classify.AddClassifyNotNull(classify)

		result.Data = (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)

	}

	return result
}
func addArticlePage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{Params: nil}
}
func articlePage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func defaultPage(context *gweb.Context) gweb.Result {

	return &gweb.RedirectToUrlResult{"/admin/index"}
}
func mainPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
