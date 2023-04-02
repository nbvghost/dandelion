package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/express"
	"github.com/nbvghost/gpa/types"
)

type ExpressInfo struct {
	ExpressTemplateService express.ExpressTemplateService
	Post                   struct {
		OrdersID     types.PrimaryKey `form:"OrdersID"`
		LogisticCode string           `form:"LogisticCode"`
		ShipperName  string           `form:"ShipperName"`
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
	Result := m.ExpressTemplateService.GetExpressInfo(m.Post.OrdersID, m.Post.LogisticCode, m.Post.ShipperName)

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: Result}}, nil
}

func (m *ExpressInfo) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
