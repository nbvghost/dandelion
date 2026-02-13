package content_sub_type

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type GetContentSubTypeID struct {
	GET struct {
		ContentSubTypeID uint `uri:"ContentSubTypeID"`
	} `method:"GET"`
}

func (m *GetContentSubTypeID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *GetContentSubTypeID) HandleGet(ctx constrain.IContext) (r constrain.IResult, err error) {
	//ContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ContentSubTypeID"], 10, 64)
	//ContentSubTypeID := object.ParseUint(context.PathParams["ContentSubTypeID"])
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	item := repository.ContentSubTypeDao.GetContentSubTypeByID(ctx, dao.PrimaryKey(m.GET.ContentSubTypeID))

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", item)}, nil
}
