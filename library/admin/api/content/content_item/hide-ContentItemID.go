package content_item

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type HideContentItemID struct {
	PUT struct {
		ContentItemID uint `uri:"ContentItemID"`
		*model.ContentItem
	} `method:"PUT"`
}

func (m *HideContentItemID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *HideContentItemID) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ContentItemID"])
	//item := &model.ContentItem{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//if err != nil {
	//	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	//}
	err = dao.UpdateByPrimaryKey(Orm, entity.ContentItem, dao.PrimaryKey(m.PUT.ContentItemID), map[string]interface{}{"ShowAtMenu": m.PUT.ContentItem.ShowAtMenu})
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "index成功", nil)}, err
}
