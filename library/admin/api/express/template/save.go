package template

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Save struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		*model.ExpressTemplate
	} `method:"Post"`
}

func (m *Save) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Save) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//item := &model.ExpressTemplate{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//log.Println(err)
	m.Post.ExpressTemplate.OID = m.Organization.ID
	err = service.Express.ExpressTemplate.SaveExpressTemplate(m.Post.ExpressTemplate)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "保存成功", nil)}, err
}
