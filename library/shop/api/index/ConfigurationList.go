package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type ConfigurationList struct {
	Get []model.ConfigurationKey `method:"get"`
}

func (m *ConfigurationList) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)
	//var ks []sqltype.ConfigurationKey
	//util.RequestBodyToJSON(context.Request.Body, &ks)
	list := service.Configuration.GetConfigurations(ctx, 0, m.Get...)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", list)}, nil
}
func (m *ConfigurationList) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	return nil, nil
}
