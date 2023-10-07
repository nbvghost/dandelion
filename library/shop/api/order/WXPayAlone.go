package order

import (
	"github.com/nbvghost/dandelion/library/db"
	"log"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/wechat"
)

type WXPayAlone struct {
	WxService     wechat.WxService
	OrdersService order.OrdersService
	User          *model.User         `mapping:""`
	WechatConfig  *model.WechatConfig `mapping:""`
	Get           struct {
		OrderNo string `form:"OrderNo"`
	} `method:"get"`
}

func (m *WXPayAlone) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//OrderNo := context.Request.URL.Query().Get("OrderNo")
	//OrderType := context.Request.URL.Query().Get("OrderType")
	contextValue := contexext.FromContext(ctx)
	WxConfig := m.WechatConfig
	ip := util.GetIP(contextValue.Request)

	//package
	orders := m.OrdersService.GetOrdersByOrderNo(m.Get.OrderNo)
	if strings.EqualFold(orders.PrepayID, "") == false {

		outData, err := m.WxService.GetWXAConfig(orders.PrepayID, WxConfig)
		if err != nil {
			return nil, err
		}
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: outData}}, nil
	}

	Success, Message, Result := m.WxService.MPOrder(ctx, orders.OrderNo, "购物", "商品消费", []model.OrdersGoods{}, m.User.OpenID, ip, orders.PayMoney, play.OrdersTypeGoods, WxConfig)
	if Success != result.Success {
		return &result.JsonResult{Data: &result.ActionResult{Code: Success, Message: Message, Data: Result}}, nil
	}

	outData, err := m.WxService.GetWXAConfig(*Result.PrepayId, WxConfig)
	if err != nil {
		return nil, err
	}

	err = dao.UpdateByPrimaryKey(db.Orm(), entity.Orders, orders.ID, map[string]interface{}{"PrepayID": Result.PrepayId})
	if err != nil {
		log.Println(err)
	}
	//outData["OrdersID"] = strconv.Itoa(int(orders.ID))
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: outData}}, nil

}
