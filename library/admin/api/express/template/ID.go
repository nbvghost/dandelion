package template

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ID struct {
	Get struct {
		ID uint `uri:"ID"`
	} `method:"get"`
	Delete struct {
		ID uint `uri:"ID"`
	} `method:"Delete"`
}

func (m *ID) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.GetDB(ctx)
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//var item model.ExpressTemplate
	item := dao.GetByPrimaryKey(Orm, entity.ExpressTemplate, dao.PrimaryKey(m.Get.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", item)}, err
}

func (m *ID) HandleDelete(ctx constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.GetDB(ctx)
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//err = m.ExpressTemplate.Delete(Orm, &model.ExpressTemplate{}, dao.PrimaryKey(m.Delete.ID))
	err = dao.DeleteByPrimaryKey(Orm, entity.ExpressTemplate, dao.PrimaryKey(m.Delete.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}
