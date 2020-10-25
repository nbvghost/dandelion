package articles

import (
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/dao"
	"net/http"
	"strconv"

	"github.com/nbvghost/glog"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nbvghost/gweb"
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
	Wx      service.WxService
	User    service.UserService
}

func (controller *Controller) Init() {
	//Index.RequestMapping = make(map[string]mvc.Function)
	controller.Interceptors.Add(&InterceptorMp{})

	controller.AddHandler(gweb.GETMethod("list/{OID}/new", controller.listNewAction))
	controller.AddHandler(gweb.GETMethod("article", controller.articleSeftPage))

	controller.AddHandler(gweb.GETMethod("index", controller.indexPage))
	controller.AddHandler(gweb.GETMethod("content/{ContentItemID}/index", controller.contentPage))
	controller.AddHandler(gweb.GETMethod("content/{ContentItemID}/article/{ArticleID}", controller.articlePage))
	controller.AddHandler(gweb.GETMethod("content/{ContentItemID}/list/sub/new/{ContentSubTypeID}", controller.listSubContentNewAction))
	controller.AddHandler(gweb.GETMethod("content/{ContentItemID}/list/hot", controller.listContentHotAction))
	controller.AddHandler(gweb.GETMethod("content/{ContentItemID}/list/new", controller.listContentNewAction))
	controller.AddHandler(gweb.GETMethod("content/{ContentItemID}/list/subtype", controller.ListContentSubTypeAction))

	controller.AddHandler(gweb.GETMethod("get/{ArticleID}", controller.getArticleAction))
	controller.AddHandler(gweb.GETMethod("new_list", controller.listNewPage))
}
func (controller *Controller) ListContentSubTypeAction(context *gweb.Context) gweb.Result {
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	PID, _ := strconv.ParseUint(context.Request.URL.Query().Get("PID"), 10, 64)

	content := controller.Content.GetContentItemByID(ContentItemID)

	csts := controller.Content.FindContentSubTypesByContentItemIDAndParentContentSubTypeID(content.ID, PID)

	return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionOK, Message: "OK", Data: csts}}
}
func (controller *Controller) articleSeftPage(context *gweb.Context) gweb.Result {
	ArticleID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
	article := controller.Content.GetContentByID(ArticleID)
	//article.Content=template.HTML(article.Content)
	return &gweb.HTMLResult{Params: map[string]interface{}{"Article": article}}
}
func (controller *Controller) listNewAction(context *gweb.Context) gweb.Result {
	Offset, _ := strconv.ParseInt(context.Request.URL.Query().Get("Offset"), 10, 64)
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	OID, _ := strconv.ParseUint(context.PathParams["OID"], 10, 64)
	ContentItemIDs := controller.Content.GetContentItemIDs(OID)

	//var articles []dao.Content
	//controller.Content.FindOrderWhereLength(service.Orm(),"Look desc",&articles,)
	pagin := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,ContentItemID,ContentSubTypeID,Author,Look,FromUrl", "CreatedAt desc", dao.Content{}, int(Offset), "ContentItemID in (?)", ContentItemIDs)
	return &gweb.JsonResult{Data: &pagin}
}
func (controller *Controller) articlePage(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	ArticleID, _ := strconv.ParseUint(context.PathParams["ArticleID"], 10, 64)
	article := controller.Content.GetContentAndAddLook(context, ArticleID)
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
				glog.Trace(err)

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
	article := controller.Content.GetContentByID(ArticleID)
	return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionOK, Message: "OK", Data: article}}
}
func (controller *Controller) listSubContentNewAction(context *gweb.Context) gweb.Result {
	//ContentSubTypeID
	ContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ContentSubTypeID"], 10, 64)
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)

	Offset, _ := strconv.ParseInt(context.Request.URL.Query().Get("Offset"), 10, 64)
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	//var articles []dao.Content
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	pagin := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,Introduce,ContentItemID,ContentSubTypeID,Author,Look,FromUrl,CreatedAt", "CreatedAt desc", dao.Content{}, int(Offset), "ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID)
	return &gweb.JsonResult{Data: &pagin}
}
func (controller *Controller) listContentNewAction(context *gweb.Context) gweb.Result {

	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	Offset, _ := strconv.ParseInt(context.Request.URL.Query().Get("Offset"), 10, 64)
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	//var articles []dao.Content
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	pagin := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,Introduce,ContentItemID,ContentSubTypeID,Author,Look,FromUrl,CreatedAt", "CreatedAt desc", dao.Content{}, int(Offset), "ContentItemID=?", ContentItemID)
	return &gweb.JsonResult{Data: &pagin}
}
func (controller *Controller) listContentHotAction(context *gweb.Context) gweb.Result {
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	Offset, _ := strconv.ParseInt(context.Request.URL.Query().Get("Offset"), 10, 64)
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	//var articles []dao.Content
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	pager := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,Introduce,ContentItemID,ContentSubTypeID,Author,Look,FromUrl,CreatedAt", "Look desc", dao.Content{}, int(Offset), "ContentItemID=?", ContentItemID)
	return &gweb.JsonResult{Data: &pager}
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

	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)

	content := controller.Content.GetContentItemByID(ContentItemID)

	csts := controller.Content.FindContentSubTypesByContentItemID(content.ID)

	result := make(map[string]interface{})
	result["ContentSubTypes"] = csts
	result["Content"] = content

	return &gweb.HTMLResult{Name: "articles/content", Params: result}
}
