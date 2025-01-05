package store

import (
	"errors"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/entity/model"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ID struct {
	PUT struct {
		ID    uint         `uri:"ID"`
		Store *model.Store `body:""`
	} `method:"PUT"`
	Delete struct {
		ID uint `uri:"ID"`
	} `method:"Delete"`
	Get struct {
		ID uint `uri:"ID"`
	} `method:"Get"`
}

func (m *ID) HandlePut(context constrain.IContext) (constrain.IResult, error) {
	Orm := db.Orm()
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Store{}
	//err := util.RequestBodyToJSON(context.Request.Body, item)
	/*if err != nil {
		return nil, err
	}*/

	var _store model.Store
	service.Company.Store.GetByPhone(m.PUT.Store.Phone)
	if _store.ID > 0 && _store.ID != m.PUT.Store.ID {
		return nil, errors.New("手机号：" + m.PUT.Store.Phone + "已经被使用")
	}

	err := dao.UpdateByPrimaryKey(Orm, entity.Store, dao.PrimaryKey(m.PUT.ID), m.PUT.Store)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}

func (m *ID) HandleDelete(context constrain.IContext) (constrain.IResult, error) {
	Orm := db.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Store{}
	err := dao.DeleteByPrimaryKey(Orm, entity.Store, dao.PrimaryKey(m.Delete.ID))
	return nil, err
}

func (m *ID) Handle(context constrain.IContext) (constrain.IResult, error) {
	panic("implement me")
}

func (m *ID) HandleGet(context constrain.IContext) (constrain.IResult, error) {
	Orm := db.Orm()
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Store{}
	item := dao.GetByPrimaryKey(Orm, entity.Store, dao.PrimaryKey(m.Get.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", item)}, nil
}
