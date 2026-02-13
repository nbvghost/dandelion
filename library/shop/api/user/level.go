package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Level struct {
	Get struct {
		UserID dao.PrimaryKey `uri:"UserID"`
	} `method:"Get"`
}

func (m *Level) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {

	//UserID, _ := strconv.ParseUint(context.PathParams["UserID"], 10, 64)

	leve1UserIDs := service.User.Leve1(ctx, m.Get.UserID)

	users := service.User.FindUserByIDs(ctx, leve1UserIDs)

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: users}}, nil
}
