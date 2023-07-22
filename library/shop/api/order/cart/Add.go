package cart

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/viewmodel"
	"github.com/nbvghost/dandelion/service/order"
)

type Add struct {
	OrdersService order.OrdersService
	User          *model.User                  `mapping:""`
	Post          viewmodel.GoodsSpecification `method:"post"`
}

func (m *Add) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	//context.Request.ParseForm()
	//GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	//SpecificationID, _ := strconv.ParseUint(context.Request.FormValue("SpecificationID"), 10, 64)
	//Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)

	//GoodsID := object.ParseUint(context.Request.FormValue("GoodsID"))
	//SpecificationID := object.ParseUint(context.Request.FormValue("SpecificationID"))
	//Quantity := object.ParseUint(context.Request.FormValue("Quantity"))

	err := m.OrdersService.AddCartOrders(m.User.ID, dao.PrimaryKey(m.Post.GoodsID), dao.PrimaryKey(m.Post.SpecificationID), uint(m.Post.Quantity))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "已添加到购物车", nil)}, nil
}

func (m *Add) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
