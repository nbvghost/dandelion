package web

import (
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/sites"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool/number"
)

type Controller struct {
	gweb.BaseController
	Template sites.TemplateService
	Content  service.ContentService
}

func (controller *Controller) Init() {
	controller.AddHandler(gweb.ALLMethod("/index", controller.index))
	controller.AddHandler(gweb.ALLMethod("/gallery", controller.gallery))
	controller.AddHandler(gweb.ALLMethod("/contents", controller.contents))
	controller.AddHandler(gweb.ALLMethod("/js/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/css/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/img/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/font/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/lib/", controller.AddProjectdsfdsfsdAction))
	//controller.AddHandler(gweb.ALLMethod("/", controller.defaultPage))
}

func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {

	//return &gweb.HTMLResult{}
	return &gweb.RedirectToUrlResult{Url: "index"}

}
func (controller *Controller) contents(context *gweb.Context) gweb.Result {
	params := make(map[string]interface{})

	ContentItemID := number.ParseInt(context.Request.URL.Query().Get("id"))
	ContentSubTypeID := number.ParseInt(context.Request.URL.Query().Get("sid"))

	//item := controller.Content.GetContentItemByID(uint64(ContentItemID))
	//params["Item"] = item

	var ContentSubTypeIDList []uint64
	if ContentSubTypeID == 0 {
		ContentSubTypeIDList = []uint64{}
	} else {
		ContentSubTypeIDList = controller.Content.GetContentSubTypeAllIDByID(uint64(ContentItemID), uint64(ContentSubTypeID))
	}

	contentList := controller.Content.FindContentByContentItemIDAndContentSubTypeID(uint64(ContentItemID), ContentSubTypeIDList)
	params["ContentList"] = contentList

	item, menusPath := controller.Template.MenusTemplate(context, uint64(ContentItemID), uint64(ContentSubTypeID), params)
	commonPath := controller.Template.CommonTemplate(context, params)

	return &gweb.HTMLResult{
		Name: item.TemplateName,
		Template: []string{
			menusPath, commonPath,
		},
		Params: params,
	}
}
func (controller *Controller) gallery(context *gweb.Context) gweb.Result {

	//siteName := context.PathParams["siteName"]

	params := make(map[string]interface{})

	_, menusPath := controller.Template.MenusTemplate(context, 0, 0, params)
	commonPath := controller.Template.CommonTemplate(context, params)

	return &gweb.HTMLResult{
		Template: []string{
			menusPath, commonPath,
		},
		Params: params,
	}
}
func (controller *Controller) index(context *gweb.Context) gweb.Result {

	//siteName := context.PathParams["siteName"]

	params := make(map[string]interface{})

	_, menusPath := controller.Template.MenusTemplate(context, 0, 0, params)
	commonPath := controller.Template.CommonTemplate(context, params)

	return &gweb.HTMLResult{
		Template: []string{
			menusPath, commonPath,
		},
		Params: params,
	}
}
func (controller *Controller) AddProjectdsfdsfsdAction(context *gweb.Context) gweb.Result {

	return &gweb.FileServerResult{}
}
func (controller *Controller) AddProjectAction(context *gweb.Context) gweb.Result {

	glog.Trace(context.Request.URL)
	//var project dao.Project

	//util.RequestBodyToJSON(context.Request.Body, &project)

	//fmt.Println(project)

	//controller.Dao.Add(service.Orm, &project)

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "信息已经提交，我们会在第一时间联系您。", Data: nil}}
}
