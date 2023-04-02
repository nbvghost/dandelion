package cart

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/order"
)

type List struct {
	User                *model.User `mapping:""`
	ShoppingCartService order.ShoppingCartService
}

func (m *List) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	list, _, err := m.ShoppingCartService.FindShoppingCartListDetails(m.User.ID)
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", list)}, nil

}
