package content

import (
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/dao"
	"strconv"

	"github.com/nbvghost/gweb"
)

type ContentController struct {
	gweb.BaseController
	Content service.ContentService
}

func (controller *ContentController) Init() {
	controller.AddHandler(gweb.GETMethod("{ContentItemID}/list/hot", controller.listHotAction))
	controller.AddHandler(gweb.GETMethod("{ContentItemID}/list/new", controller.listNewAction))
	controller.AddHandler(gweb.GETMethod("article/{ArticleID}", controller.articleAction))
	controller.AddHandler(gweb.GETMethod("{ContentItemID}/list/subtype", controller.ListContentSubTypeAction))
	controller.AddHandler(gweb.GETMethod("{ContentItemID}/related/{ContentSubTypeID}", controller.relatedAction))
}
func (controller *ContentController) ListContentSubTypeAction(context *gweb.Context) gweb.Result {
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	PID, _ := strconv.ParseUint(context.Request.URL.Query().Get("PID"), 10, 64)

	content := controller.Content.GetContentItemByID(ContentItemID)

	csts := controller.Content.FindContentSubTypesByContentItemIDAndParentContentSubTypeID(content.ID, PID)

	return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionOK, Message: "OK", Data: csts}}
}
func (controller *ContentController) articleAction(context *gweb.Context) gweb.Result {
	ArticleID, _ := strconv.ParseUint(context.PathParams["ArticleID"], 10, 64)
	//article := controller.Article.GetArticle(ArticleID)
	article := controller.Content.GetContentAndAddLook(context, ArticleID)
	return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionOK, Message: "OK", Data: article}}
}
func (controller *ContentController) relatedAction(context *gweb.Context) gweb.Result {
	Offset, _ := strconv.Atoi(context.Request.URL.Query().Get("Offset"))
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	ContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ContentSubTypeID"], 10, 64)

	var articles []dao.Content
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,ContentItemID,ContentSubTypeID,Author,Look,FromUrl", "CreatedAt desc", &articles, Offset, "ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID)
	return &gweb.JsonResult{Data: &result.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
func (controller *ContentController) listNewAction(context *gweb.Context) gweb.Result {
	Offset, _ := strconv.Atoi(context.Request.URL.Query().Get("Offset"))
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)

	var articles []dao.Content
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,ContentItemID,ContentSubTypeID,Author,Look,FromUrl", "CreatedAt desc", &articles, Offset, "ContentItemID=?", ContentItemID)
	return &gweb.JsonResult{Data: &result.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
func (controller *ContentController) listHotAction(context *gweb.Context) gweb.Result {
	Offset, _ := strconv.Atoi(context.Request.URL.Query().Get("Offset"))
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)

	var articles []dao.Content
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)

	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,ContentItemID,ContentSubTypeID,Author,Look,FromUrl", "Look desc", &articles, Offset, "ContentItemID=?", ContentItemID)

	return &gweb.JsonResult{Data: &result.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
