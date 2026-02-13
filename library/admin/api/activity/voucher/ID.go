package voucher

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ID struct {
	Organization *model.Organization `mapping:""`
	Put          struct {
		ID      uint           `uri:"ID"`
		Voucher *model.Voucher `body:""`
	} `method:"Put"`
	Delete struct {
		ID uint `uri:"ID"`
	} `method:"Delete"`
	Get struct {
		ID uint `uri:"ID"`
	} `method:"Get"`
}

func (m *ID) HandlePut(ctx constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.GetDB(ctx)
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Voucher{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//if err != nil {
	//	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	//}
	err = dao.UpdateByPrimaryKey(Orm, entity.Voucher, dao.PrimaryKey(m.Put.ID), m.Put.Voucher)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}

func (m *ID) HandleDelete(ctx constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.GetDB(ctx)
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Voucher{}
	err = dao.DeleteByPrimaryKey(Orm, entity.Voucher, dao.PrimaryKey(m.Delete.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}

func (m *ID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *ID) HandleGet(ctx constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.GetDB(ctx)
	///ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Voucher{}
	item := dao.GetByPrimaryKey(Orm, entity.Voucher, dao.PrimaryKey(m.Get.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", item)}, err
}
