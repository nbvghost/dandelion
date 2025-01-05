package content_sub_type

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type ChildListContentItemIDParentContentSubTypeID struct {
	GET struct {
		ParentContentSubTypeID uint `uri:"ParentContentSubTypeID"`
		ContentItemID          uint `uri:"ContentItemID"`
	} `method:"GET"`
}

func (m *ChildListContentItemIDParentContentSubTypeID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//ParentContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ParentContentSubTypeID"], 10, 64)
	//ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//ParentContentSubTypeID := object.ParseUint(context.PathParams["ParentContentSubTypeID"])
	//ContentItemID := object.ParseUint(context.PathParams["ContentItemID"])
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	list := repository.ContentSubTypeDao.FindContentSubTypesByContentItemIDAndParentContentSubTypeID(dao.PrimaryKey(m.GET.ContentItemID), dao.PrimaryKey(m.GET.ParentContentSubTypeID))

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", list)}, nil
}
