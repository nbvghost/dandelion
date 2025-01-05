package content_item

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type List struct {
	Organization *model.Organization `mapping:""`
}

func (m *List) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	dts := repository.ContentItemDao.ListContentItemByOID(m.Organization.ID)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", dts)}, nil
}
