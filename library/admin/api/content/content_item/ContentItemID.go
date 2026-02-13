package content_item

import (
	"errors"

	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ContentItemID struct {
	Organization *model.Organization `mapping:""`
	Delete       struct {
		ContentItemID uint `uri:"ContentItemID"`
	} `method:"Delete"`
	Put struct {
		ContentItemID uint `uri:"ContentItemID"`
		*model.ContentItem
	} `method:"Put"`
	Get struct {
		ContentItemID uint `uri:"ContentItemID"`
	} `method:"Get"`
}

func (m *ContentItemID) HandleDelete(ctx constrain.IContext) (constrain.IResult, error) {
	Orm := db.GetDB(ctx)
	//ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//ContentItemID := object.ParseUint(context.PathParams["ContentItemID"])

	css := repository.ContentSubTypeDao.FindContentSubTypesByContentItemID(ctx, m.Delete.ContentItemID)
	if len(css) > 0 {

		return nil, errors.New("包含子项内容无法删除")
	}

	err := dao.DeleteByPrimaryKey(Orm, entity.ContentItem, dao.PrimaryKey(m.Delete.ContentItemID))
	if err == nil {
		err = dao.DeleteBy(Orm, entity.Content, map[string]interface{}{
			"ContentItemID":    m.Delete.ContentItemID,
			"ContentSubTypeID": 0,
		})
	}
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}

func (m *ContentItemID) Handle(context constrain.IContext) (constrain.IResult, error) {
	panic("implement me")
}

func (m *ContentItemID) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	//Orm := db.GetDB(ctx)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//ID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ContentItemID"])
	/*item := &model.ContentItem{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}*/

	err := service.Content.SaveContentItem(ctx, m.Organization.ID, m.Put.ContentItem)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
func (m *ContentItemID) HandleGet(ctx constrain.IContext) (constrain.IResult, error) {
	Orm := db.GetDB(ctx)
	//ID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ContentItemID"])
	//item := &model.ContentItem{}
	item := dao.GetByPrimaryKey(Orm, entity.ContentItem, dao.PrimaryKey(m.Get.ContentItemID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", item)}, nil
}
