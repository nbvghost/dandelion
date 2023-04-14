package cart

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/gpa/types"
)

type List struct {
	User                *model.User `mapping:""`
	ShoppingCartService order.ShoppingCartService
	Get                 struct {
		AddressID types.PrimaryKey `form:"address-id"`
	} `method:"get"`
}

func (m *List) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	var address *model.Address
	if m.Get.AddressID > 0 {
		address = dao.GetByPrimaryKey(singleton.Orm(), &model.Address{}, m.Get.AddressID).(*model.Address)
	} else {
		addressList := dao.Find(singleton.Orm(), &model.Address{}).Where(`"UserID"=? and "DefaultShipping"=true`, ctx.UID()).List()
		if len(addressList) > 0 {
			address = addressList[0].(*model.Address)
		}
	}

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	list, _, err := m.ShoppingCartService.FindShoppingCartListDetails(m.User.ID, address)
	if err != nil {
		return nil, err
	}
	return result.NewData(list), nil //&result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", list)}, nil

}
