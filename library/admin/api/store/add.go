package store

import (
	"errors"

	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/entity/model"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Add struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		Store *model.Store `body:""`
	} `method:"POST"`
}

func (controller *Add) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Add) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	Orm := db.GetDB(ctx)
	item := &model.Store{}
	m.POST.Store.OID = m.Organization.ID
	/*err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}*/
	var _store model.Store
	_store = service.Company.Store.GetByPhone(ctx, item.Phone)
	if _store.ID > 0 {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("手机号："+item.Phone+"已经被使用"), "", nil)}, nil
	}

	err = dao.Create(Orm, item)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}, err
}
