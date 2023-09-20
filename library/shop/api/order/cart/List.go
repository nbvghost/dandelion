package cart

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/order"
)

type List struct {
	User          *model.User `mapping:""`
	OrdersService order.OrdersService
	Get           struct {
		AddressID dao.PrimaryKey `form:"address-id"`
	} `method:"get"`
}

func (m *List) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	var address = &model.Address{}
	if m.Get.AddressID > 0 {
		address = dao.GetByPrimaryKey(db.Orm(), &model.Address{}, m.Get.AddressID).(*model.Address)
	} else {
		addressList := dao.Find(db.Orm(), &model.Address{}).Where(`"UserID"=? and "DefaultShipping"=true`, ctx.UID()).List()
		if len(addressList) > 0 {
			address = addressList[0].(*model.Address)
		}
	}

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	list, err := m.OrdersService.FindShoppingCartListDetails(m.User.OID, m.User.ID, address)
	if err != nil {
		return nil, err
	}
	return result.NewData(list), nil //&result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", list)}, nil

}
