package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type InfoUser struct {
	Get struct {
		UserID dao.PrimaryKey `uri:"UserID"`
	} `method:"Get"`
}

func (m *InfoUser) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//UserID, _ := strconv.ParseUint(context.PathParams["UserID"], 10, 64)

	//var user model.User
	user := dao.GetByPrimaryKey(db.Orm(), entity.User, m.Get.UserID)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", user)}, nil
}
