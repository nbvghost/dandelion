package order

import (
	"encoding/json"
	"log"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/server/redis"
	"github.com/nbvghost/dandelion/service/order"
)

type ConfirmList struct {
	OrdersService order.OrdersService
	User          *model.User `mapping:""`
	Post          struct {
		PostType int    `form:"PostType"`
		Address  string `form:"Address"`
	} `method:"post"`
}

func (m *ConfirmList) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	var ogs []model.OrdersGoods

	confirmOrdersJson, err := ctx.Redis().Get(ctx, redis.NewConfirmOrders(ctx.UID()))
	if err == nil {
		err = json.Unmarshal([]byte(confirmOrdersJson), &ogs)
		if err != nil {
			return nil, err
		}
	} else {
		log.Println(err)
	}

	/*if context.Session.Attributes.Get(play.SessionConfirmOrders) == nil {
		ogs = make([]entity.OrdersGoods, 0)
	} else {
		ogs = *(context.Session.Attributes.Get(play.SessionConfirmOrders)).(*[]entity.OrdersGoods)
	}*/
	//context.Request.ParseForm()

	//PostType, _ := strconv.ParseInt(context.Request.FormValue("PostType"), 10, 64)
	//AddressTxt := context.Request.FormValue("Address")
	/*

	 */
	address := model.Address{}
	util.JSONToStruct(m.Post.Address, &address)

	results, _, err := m.OrdersService.AnalyseOrdersGoodsList(m.User.ID, address, int(m.Post.PostType), ogs)

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", results)}, err
}

func (m *ConfirmList) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
