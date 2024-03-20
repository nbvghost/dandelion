package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type ExpressInfo struct {
	Post struct {
		OrdersID dao.PrimaryKey `form:"OrdersID"`
	} `method:"post"`
}

func (m *ExpressInfo) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	//et := service.ExpressTemplateService{}
	//et.GetExpressInfo(4545458, "3957600136312", "韵达快递")
	//context.Request.ParseForm()
	//OrdersID, _ := strconv.ParseUint(context.Request.FormValue("OrdersID"), 10, 64)
	//OrdersID := object.ParseUint(context.Request.FormValue("OrdersID"))
	//LogisticCode := context.Request.FormValue("LogisticCode")
	//ShipperName := context.Request.FormValue("ShipperName")
	//LogisticCode, ShipperName
	Result, err := service.Express.ExpressTemplate.GetExpressInfo(ctx, m.Post.OrdersID)

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: Result}}, err
}

func (m *ExpressInfo) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}
