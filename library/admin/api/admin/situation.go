package admin

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Situation struct {
	Admin *model.Admin `mapping:""`
	Post  struct {
		StartTime int64 `form:"StartTime"`
		EndTime   int64 `form:"EndTime"`
	} `method:"Post"`
}

func (m *Situation) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Situation) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	//context.Request.ParseForm()
	//StartTime, _ := strconv.ParseInt(context.Request.FormValue("StartTime"), 10, 64)
	//EndTime, _ := strconv.ParseInt(context.Request.FormValue("EndTime"), 10, 64)

	results := make(map[string]interface{})

	results["Orders"] = service.Order.Orders.Situation(ctx, m.Post.StartTime, m.Post.EndTime)
	results["Voucher"] = service.Activity.Voucher.Situation(ctx, m.Post.StartTime, m.Post.EndTime)
	results["ScoreGoods"] = service.Activity.ScoreGoods.Situation(ctx, m.Post.StartTime, m.Post.EndTime)
	results["User"] = service.User.Situation(ctx, m.Post.StartTime, m.Post.EndTime)

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", results)}, nil
}
