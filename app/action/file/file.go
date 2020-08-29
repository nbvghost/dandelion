package file

import (
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/gweb"
	"strconv"
)

type Controller struct {
	gweb.BaseController
}

func (controller *Controller) Init() {
	controller.AddHandler(gweb.POSMethod("up", controller.upAction))
}
func (controller *Controller) upAction(context *gweb.Context) gweb.Result {

	if context.Session.Attributes.Get(play.SessionAdmin) == nil {

		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "[ADMIN]没有登陆"}}
	}

	admin := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)
	dynamic := strconv.Itoa(int(admin.OID))
	gweb.FileUploadAction(context, dynamic)

	return &gweb.EmptyResult{}
}
