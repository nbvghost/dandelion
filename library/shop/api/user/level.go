package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/nbvghost/gpa/types"
)

type Level struct {
	UserService user.UserService
	Get         struct {
		UserID types.PrimaryKey `uri:"UserID"`
	} `method:"Get"`
}

func (m *Level) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	//UserID, _ := strconv.ParseUint(context.PathParams["UserID"], 10, 64)

	leve1UserIDs := m.UserService.Leve1(m.Get.UserID)

	users := m.UserService.FindUserByIDs(leve1UserIDs)

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: users}}, nil
}
