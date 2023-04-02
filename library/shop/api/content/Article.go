package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/content"
	"github.com/nbvghost/gpa/types"
)

type Article struct {
	ContentService content.ContentService
	Get            struct {
		ArticleID types.PrimaryKey `form:"ArticleID"`
	} `method:"Get"`
}

func (g *Article) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Article) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//ArticleID, _ := strconv.ParseUint(context.PathParams["ArticleID"], 10, 64)
	//ArticleID := object.ParseUint(context.PathParams["ArticleID"])
	//article := controller.Article.GetArticle(ArticleID)

	article := g.ContentService.GetContentAndAddLook(ctx, types.PrimaryKey(g.Get.ArticleID))
	return result.NewData(&result.ActionResult{Code: result.Success, Message: "OK", Data: article}), nil //{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: article}}

}
