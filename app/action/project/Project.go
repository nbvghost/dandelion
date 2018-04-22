package project

import (
	"fmt"

	"dandelion/app/service"
	"dandelion/app/service/dao"
	"dandelion/app/util"

	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	Dao dao.BaseDao
}

func (controller *Controller) Apply() {
	controller.AddHandler(gweb.ALLMethod("add", controller.AddProjectAction))

}
func (controller *Controller) AddProjectAction(context *gweb.Context) gweb.Result {

	var project dao.Project

	util.RequestBodyToJSON(context.Request.Body, &project)

	fmt.Println(project)

	controller.Dao.Add(service.Orm, &project)

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "信息已经提交，我们会在第一时间联系您。", Data: nil}}
}
