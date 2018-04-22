package manager

import (
	"net/http"
	"net/url"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"

	"strconv"

	"encoding/json"

	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/service/dao"
)

type InterceptorManager struct {
}

func (this InterceptorManager) Execute(context *gweb.Context) bool {

	//util.Trace(context.Session,"context.Session")
	if context.Session.Attributes.Get(play.SessionManager) == nil {
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
		http.Redirect(context.Response, context.Request, "/account/loginManagerPage?redirect="+url.QueryEscape(redirect), http.StatusFound)
		return false
	} else {
		return true
	}
}

type Controller struct {
	gweb.BaseController
}

func (i *Controller) Apply() {
	i.Interceptors.Add(&InterceptorManager{})
	//Index.RequestMapping = make(map[string]mvc.Function)
	i.AddHandler(gweb.ALLMethod("", rootPage))
	i.AddHandler(gweb.ALLMethod("*", defaultPage))
	i.AddHandler(gweb.ALLMethod("index", indexPage))
	i.AddHandler(gweb.ALLMethod("articlePage", articlePage))
	i.AddHandler(gweb.ALLMethod("article", articleAction))
	i.AddHandler(gweb.ALLMethod("admin", adminAction))
	i.AddHandler(gweb.ALLMethod("add_article", addArticlePage))
	i.AddHandler(gweb.ALLMethod("categoryAction", categoryAction))

}
func addArticlePage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{Params: nil}
}
func categoryAction(context *gweb.Context) gweb.Result {
	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	case "list":
		result.Data = &dao.ActionStatus{true, "ok", service.Category.FindCategory()}
	case "add":
		context.Request.ParseForm()
		label := context.Request.Form.Get("label")
		_, su := service.Category.AddCategory(label)
		result.Data = (&dao.ActionStatus{}).Smart(su, "添加成功", "添加失败,")
	case "del":
		context.Request.ParseForm()
		id, _ := strconv.ParseUint(context.Request.Form.Get("id"), 10, 64)
		service.Category.DelCategory(id)
		result.Data = &dao.ActionStatus{true, "删除成功", nil}
	}

	return result
}

func adminAction(context *gweb.Context) gweb.Result {
	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	//list del
	case play.ActionKey_list:
		result.Data = &dao.ActionStatus{true, "ok", service.Admin.FindAdmin(service.Orm)}
	case play.ActionKey_del:
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		result.Data = service.Admin.DelAdmin(ID)
	case play.ActionKey_get:
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		admin := service.Admin.GetAdmin(service.Orm, ID)
		result.Data = &dao.ActionStatus{true, "ok", admin}
	case play.ActionKey_change:
		context.Request.ParseForm()
		Name := context.Request.Form.Get("Name")
		Password := context.Request.Form.Get("Password")
		Email := context.Request.Form.Get("Email")
		Tel := context.Request.Form.Get("Tel")
		ID, _ := strconv.ParseUint(context.Request.Form.Get("ID"), 10, 64)

		if err := service.Admin.ChangeAdmin(Name, Password, Email, Tel, ID); err != nil {
			result.Data = &dao.ActionStatus{false, "修改失败", nil}
		} else {
			result.Data = &dao.ActionStatus{true, "修改成功", nil}
		}

	case play.ActionKey_add:
		context.Request.ParseForm()
		Name := context.Request.Form.Get("Name")
		Password := context.Request.Form.Get("Password")
		Email := context.Request.Form.Get("Email")
		Tel := context.Request.Form.Get("Tel")
		result.Data = service.Admin.AddAdmin(Name, Password, Email, Tel)
	}
	return result
}
func articleAction(context *gweb.Context) gweb.Result {

	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	case "listByCategory":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		d := service.Article.FindArticleByCategoryID(ID)
		result.Data = &dao.ActionStatus{true, "SUCCESS", d}
	case "add":
		context.Request.ParseForm()
		jsonText := context.Request.Form.Get("json")
		article := &dao.Article{}
		err := json.Unmarshal([]byte(jsonText), article)
		tool.CheckError(err)
		result.Data = service.Article.AddArticle(article)
	case "change":
		context.Request.ParseForm()
		jsonText := context.Request.Form.Get("json")
		article := &dao.Article{}
		err := json.Unmarshal([]byte(jsonText), article)
		tool.CheckError(err)
		service.Article.ChangeArticle(article)
		result.Data = &dao.ActionStatus{true, "SUCCESS", article}
	case "one":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		article := service.Article.GetArticle(ID)
		result.Data = &dao.ActionStatus{true, "SUCCESS", article}
	case "del":
		context.Request.ParseForm()
		id, _ := strconv.ParseUint(context.Request.Form.Get("id"), 10, 64)
		service.Article.DelArticle(id)
		result.Data = &dao.ActionStatus{true, "删除成功", nil}
	}

	return result
}
func articlePage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}

func indexPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}

func defaultPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func rootPage(context *gweb.Context) gweb.Result {

	return &gweb.RedirectToUrlResult{"index"}
}
