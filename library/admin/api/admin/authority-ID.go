package admin

import (
	"errors"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/entity/model"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type AuthorityID struct {
	Admin *model.Admin `mapping:""`
	Put   struct {
		ID    uint         `uri:"ID"`
		Admin *model.Admin `body:""`
	} `method:"Put"`
}

func (m *AuthorityID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *AuthorityID) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {

	Orm := db.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	/*item := &model.Admin{}
	err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}*/

	//var _admin model.Admin
	_admin := dao.GetByPrimaryKey(Orm, entity.Admin, dao.PrimaryKey(m.Put.ID)).(*model.Admin)
	if strings.EqualFold(_admin.Account, "admin") {
		//admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)
		if strings.EqualFold(m.Admin.Account, _admin.Account) {

		} else {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无权修改admin账号权限"), "", nil)}, err
		}
	}

	//err = dao.UpdateByPrimaryKey(Orm, entity.Admin, dao.PrimaryKey(m.Put.ID), &model.Admin{Authority: m.Put.Admin.Authority})
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
