package cart

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/viewmodel"
	"github.com/nbvghost/dandelion/service"
)

type Add struct {
	User *model.User                  `mapping:""`
	Post viewmodel.GoodsSpecification `method:"post"`
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

	err := service.Order.Orders.AddCartOrders(ctx, m.User.ID, m.Post.GoodsID, m.Post.SpecificationID, m.Post.Quantity)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "已添加到购物车", nil)}, nil
}

func (m *Add) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
