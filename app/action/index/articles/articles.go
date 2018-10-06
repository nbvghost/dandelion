package articles

import (
	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/service/dao"
	"net/http"
	"strconv"

	"github.com/nbvghost/gweb"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nbvghost/gweb/tool"
)

type InterceptorMp struct {
	Organization service.OrganizationService
}

func (controller InterceptorMp) Execute(Context *gweb.Context) (bool, gweb.Result) {
	return true, nil
}

type Controller struct {
	gweb.BaseController
	Content service.ContentService
	Article service.ArticleService
	Wx      service.WxService
	User    service.UserService
}

func (controller *Controller) Apply() {
	//Index.RequestMapping = make(map[string]mvc.Function)
	controller.Interceptors.Add(&InterceptorMp{})
	controller.AddHandler(gweb.ALLMethod("", controller.defaultPage))
	controller.AddHandler(gweb.ALLMethod("*", controller.indexPage))

	controller.AddHandler(gweb.GETMethod("list/:OID/new", controller.listNewAction))
	controller.AddHandler(gweb.GETMethod("article", controller.articleSeftPage))

	//csts := controller.Content.FindContentSubTypesByContentID(content.ID)
	controller.AddHandler(gweb.GETMethod("index", controller.indexPage))
	controller.AddHandler(gweb.GETMethod("content/:ContentID/index", controller.contentPage))
	controller.AddHandler(gweb.GETMethod("content/:ContentID/article/:ArticleID", controller.articlePage))

	controller.AddHandler(gweb.GETMethod("content/:ContentID/list/sub/new/:ContentSubTypeID", controller.listSubContentNewAction))
	controller.AddHandler(gweb.GETMethod("content/:ContentID/list/hot", controller.listContentHotAction))
	controller.AddHandler(gweb.GETMethod("content/:ContentID/list/new", controller.listContentNewAction))

	controller.AddHandler(gweb.GETMethod("content/:ContentID/list/subtype", controller.Content.ListContentSubTypeAction))

	controller.AddHandler(gweb.GETMethod("get/:ArticleID", controller.getArticleAction))
	controller.AddHandler(gweb.GETMethod("new_list", controller.listNewPage))
}
func (controller *Controller) articleSeftPage(context *gweb.Context) gweb.Result {
	ArticleID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
	article := controller.Article.GetArticle(ArticleID)
	//article.Content=template.HTML(article.Content)
	return &gweb.HTMLResult{Params: map[string]interface{}{"Article": article}}
}
func (controller *Controller) listNewAction(context *gweb.Context) gweb.Result {
	Offset, _ := strconv.ParseInt(context.Request.URL.Query().Get("Offset"), 10, 64)
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	OID, _ := strconv.ParseUint(context.PathParams["OID"], 10, 64)
	ContentIDs := controller.Content.GetContentIDs(OID)

	var articles []dao.Article
	//controller.Content.FindOrderWhereLength(service.Orm(),"Look desc",&articles,)
	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Thumbnail,ContentID,ContentSubTypeID,ContentSubTypeChildID,Author,Look,FromUrl", "CreatedAt desc", &articles, int(Offset), "ContentID in (?)", ContentIDs)
	return &gweb.JsonResult{Data: &dao.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
func (controller *Controller) articlePage(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	ArticleID, _ := strconv.ParseUint(context.PathParams["ArticleID"], 10, 64)
	article := controller.Article.GetArticleAndAddLook(context, ArticleID)
	//article.Content=template.HTML(article.Content)

	//fmt.Println(controller.Wx.MWQRCodeTemp(company.ID, 145, play.QRCodeCreateType_Article, strconv.Itoa(int(ArticleID))))

	//todo:test
	if context.Session.Attributes.Get(play.SessionUser) != nil {

		user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

		if user.Subscribe == 0 {
			UserAgent := strings.ToLower(context.Request.Header.Get("User-Agent"))
			if strings.Contains(UserAgent, "micromessenger") {

				WxConfig := controller.Wx.MiniWeb()
				access_token := controller.Wx.GetAccessToken(WxConfig)
				//https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
				resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + access_token + "&openid=" + user.OpenID + "&lang=zh_CN")
				tool.Trace(err)
				b, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()

				fmt.Println(string(b))
				//controller.User.ChangeModel(dao.Orm(), user.ID, user)

				mapData := make(map[string]interface{})
				err = json.Unmarshal(b, &mapData)

				if mapData["errcode"] != nil {
					return &gweb.HTMLResult{Name: "articles/article", Params: map[string]interface{}{"Article": article}}
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
				context.Session.Attributes.Put(play.SessionUser, user)

			}
		}

	}

	return &gweb.HTMLResult{Name: "articles/article", Params: map[string]interface{}{"Article": article}}
}
func (controller *Controller) listNewPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func (controller *Controller) getArticleAction(context *gweb.Context) gweb.Result {
	ArticleID, _ := strconv.ParseUint(context.PathParams["ArticleID"], 10, 64)
	article := controller.Article.GetArticle(ArticleID)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: article}}
}
func (controller *Controller) listSubContentNewAction(context *gweb.Context) gweb.Result {
	//ContentSubTypeID
	ContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ContentSubTypeID"], 10, 64)
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)

	Offset, _ := strconv.ParseInt(context.Request.URL.Query().Get("Offset"), 10, 64)
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	//ContentIDs:=controller.Content.GetContentIDs(Organization.ID)

	var articles []dao.Article
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Thumbnail,Introduce,ContentID,ContentSubTypeID,Author,Look,FromUrl,CreatedAt", "CreatedAt desc", &articles, int(Offset), "ContentID=? and ContentSubTypeID=?", ContentID, ContentSubTypeID)
	return &gweb.JsonResult{Data: &dao.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
func (controller *Controller) listContentNewAction(context *gweb.Context) gweb.Result {
	//ContentID
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)
	Offset, _ := strconv.ParseInt(context.Request.URL.Query().Get("Offset"), 10, 64)
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	//ContentIDs:=controller.Content.GetContentIDs(Organization.ID)
	var articles []dao.Article
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Thumbnail,Introduce,ContentID,ContentSubTypeID,Author,Look,FromUrl,CreatedAt", "CreatedAt desc", &articles, int(Offset), "ContentID=?", ContentID)
	return &gweb.JsonResult{Data: &dao.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
func (controller *Controller) listContentHotAction(context *gweb.Context) gweb.Result {
	//ContentID
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)
	Offset, _ := strconv.ParseInt(context.Request.URL.Query().Get("Offset"), 10, 64)
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	//ContentIDs:=controller.Content.GetContentIDs(Organization.ID)
	var articles []dao.Article
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Thumbnail,Introduce,ContentID,ContentSubTypeID,Author,Look,FromUrl,CreatedAt", "Look desc", &articles, int(Offset), "ContentID=?", ContentID)
	return &gweb.JsonResult{Data: &dao.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {
	return &gweb.RedirectToUrlResult{Url: "index"}
}

//6c0420c5e926a2ac8d56aa4192ab10fa
func (controller *Controller) indexPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func (controller *Controller) contentPage(context *gweb.Context) gweb.Result {
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)

	content := controller.Content.GetContentByID(ContentID)

	csts := controller.Content.FindContentSubTypesByContentID(content.ID)

	result := make(map[string]interface{})
	result["ContentSubTypes"] = csts
	result["Content"] = content

	return &gweb.HTMLResult{Name: "articles/content", Params: result}
}
