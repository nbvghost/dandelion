package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/entity/model"

	"github.com/nbvghost/dandelion/library/result"
)

type SubType struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		*model.ContentSubType
	} `method:"POST"`
}

func (m *SubType) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *SubType) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//Orm := db.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	/*item := &model.ContentSubType{}
	err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}*/
	err = service.Content.SaveContentSubType(m.Organization.ID, m.POST.ContentSubType)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}, err
}
