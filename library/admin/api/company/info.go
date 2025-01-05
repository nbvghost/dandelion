package company

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/library/result"
)

type Info struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		*model.Organization
	} `method:"Post"`
}

func (m *Info) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//sessionCompany := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//company := &model.Organization{}
	//util.RequestBodyToJSON(context.Request.Body, company)

	err = service.Company.Organization.ChangeOrganization(m.Organization.ID, m.Post.Organization)

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", m.Post.Organization)}, err

}

func (m *Info) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//sessionCompany := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	company := service.Company.Organization.GetOrganization(m.Organization.ID)

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", company)}, nil
}
