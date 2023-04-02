package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/nbvghost/gpa/types"
)

type InfoUser struct {
	UserService user.UserService
	Get         struct {
		UserID types.PrimaryKey `uri:"UserID"`
	} `method:"Get"`
}

func (m *InfoUser) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//UserID, _ := strconv.ParseUint(context.PathParams["UserID"], 10, 64)

	//var user model.User
	user := dao.GetByPrimaryKey(singleton.Orm(), entity.User, m.Get.UserID)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", user)}, nil
}
