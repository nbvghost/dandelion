package user

import (
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
)

type GrowthList struct {
	Get struct {
		Order string `uri:"Order"`
	} `method:"Get"`
}

func (m *GrowthList) Handle(context constrain.IContext) (constrain.IResult, error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)
	var Order string
	if strings.EqualFold(m.Get.Order, "asc") {
		Order = `"Growth" asc,"Score" asc,"Amount"+"BlockAmount" asc`
	} else if strings.EqualFold(m.Get.Order, "desc") {
		Order = `"Growth" desc,"Score" desc,"Amount"+"BlockAmount" desc`
	} else {
		Order = `"Growth" asc,"Score" asc,"Amount"+"BlockAmount" asc`
	}
	//var users []model.User
	//err := m.UserService.FindOrderWhereLength(db.Orm(), Order, &users, 20)
	//DB.Model(target).Order(Order).Limit(Length).Find(target).Error
	users := dao.Find(db.Orm(), &model.User{}).Order(Order).Where(`"Growth">0`).List()
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", users)}, nil

}
