package score_goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"

	"github.com/nbvghost/dandelion/library/result"
)

type ID struct {
	Put struct {
		ID         uint              `uri:"ID"`
		ScoreGoods *model.ScoreGoods `body:""`
	} `method:"Put"`
	Delete struct {
		ID uint `uri:"ID"`
	} `method:"Delete"`
	Get struct {
		ID uint `uri:"ID"`
	} `method:"Get"`
}

func (m *ID) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.Orm()
	//ID := object.ParseUint(context.PathParams["ID"])
	/*item := &model.ScoreGoods{}
	err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}*/
	err = dao.UpdateByPrimaryKey(Orm, entity.ScoreGoods, dao.PrimaryKey(m.Put.ID), m.Put.ScoreGoods)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}

func (m *ID) HandleDelete(context constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.Orm()
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.ScoreGoods{}
	err = dao.DeleteByPrimaryKey(Orm, entity.ScoreGoods, dao.PrimaryKey(m.Delete.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}

func (m *ID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *ID) HandleGet(context constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.Orm()
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.ScoreGoods{}
	item := dao.GetByPrimaryKey(Orm, entity.ScoreGoods, dao.PrimaryKey(m.Get.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", item)}, err
}
