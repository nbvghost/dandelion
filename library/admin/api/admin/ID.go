package admin

import (
	"errors"
	"github.com/nbvghost/dandelion/service"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/entity/model"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ID struct {
	Admin *model.Admin `mapping:""`
	Put   struct {
		ID           uint `uri:"ID"`
		*model.Admin `body:""`
	} `method:"Put"`
	Delete struct {
		ID uint `uri:"ID"`
	} `method:"Delete"`
	Get struct {
		ID uint `uri:"ID"`
	} `method:"Get"`
	Post struct {
		*model.Datatables `body:""`
	} `method:"Post"`
}

func (m *ID) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Admin{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//if err != nil {
	//	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	//}

	//var _admin model.Admin
	_admin := dao.GetByPrimaryKey(Orm, entity.Admin, dao.PrimaryKey(m.Put.ID)).(*model.Admin)
	if strings.EqualFold(_admin.Account, "admin") {
		//admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)
		if strings.EqualFold(m.Admin.Account, _admin.Account) {

		} else {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无权修改admin账号密码"), "", nil)}, err
		}
	}

	//m.Put.Admin.PassWord = encryption.Md5ByString(m.Put.Admin.PassWord)

	err = dao.UpdateByPrimaryKey(Orm, entity.Admin, dao.PrimaryKey(m.Put.ID), m.Put.Admin)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}

func (m *ID) HandleDelete(context constrain.IContext) (r constrain.IResult, err error) {
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Admin{}
	Orm := db.Orm()

	item := dao.GetByPrimaryKey(Orm, entity.Admin, dao.PrimaryKey(m.Delete.ID)).(*model.Admin)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}
	if strings.EqualFold(item.Account, "admin") {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("admin不能删除"), "", nil)}, err
	}

	err = dao.DeleteByPrimaryKey(Orm, entity.Admin, dao.PrimaryKey(m.Delete.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}

func (m *ID) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.Admin.Service.DatatablesListOrder(db.Orm(), m.Post.Datatables, &[]model.Admin{}, m.Admin.OID, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, err
}

func (m *ID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *ID) HandleGet(context constrain.IContext) (r constrain.IResult, err error) {
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Admin{}
	item := dao.GetByPrimaryKey(db.Orm(), entity.Admin, dao.PrimaryKey(m.Get.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", item)}, err
}
