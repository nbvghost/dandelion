package sites

import (
	"github.com/nbvghost/dandelion/app/action/sites/shop"
	"github.com/nbvghost/dandelion/app/action/sites/web"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service/company"
	"github.com/nbvghost/dandelion/app/service/dao"
	"strings"

	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	Dao dao.BaseDao
}
type InterceptorManager struct {
	Organization company.OrganizationService
}

func (m InterceptorManager) Execute(context *gweb.Context) (bool, gweb.Result) {

	siteName := context.PathParams["siteName"]

	if strings.EqualFold(siteName, "") {
		return false, &gweb.HTMLResult{Name: "404"}
	}
	//util.Trace(context.Session,"context.Session")
	if context.Session.Attributes.Get(play.SessionOrganization) == nil {
		org := m.Organization.FindByDomain(dao.Orm(), siteName)
		if org == nil {
			return false, &gweb.HTMLResult{Name: "404"}
		}
		context.Session.Attributes.Put(play.SessionOrganization, org)
		return true, nil
	} else {
		org := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
		if strings.EqualFold(org.Domain, siteName) {
			return true, nil
		} else {
			org := m.Organization.FindByDomain(dao.Orm(), siteName)
			if org == nil {
				return false, &gweb.HTMLResult{Name: "404"}
			}
			context.Session.Attributes.Put(play.SessionOrganization, org)
			return true, nil
		}
	}
}
func (controller *Controller) Init() {
	controller.Interceptors.Add(&InterceptorManager{})

	shop := &shop.Controller{}
	shop.Interceptors = controller.Interceptors
	controller.AddSubController("/{siteName}/shop/", shop)

	web := &web.Controller{}
	web.Interceptors = controller.Interceptors
	controller.AddSubController("/{siteName}/", web)

	//controller.AddHandler(gweb.ALLMethod("/", controller.defaultPage))
}

func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {

	//return &gweb.HTMLResult{}
	return &gweb.HTMLResult{}

}
