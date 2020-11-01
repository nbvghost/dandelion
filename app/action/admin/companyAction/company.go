package companyAction

import (
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/content"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	Content content.ContentService
}

func (controller *Controller) Init() {

	//------------------ArticleService.go-datatables------------------------
	controller.AddHandler(gweb.GETMethod("info", controller.getInfoAction))
	controller.AddHandler(gweb.POSMethod("info", controller.changeInfoAction))

}
func (controller *Controller) changeInfoAction(context *gweb.Context) gweb.Result {

	company := &dao.Organization{}
	util.RequestBodyToJSON(context.Request.Body, company)

	return &gweb.JsonResult{}
}
func (controller *Controller) getInfoAction(context *gweb.Context) gweb.Result {

	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", company)}
}
