package content_item

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ShowAtHome struct {
	PUT struct {
		*model.ContentItem
	} `method:"PUT"`
}

func (m *ShowAtHome) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *ShowAtHome) HandlePut(ctx constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.GetDB(ctx)
	err = dao.UpdateByPrimaryKey(Orm, entity.ContentItem, dao.PrimaryKey(m.PUT.ID), map[string]interface{}{"ShowAtHome": m.PUT.ContentItem.ShowAtHome})
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
