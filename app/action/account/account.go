package account

import (
	"fmt"
	"net/url"
	"strings"

	"strconv"

	"net/http"

	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/service/dao"

	"dandelion/app/util"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type Controller struct {
	gweb.BaseController
}

func (i *Controller) Apply() {

	//Index.RequestMapping = make(map[string]mvc.Function)
	i.AddHandler(gweb.ALLMethod("loginAdminPage", loginAdminPage))
	i.AddHandler(gweb.ALLMethod("loginManagerPage", loginManagerPage))
	i.AddHandler(gweb.ALLMethod("loginUserPage", loginUserPage))
	i.AddHandler(gweb.ALLMethod("forget", forgetPage))
	i.AddHandler(gweb.ALLMethod("open.do", openDo))
	i.AddHandler(gweb.ALLMethod("orderQuery", sdfsda))
	i.AddHandler(gweb.ALLMethod("user", sdfsda))
	i.AddHandler(gweb.ALLMethod("login", loginAction))
	i.AddHandler(gweb.ALLMethod("loginManager", loginManager))
	i.AddHandler(gweb.ALLMethod("loginUser", loginUserAction))
	i.AddHandler(gweb.ALLMethod(":shopID/user", sdfsda))
	i.AddHandler(gweb.ALLMethod("loginOut", sdfsda))
	i.AddHandler(gweb.ALLMethod("userLogin/:action", sdfsda))
	i.AddHandler(gweb.ALLMethod(":shopID/register", sdfsda))
	i.AddHandler(gweb.ALLMethod("payNotify", sdfsda))
	i.AddHandler(gweb.ALLMethod("popularize_info/:shopID", sdfsda))
	i.AddHandler(gweb.ALLMethod("popularize/:shopID", sdfsda))
	i.AddHandler(gweb.ALLMethod("pay/platform_pay", sdfsda))
	i.AddHandler(gweb.ALLMethod("pay/platform_order_create", sdfsda))
	i.AddHandler(gweb.ALLMethod(":action/:shopID/expire", sdfsda))
	i.AddHandler(gweb.ALLMethod("transfers/:orderID", sdfsda))
	i.AddHandler(gweb.ALLMethod("transfers/:orderID/get", sdfsda))
	i.AddHandler(gweb.ALLMethod("article/get/:id", articleGet))
	i.AddHandler(gweb.ALLMethod("heartbeat", heartbeatAction))
	i.AddHandler(gweb.ALLMethod("*", defaultPage))

}
func defaultPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func heartbeatAction(context *gweb.Context) gweb.Result {

	return &gweb.JsonResult{}
}
func articleGet(context *gweb.Context) gweb.Result {
	id, _ := strconv.ParseUint(context.PathParams["id"], 10, 64)

	article := service.Article.GetArticle(id)

	return &gweb.JsonResult{Data: article}
}

func loginAction(context *gweb.Context) gweb.Result {
	//fmt.Println(context.Request.ParseForm())
	account := context.Request.FormValue("account")
	password := context.Request.FormValue("password")

	account = strings.ToLower(account) //小写

	as := &dao.ActionStatus{}

	admin := service.Admin.GetAdminByEmail(service.Orm, account)
	if admin.ID == 0 {
		admin = service.Admin.GetAdminByTel(service.Orm, account)
	}
	if admin.ID == 0 {

		as.Success = false
		as.Message = "手机/邮箱/密码不正确！"
	} else {
		md5Password := tool.Md5(password)
		if strings.EqualFold(admin.Password, md5Password) {
			as.Success = true
			as.Message = ""
			shop := service.Company.GetCompany(admin.CompanyID)
			context.Session.Attributes.Put(play.SessionAdmin, admin)
			context.Session.Attributes.Put(play.SessionShop, shop)
		} else {
			as.Success = false
			as.Message = "手机/邮箱/密码不正确！"
		}

	}

	return &gweb.JsonResult{Data: as}
}
func loginUserAction(context *gweb.Context) gweb.Result {

	account := context.Request.FormValue("account")
	password := context.Request.FormValue("password")

	account = strings.ToLower(account) //小写

	as := &dao.ActionStatus{}

	user := service.User.FindUserByTel(service.Orm, account)

	context.Session.Attributes.Put(play.SessionUser, &dao.User{})

	if user.ID == 0 {
		as.Success = true
		as.Message = "账号/密码不正确！"
	} else {
		md5Password := tool.Md5(password)
		if strings.EqualFold(user.Password, md5Password) {
			as.Success = true
			as.Message = ""
			context.Session.Attributes.Put(play.SessionUser, user)
		} else {
			as.Success = true
			as.Message = "账号/密码不正确！"
		}

	}

	return &gweb.JsonResult{Data: as}
}
func loginManager(context *gweb.Context) gweb.Result {

	account := context.Request.FormValue("account")
	password := context.Request.FormValue("password")

	account = strings.ToLower(account) //小写

	as := &dao.ActionStatus{}

	user := service.Manager.FindManagerByAccount(account)

	if user.ID == 0 {
		as.Success = false
		as.Message = "账号/密码不正确！"
	} else {
		md5Password := tool.Md5(password)
		if strings.EqualFold(user.PassWord, md5Password) {
			as.Success = true
			as.Message = ""
			context.Session.Attributes.Put(play.SessionManager, user)
		} else {
			as.Success = false
			as.Message = "账号/密码不正确！"
		}
	}
	return &gweb.JsonResult{Data: as}
}
func openDo(context *gweb.Context) gweb.Result {
	fmt.Println(util.IsMobile(context))
	///account/open.do?redirect=%2Ffront%2Fappointment%2F20002%2Findex
	redirect := context.Request.URL.Query().Get("redirect")

	if util.IsMobile(context) {

	} else {
		http.Redirect(context.Response, context.Request, "/account/loginUserPage?redirect="+url.QueryEscape(redirect), http.StatusFound)
	}
	return &gweb.ViewResult{}
}
func forgetPage(context *gweb.Context) gweb.Result {
	return &gweb.HTMLResult{}
}
func loginAdminPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func loginManagerPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func loginUserPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func sdfsda(context *gweb.Context) gweb.Result {

	return &gweb.JsonResult{}
}
