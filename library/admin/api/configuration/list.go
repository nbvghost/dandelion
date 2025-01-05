package configuration

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/library/result"
)

type List struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		ConfigurationKey []model.ConfigurationKey `body:""`
	} `method:"POST"`
}

func (m *List) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *List) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//var ks []sqltype.ConfigurationKey
	//util.RequestBodyToJSON(context.Request.Body, &ks)
	list := service.Configuration.GetConfigurations(m.Organization.ID, m.Post.ConfigurationKey...)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", list)}, err
}
