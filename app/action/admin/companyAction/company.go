package companyAction

import (
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/company"
	"github.com/nbvghost/dandelion/app/service/content"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	Content      content.ContentService
	Organization company.OrganizationService
}

func (controller *Controller) Init() {

	//------------------ArticleService.go-datatables------------------------
	controller.AddHandler(gweb.GETMethod("info", controller.getInfoAction))
	controller.AddHandler(gweb.POSMethod("info", controller.changeInfoAction))

}
func (controller *Controller) changeInfoAction(context *gweb.Context) gweb.Result {
	sessionCompany := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	company := &dao.Organization{}
	util.RequestBodyToJSON(context.Request.Body, company)

	err := controller.Organization.ChangeOrganization(sessionCompany.ID, company)

	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", company)}
}
func (controller *Controller) getInfoAction(context *gweb.Context) gweb.Result {

	sessionCompany := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	company := controller.Organization.GetOrganization(sessionCompany.ID)

	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", company)}
}
