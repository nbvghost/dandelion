package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Article struct {
	Get struct {
		ArticleID dao.PrimaryKey `form:"ArticleID"`
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

	article := service.Content.GetContentAndAddLook(ctx, dao.PrimaryKey(g.Get.ArticleID))
	return result.NewData(&result.ActionResult{Code: result.Success, Message: "OK", Data: article}), nil //{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: article}}

}
