package sites

import (
	"github.com/nbvghost/dandelion/app/action/sites/shop"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/glog"

	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	Dao dao.BaseDao
}

func (controller *Controller) Init() {
	controller.AddHandler(gweb.ALLMethod("index", controller.AddProjectAction))
	controller.AddHandler(gweb.ALLMethod("/{siteName}/js/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/{siteName}/css/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/{siteName}/img/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/{siteName}/font/", controller.AddProjectdsfdsfsdAction))

	shop := &shop.Controller{}
	shop.Interceptors = controller.Interceptors
	controller.AddSubController("/{siteName}/shop/", shop)

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
