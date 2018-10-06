package index

import (
	"dandelion/app/action/index/articles"
	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"fmt"
	"io/ioutil"

	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type InterceptorMp struct {
	Organization service.OrganizationService
	Wx           service.WxService
}

//Execute(Session *Session,Request *http.Request)(bool,Result)
func (controller InterceptorMp) Execute(context *gweb.Context) (bool, gweb.Result) {
	/*su, re := controller.Organization.ReadOrganization(context)
	if su {
		Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
		wxconfig := controller.Wx.MiniWeb(Organization.ID)
		if !strings.EqualFold(wxconfig.AppID, "") && !strings.EqualFold(wxconfig.AppSecret, "") {
			//micromessenger
			UserAgent := strings.ToLower(context.Request.Header.Get("User-Agent"))
			if strings.Contains(UserAgent, "micromessenger") {

				if context.Session.Attributes.Get(play.SessionUser) == nil {

					redirect := util.GetFullPath(context.Request)
					return false, &gweb.RedirectToUrlResult{Url: "/account/wx/authorize?OID=" + strconv.Itoa(int(Organization.ID)) + "&redirect=" + redirect}
				}
			}
		}
	}
	context.Response.Header().Set("Access-Control-Allow-Origin", "*")
	return su, re*/
	return true, nil
}

type Controller struct {
	gweb.BaseController
	Content       service.ContentService
	Article       service.ArticleService
	Configuration service.ConfigurationService
}

func (controller *Controller) Apply() {
	//Index.RequestMapping = make(map[string]mvc.Function)
	controller.Interceptors.Add(&InterceptorMp{})
	controller.AddHandler(gweb.ALLMethod("", controller.defaultPage))
	controller.AddHandler(gweb.ALLMethod("*", controller.indexPage))
	controller.AddHandler(gweb.GETMethod("index", controller.indexPage))
	controller.AddHandler(gweb.POSMethod("/content/new/article/webhook", controller.newArticleWebhookAction))
	controller.AddHandler(gweb.POSMethod("/content/new/article/post", controller.newArticlePostAction))
	controller.AddHandler(gweb.POSMethod("/configuration/list", controller.configurationListAction))
	//https://dandelion.nutsy.cc/2000/content/2000/new/article/webhook

	articles := &articles.Controller{}
	articles.Interceptors = controller.Interceptors
	controller.AddSubController("/articles/", articles)
}
func (controller *Controller) configurationListAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var ks []uint64
	util.RequestBodyToJSON(context.Request.Body, &ks)
	list := controller.Configuration.GetConfigurations(company.ID, ks)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}
}
func (controller *Controller) newArticlePostAction(context *gweb.Context) gweb.Result {
	//ContentID,_:=strconv.ParseUint(context.PathParams["ContentID"],10,64)
	organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	context.Request.ParseForm()
	__sign := context.Request.FormValue("__sign")
	/*context.Request.ParseForm()

	__dataId:=context.Request.FormValue("__dataId")
	__crawlTime:=context.Request.FormValue("__crawlTime")
	__dataVersion:=context.Request.FormValue("__dataVersion")
	__dataLatest:=context.Request.FormValue("__dataLatest")
	__sourceId:=context.Request.FormValue("__sourceId")
	__dataKey:=context.Request.FormValue("__dataKey")
	__crawlUrl:=context.Request.FormValue("__crawlUrl")*/

	msSign := tool.Md5ByString("274455411" + "shenjianshou.cn")
	if strings.EqualFold(strings.ToUpper(__sign), msSign) {

		article_title := context.Request.FormValue("article_title")
		weixin_tmp_url := context.Request.FormValue("weixin_tmp_url")
		article_publish_time := context.Request.FormValue("article_publish_time")
		weixin_introduce := context.Request.FormValue("weixin_introduce")
		article_thumbnail := context.Request.FormValue("article_thumbnail")
		article_content := context.Request.FormValue("article_content")

		//var images []interface{}
		//article_images:=context.Request.FormValue("article_images")
		//json.Unmarshal([]byte(article_images),&images)

		weixin_nickname := context.Request.FormValue("weixin_nickname")
		//fmt.Println(article_title,weixin_tmp_url,weixin_introduce,article_thumbnail,article_content,images,weixin_nickname)

		var createTime time.Time
		ts, _ := strconv.ParseInt(article_publish_time, 10, 64)
		createTime = time.Unix(ts, 0)

		//OID,			   ContentName,ContentSubTypeName,Author,Title,          FromUrl,       Introduce,       Thumbnail,                        Content,        CreatedAt
		controller.addArticle(organization.ID, "头条文摘", weixin_nickname, weixin_nickname, article_title, weixin_tmp_url, weixin_introduce, article_thumbnail, article_content, createTime)

		return &gweb.JsonResult{Data: map[string]interface{}{"result": 1, "data": "发布成功"}}

	} else {
		return &gweb.JsonResult{Data: map[string]interface{}{"result": 2, "reason": "发布失败, 错误原因: 发布密码验证失败"}}
	}

	//
}
func (controller *Controller) addArticle(OID uint64, ContentName string, ContentSubTypeName string, Author, Title string, FromUrl string, Introduce string, Thumbnail string, Content string, CreatedAt time.Time) {

	controller.Article.AddSpiderArticle(OID, ContentName, ContentSubTypeName, Author, Title, FromUrl, Introduce, Thumbnail, Content, CreatedAt)

}
func (controller *Controller) newArticleWebhookAction(context *gweb.Context) gweb.Result {
	//ContentID,_:=strconv.ParseUint(context.PathParams["ContentID"],10,64)
	organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	context.Request.ParseForm()
	url := context.Request.FormValue("url")
	data := context.Request.FormValue("data")
	timestamp := context.Request.FormValue("timestamp")
	crawl_time := context.Request.FormValue("crawl_time")
	sign2 := context.Request.FormValue("sign2")
	data_key := context.Request.FormValue("data_key")
	event_type := context.Request.FormValue("event_type")

	mySign2 := tool.Md5ByString(url + "k4Mjg3NDhmYTJlYT-a1910cf788d26ec" + timestamp)
	fmt.Println(data, crawl_time, data_key, event_type)
	fmt.Println(mySign2)
	fmt.Println(sign2)
	if strings.EqualFold(mySign2, strings.ToUpper(sign2)) {
		dataMap := make(map[string]interface{})
		util.JSONToStruct(data, &dataMap)

		//article_publish_time :=dataMap["article_publish_time"].(string)//context.Request.FormValue("article_publish_time")

		var weixin_nickname string
		var article_title string
		var weixin_tmp_url string
		var weixin_introduce string
		var article_thumbnail string
		var article_content string
		var article_publish_time string

		if dataMap["weixin_nickname"] != nil {
			weixin_nickname = dataMap["weixin_nickname"].(string)
		}
		if dataMap["article_title"] != nil {
			article_title = dataMap["article_title"].(string)
		}
		if dataMap["weixin_tmp_url"] != nil {
			weixin_tmp_url = dataMap["weixin_tmp_url"].(string)
		}
		if dataMap["weixin_introduce"] != nil {
			weixin_introduce = dataMap["weixin_introduce"].(string)
		}
		if dataMap["article_thumbnail"] != nil {
			article_thumbnail = dataMap["article_thumbnail"].(string)
		}
		if dataMap["article_content"] != nil {
			article_content = dataMap["article_content"].(string)
		}
		if dataMap["article_publish_time"] != nil {
			article_publish_time = dataMap["article_publish_time"].(string)
		}

		var createTime time.Time
		ts, _ := strconv.ParseInt(article_publish_time, 10, 64)
		createTime = time.Unix(ts, 0)

		//controller.addArticle(organization.ID,"头条文摘",weixin_nickname,weixin_nickname,article_title,weixin_tmp_url,weixin_introduce,article_thumbnail,article_content,createTime)
		controller.addArticle(organization.ID, "头条文摘", weixin_nickname, weixin_nickname, article_title, weixin_tmp_url, weixin_introduce, article_thumbnail, article_content, createTime)
		return &gweb.TextResult{Data: data_key}
	} else {
		return &gweb.TextResult{Data: ""}
	}
	b, err := ioutil.ReadAll(context.Request.Body)
	fmt.Println(err)
	fmt.Println(string(b))

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: nil}}
}

func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {
	return &gweb.RedirectToUrlResult{Url: "index"}
}

//6c0420c5e926a2ac8d56aa4192ab10fa
func (controller *Controller) indexPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
