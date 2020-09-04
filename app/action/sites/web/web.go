package web

import (
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/sites"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	Template sites.TemplateService
}

func (controller *Controller) Init() {
	controller.AddHandler(gweb.ALLMethod("/index", controller.index))
	controller.AddHandler(gweb.ALLMethod("/gallery", controller.gallery))
	controller.AddHandler(gweb.ALLMethod("/js/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/css/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/img/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/font/", controller.AddProjectdsfdsfsdAction))
	//controller.AddHandler(gweb.ALLMethod("/", controller.defaultPage))
}

func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {

	//return &gweb.HTMLResult{}
	return &gweb.RedirectToUrlResult{Url: "index"}

}
func (controller *Controller) gallery(context *gweb.Context) gweb.Result {

	//siteName := context.PathParams["siteName"]

	params := make(map[string]interface{})

	menusPath := controller.Template.MenusTemplate(context, params)
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

	menusPath := controller.Template.MenusTemplate(context, params)
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
