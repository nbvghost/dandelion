package content_item

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Add struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		*model.ContentItem
	} `method:"POST"`
}

func (m *Add) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Add) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	/*item := &model.ContentItem{}
	err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}*/

	return &result.JsonResult{Data: service.Content.SaveContentItem(ctx, m.Organization.ID, m.POST.ContentItem)}, nil
}
