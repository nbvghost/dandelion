package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/pkg/errors"
)

type BuyCollage struct {
	OrdersService order.OrdersService
	User          *model.User `mapping:""`
	Post          struct {
		GoodsID         dao.PrimaryKey `form:"GoodsID"`
		SpecificationID dao.PrimaryKey `form:"SpecificationID"`
		Quantity        uint           `form:"Quantity"`
	} `method:"post"`
}

func (m *BuyCollage) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//context.Request.ParseForm()
	//No := context.Request.FormValue("No")
	//GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	//SpecificationID, _ := strconv.ParseUint(context.Request.FormValue("SpecificationID"), 10, 64)
	//Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)

	//GoodsID := object.ParseUint(context.Request.FormValue("GoodsID"))
	//SpecificationID := object.ParseUint(context.Request.FormValue("SpecificationID"))
	//Quantity := object.ParseUint(context.Request.FormValue("Quantity"))

	if m.Post.GoodsID != 0 && m.Post.SpecificationID != 0 && m.Post.Quantity != 0 {
		err := m.OrdersService.BuyCollageOrders(ctx, m.User.ID, dao.PrimaryKey(m.Post.GoodsID), dao.PrimaryKey(m.Post.SpecificationID), uint(m.Post.Quantity))
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "立即购买", nil)}, nil
	} else {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("订单数据出错"), "", nil)}, nil
	}
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("订单数据出错"), "", nil)}, nil
}

func (m *BuyCollage) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
