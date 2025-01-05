package content_sub_type

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type AllListContentItemID struct {
	GET struct {
		ContentItemID uint `uri:"ContentItemID"`
	} `method:"GET"`
}

func (m *AllListContentItemID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//ContentItemID := object.ParseUint(context.PathParams["ContentItemID"])
	list := repository.ContentSubTypeDao.FindContentSubTypesByContentItemID(m.GET.ContentItemID)

	resultMap := make(map[dao.PrimaryKey]interface{})

	for index := range list {
		item := list[index]
		subTypes := repository.ContentSubTypeDao.FindContentSubTypesByContentItemIDAndParentContentSubTypeID(dao.PrimaryKey(m.GET.ContentItemID), item.ID)

		childrenMap := make(map[dao.PrimaryKey]interface{})

		for sindex := range subTypes {

			childrenMap[subTypes[sindex].ID] = subTypes[sindex]

		}

		resultMap[item.ID] = map[string]interface{}{
			"SubType":         item,
			"SubTypeChildren": childrenMap,
		}

	}

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", resultMap)}, nil
}
