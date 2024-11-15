package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type ListContentSubType struct {
	Get struct {
		ContentItemID dao.PrimaryKey `form:"ContentItemID"`
		PID           dao.PrimaryKey `form:"PID"`
	} `method:"Get"`
}

func (g *ListContentSubType) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ListContentSubType) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)

	//ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//ContentItemID := object.ParseUint(context.PathParams["ContentItemID"])
	//PID, _ := strconv.ParseUint(context.Request.URL.Query().Get("PID"), 10, 64)
	//PID := object.ParseUint(context.Request.URL.Query().Get("PID"))

	content := repository.ContentItemDao.GetContentItemByID(dao.PrimaryKey(g.Get.ContentItemID))

	csts := repository.ContentSubTypeDao.FindContentSubTypesByContentItemIDAndParentContentSubTypeID(content.ID, dao.PrimaryKey(g.Get.PID))

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: csts}}, nil

}
