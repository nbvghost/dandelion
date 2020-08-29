package content

import (
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/dao"
	"strconv"

	"github.com/nbvghost/gweb"
)

type ContentController struct {
	gweb.BaseController
	Content service.ContentService
	Article service.ArticleService
}

func (controller *ContentController) Init() {
	controller.AddHandler(gweb.GETMethod("{ContentID}/list/hot", controller.listHotAction))
	controller.AddHandler(gweb.GETMethod("{ContentID}/list/new", controller.listNewAction))
	controller.AddHandler(gweb.GETMethod("article/{ArticleID}", controller.articleAction))
	controller.AddHandler(gweb.GETMethod("{ContentID}/list/subtype", controller.Content.ListContentSubTypeAction))
	controller.AddHandler(gweb.GETMethod("{ContentID}/related/{ContentSubTypeID}", controller.relatedAction))
}
func (controller *ContentController) articleAction(context *gweb.Context) gweb.Result {
	ArticleID, _ := strconv.ParseUint(context.PathParams["ArticleID"], 10, 64)
	//article := controller.Article.GetArticle(ArticleID)
	article := controller.Article.GetArticleAndAddLook(context, ArticleID)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: article}}
}
func (controller *ContentController) relatedAction(context *gweb.Context) gweb.Result {
	Offset, _ := strconv.Atoi(context.Request.URL.Query().Get("Offset"))
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)
	ContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ContentSubTypeID"], 10, 64)

	var articles []dao.Article
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,ContentID,ContentSubTypeID,Author,Look,FromUrl", "CreatedAt desc", &articles, Offset, "ContentID=? and ContentSubTypeID=?", ContentID, ContentSubTypeID)
	return &gweb.JsonResult{Data: &dao.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
func (controller *ContentController) listNewAction(context *gweb.Context) gweb.Result {
	Offset, _ := strconv.Atoi(context.Request.URL.Query().Get("Offset"))
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)

	var articles []dao.Article
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)
	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,ContentID,ContentSubTypeID,Author,Look,FromUrl", "CreatedAt desc", &articles, Offset, "ContentID=?", ContentID)
	return &gweb.JsonResult{Data: &dao.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
func (controller *ContentController) listHotAction(context *gweb.Context) gweb.Result {
	Offset, _ := strconv.Atoi(context.Request.URL.Query().Get("Offset"))
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)

	var articles []dao.Article
	//controller.Content.FindOrderWhereLength(dao.Orm(),"Look desc",&articles,)

	_Total, _Limit, _Offset := controller.Content.FindSelectWherePaging(dao.Orm(), "ID,Title,Picture,ContentID,ContentSubTypeID,Author,Look,FromUrl", "Look desc", &articles, Offset, "ContentID=?", ContentID)

	return &gweb.JsonResult{Data: &dao.Pager{Data: articles, Total: _Total, Limit: _Limit, Offset: _Offset}}
}
