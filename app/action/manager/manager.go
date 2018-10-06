package manager

import (
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

//Execute(Session *Session,Request *http.Request)(bool,Result)
func (this InterceptorManager) Execute(context *gweb.Context) (bool, gweb.Result) {

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
		//http.Redirect(context.Response, context.Request, "/account/loginManagerPage?redirect="+url.QueryEscape(redirect), http.StatusFound)
		return false, &gweb.RedirectToUrlResult{Url: "/account/loginManagerPage?redirect=" + url.QueryEscape(redirect)}
	} else {
		return true, nil
	}
}

type Controller struct {
	gweb.BaseController
	Article service.ArticleService
	Admin   service.AdminService
}

func (this *Controller) Apply() {
	//Index.RequestMapping = make(map[string]mvc.Function)
	this.Interceptors.Add(&InterceptorManager{})
	this.AddHandler(gweb.ALLMethod("", this.rootPage))
	this.AddHandler(gweb.ALLMethod("*", this.defaultPage))
	this.AddHandler(gweb.ALLMethod("index", this.indexPage))
	this.AddHandler(gweb.ALLMethod("articlePage", this.articlePage))
	this.AddHandler(gweb.ALLMethod("article", this.articleAction))
	this.AddHandler(gweb.ALLMethod("admin", this.adminAction))
	this.AddHandler(gweb.ALLMethod("add_article", this.addArticlePage))
	//this.AddHandler(gweb.ALLMethod("categoryAction", this.categoryAction))

}
func (this *Controller) addArticlePage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{Params: nil}
}

/*func (this *Controller) categoryAction(context *gweb.Context) gweb.Result {
	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	case "list":
		result.Data = &dao.ActionStatus{true, "ok", this.ArticleType.FindCategory()}
	case "add":
		context.Request.ParseForm()
		label := context.Request.Form.Get("label")
		_, su := this.ArticleType.AddCategory(label)
		result.Data = (&dao.ActionStatus{}).Smart(su, "添加成功", "添加失败,")
	case "del":
		context.Request.ParseForm()
		id, _ := strconv.ParseUint(context.Request.Form.Get("id"), 10, 64)
		this.ArticleType.DelCategory(id)
		result.Data = &dao.ActionStatus{true, "删除成功", nil}
	}

	return result
}*/

func (this *Controller) adminAction(context *gweb.Context) gweb.Result {
	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	//list del
	case play.ActionKey_list:
		result.Data = &dao.ActionStatus{true, "ok", this.Admin.FindAdmin()}
	case play.ActionKey_del:
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		result.Data = this.Admin.DelAdmin(ID)
	case play.ActionKey_get:
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		admin := this.Admin.GetAdmin(ID)
		result.Data = &dao.ActionStatus{true, "ok", admin}
	case play.ActionKey_change:
		context.Request.ParseForm()
		Name := context.Request.Form.Get("Account")
		Password := context.Request.Form.Get("PassWord")
		//Email := context.Request.Form.Get("Email")
		//Tel := context.Request.Form.Get("Tel")
		ID, _ := strconv.ParseUint(context.Request.Form.Get("ID"), 10, 64)

		if err := this.Admin.ChangeAdmin(Name, Password, ID); err != nil {
			result.Data = &dao.ActionStatus{false, "修改失败", nil}
		} else {
			result.Data = &dao.ActionStatus{true, "修改成功", nil}
		}

	case play.ActionKey_add:
		context.Request.ParseForm()
		Name := context.Request.Form.Get("Account")
		Password := context.Request.Form.Get("PassWord")
		Domain := context.Request.Form.Get("Domain")
		//Email := context.Request.Form.Get("Email")
		//Tel := context.Request.Form.Get("Tel")
		result.Data = this.Admin.AddAdmin(Name, Password, Domain)
	}
	return result
}
func (this *Controller) articleAction(context *gweb.Context) gweb.Result {

	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	case "listByCategory":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		d := this.Article.FindArticleByContentSubTypeID(ID)
		result.Data = &dao.ActionStatus{Success: true, Message: "SUCCESS", Data: d}
	case "add":
		context.Request.ParseForm()
		jsonText := context.Request.Form.Get("json")
		article := &dao.Article{}
		err := json.Unmarshal([]byte(jsonText), article)
		tool.CheckError(err)
		result.Data = this.Article.AddArticle(article)
	case "change":
		context.Request.ParseForm()
		jsonText := context.Request.Form.Get("json")
		article := &dao.Article{}
		err := json.Unmarshal([]byte(jsonText), article)
		tool.CheckError(err)
		this.Article.ChangeArticle(article)
		result.Data = &dao.ActionStatus{true, "SUCCESS", article}
	case "one":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		article := this.Article.GetArticle(ID)
		result.Data = &dao.ActionStatus{true, "SUCCESS", article}
	case "del":
		context.Request.ParseForm()
		id, _ := strconv.ParseUint(context.Request.Form.Get("id"), 10, 64)
		this.Article.DelArticle(id)
		result.Data = &dao.ActionStatus{true, "删除成功", nil}
	}

	return result
}
func (this *Controller) articlePage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}

func (this *Controller) indexPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}

func (this *Controller) defaultPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func (this *Controller) rootPage(context *gweb.Context) gweb.Result {

	return &gweb.RedirectToUrlResult{"index"}
}
