package view

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/funcmap"
	"github.com/nbvghost/dandelion/library/shop/embed/template/function"
	embedWidget "github.com/nbvghost/dandelion/library/shop/embed/template/widget"
	"github.com/nbvghost/dandelion/library/shop/view/active"
	"github.com/nbvghost/dandelion/library/shop/view/blog"
	"github.com/nbvghost/dandelion/library/shop/view/content"
	"github.com/nbvghost/dandelion/library/shop/view/page"
	"github.com/nbvghost/dandelion/library/shop/view/product"
	"github.com/nbvghost/dandelion/library/shop/view/user"
	"github.com/nbvghost/dandelion/library/shop/widget"
)

func Register(route constrain.IRoute) {

	//adminController := &admin.Controller{}
	//adminController.Interceptors.Set(&admin.Interceptor{})
	//adminController := gweb.NewController("admin", "")
	//adminController.NewController("template").DefaultHandle(&admin.Index{})
	//adminController.AddInterceptor(&admin.Interceptor{})

	route.RegisterView("search", &SearchRequest{})
	route.RegisterView("index", &IndexRequest{})
	route.RegisterView("privacy", &IndexRequest{})
	route.RegisterView("terms", &IndexRequest{})
	route.RegisterView("sitemap.xml", &SitemapRequest{})

	route.RegisterView("active/{event}", &active.Event{})

	{
		route.RegisterView("user/user", &IndexRequest{})
		route.RegisterView("user/sign-out", &user.SignOut{})

		route.RegisterView("user/{sub}/{page}", &user.IndexRequest{})
		route.RegisterView("user/{page}", &user.IndexRequest{})
	}

	{
		route.RegisterView("gallery/{TypeID}", &GalleryRequest{})
		route.RegisterView("gallery/{TypeID}/{SubTypeID}", &GalleryRequest{})
	}

	{
		route.RegisterView("blog/detail/{ContentID}", &blog.DetailRequest{})
		route.RegisterView("blog/{TypeID}", &BlogRequest{})
		route.RegisterView("blog/{TypeID}/{SubTypeID}", &BlogRequest{})
	}

	{
		/*templateNameData, err := dao.Find(db.GetDB(ctx), &model.ContentItem{}).Where(`"Type"=?`, model.ContentTypeContents).Group("TemplateName")
		if err != nil {
			panic(err)
		}
		templateNameList := templateNameData.([]string)
		for i := range templateNameList {
			templateName := templateNameList[i]
			route.RegisterView(fmt.Sprintf("%s/{TypeID}", templateName), &ContentsRequest{})
			route.RegisterView(fmt.Sprintf("%s/{TypeID}/{SubTypeID}", templateName), &ContentsRequest{})
		}*/
		route.RegisterView("contents/detail/{ContentUri}", &content.DetailRequest{})
		route.RegisterView("contents/{TypeUri}", &ContentsRequest{})
		route.RegisterView("contents/{TypeUri}/{SubTypeUri}", &ContentsRequest{})
	}
	{
		route.RegisterView("content/tag/{Tag}", &content.TagRequest{})
		//route.RegisterView("content/detail/{ContentID:[0-9]+}", &content.DetailRequest{})
		route.RegisterView("content/detail/{ContentUri}", &content.DetailRequest{})
		route.RegisterView("content/{TypeUri}", &ContentRequest{})
		route.RegisterView("content/{TypeUri}/{SubTypeUri}", &ContentRequest{})
	}
	{
		route.RegisterView("page/{Uri}", &page.PageRequest{})
		route.RegisterView("page/{Sub}/{Name}", &page.SubPageRequest{})
		//route.RegisterView("page/about", &page.AboutRequest{}, &p)
		//route.RegisterView("page/faq", &page.FaqRequest{}, &p)
		//route.RegisterView("page/privacy", &page.PrivacyRequest{}, &p)
		//route.RegisterView("page/terms", &page.TermsRequest{}, &p)
		//route.RegisterView("page/contact", &page.ContactRequest{}, &p)
	}
	{
		route.RegisterView("product/detail/{GoodsID}", &product.DetailRequest{})
		route.RegisterView("product/tag/{Tag}/page/{PageIndex}", &product.TagRequest{})
		route.RegisterView("product/tag/{Tag}", &product.TagRequest{})
		route.RegisterView("products", &ProductsRequest{})
		route.RegisterView("products/{TypeUri}", &ProductsRequest{})
		route.RegisterView("products/{TypeUri}/", &ProductsRequest{})
		route.RegisterView("products/{TypeUri}/{SubTypeUri}", &ProductsRequest{})
	}

	{
		route.RegisterView("favicon.ico", &FaviconicoRequest{})
		route.RegisterView("robots.txt", &RobotsRequest{})
		route.RegisterView("sign-in", &SignInRequest{})
		route.RegisterView("sign-up", &SignUpRequest{})
		route.RegisterView("404", &NotFound{})
		route.RegisterView("", &IndexRequest{})
		route.RegisterView("*", &DefaultRequest{})
	}

	//route.RegisterInterceptors("/")

	funcmap.RegisterFunction("OSSUrl", &function.OSSUrl{})
	funcmap.RegisterFunction("Title", &function.Title{})
	funcmap.RegisterFunction("Tags", &function.Tags{})
	funcmap.RegisterFunction("FaviconIco", &function.FaviconIco{})
	funcmap.RegisterFunction("GotoUrl", &function.GotoUrl{})
	funcmap.RegisterFunction("URLPathJoin", &function.URLPathJoin{})
	funcmap.RegisterFunction("Currency", &function.Currency{})
	funcmap.RegisterFunction("Url", &function.Url{})
	funcmap.RegisterFunction("UrlParam", &function.UrlParam{})
	funcmap.RegisterFunction("ToString", &function.ToString{})

	funcmap.RegisterWidget("GalleryBlock", &widget.GalleryBlock{})
	funcmap.RegisterWidget("Menus", &widget.Menus{})
	funcmap.RegisterWidget("Language", &widget.Language{})
	funcmap.RegisterWidget("Goods", &widget.Goods{})
	funcmap.RegisterWidget("LeaveMessage", &embedWidget.LeaveMessage{})
	funcmap.RegisterWidget("Subscribe", &widget.Subscribe{})
	funcmap.RegisterWidget("Footer", &widget.Footer{})
	funcmap.RegisterWidget("CustomerService", &widget.CustomerService{})
	funcmap.RegisterWidget("Breadcrumb", &widget.Breadcrumb{})
	funcmap.RegisterWidget("TopGoods", &widget.TopGoods{})
	funcmap.RegisterWidget("TopContent", &widget.TopContent{})
	funcmap.RegisterWidget("Social", &widget.Social{})
	funcmap.RegisterWidget("HtmlMeta", &widget.HtmlMeta{})
	funcmap.RegisterWidget("ContentPagination", &widget.ContentPagination[*model.Content]{})
	funcmap.RegisterWidget("GoodsPagination", &widget.ContentPagination[*extends.GoodsDetail]{})
	funcmap.RegisterWidget("ContentItemPagination", &widget.ContentItemPagination[*model.ContentItem]{})
	funcmap.RegisterWidget("SearchPagination", &widget.SearchPagination{})
	funcmap.RegisterWidget("TagPagination", &widget.TagPagination[*model.Content]{})
	funcmap.RegisterWidget("Categories", &widget.Categories{})

}

type RobotsRequest struct {
}

type RobotsReply struct {
	extends.ViewBase
}

func (m *RobotsRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &RobotsReply{}
	reply.Name = "robots.txt"
	reply.ContentType = "text/plain"
	return reply, nil
}
