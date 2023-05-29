package user

import (
	"github.com/nbvghost/dandelion/library/db"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/user"
)

type GrowthList struct {
	UserService user.UserService
	Get         struct {
		Order string `uri:"Order"`
	} `method:"Get"`
}

func (m *GrowthList) Handle(context constrain.IContext) (constrain.IResult, error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)
	var Order string
	if strings.EqualFold(m.Get.Order, "asc") {
		Order = `"Growth" asc`
	} else if strings.EqualFold(Order, "desc") {
		Order = `"Growth" desc`
	} else {
		Order = `"Growth" asc`
	}
	var users []model.User
	err := m.UserService.FindOrderWhereLength(db.Orm(), Order, &users, 20)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", users)}, err

}
