package content_sub_type

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type ListContentItemID struct {
	GET struct {
		ContentItemID uint `uri:"ContentItemID"`
	} `method:"GET"`
}

func (m *ListContentItemID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//ContentItemID := object.ParseUint(context.PathParams["ContentItemID"])

	list := repository.ContentSubTypeDao.FindContentSubTypesByContentItemID(m.GET.ContentItemID)

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", list)}, nil
}
