package order

import (
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"log"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
)

type WXPayPackage struct {
	User *model.User `mapping:""`
	//WechatConfig *model.WechatConfig `mapping:""`
	Get struct {
		OrderNo string `form:"OrderNo"`
	} `method:"get"`
}

func (m *WXPayPackage) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(ctx)
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//OrderNo := context.Request.URL.Query().Get("OrderNo")
	//OrderType := context.Request.URL.Query().Get("OrderType")

	ip := util.GetIP(contextValue.Request)

	wechat := service.Payment.NewWechat(ctx, m.User.OID)

	//package
	orders := service.Order.Orders.GetOrdersPackageByOrderNo(m.Get.OrderNo)
	if strings.EqualFold(orders.PrepayID, "") == false {

		outData, err := wechat.GetWXAConfig(orders.PrepayID)
		if err != nil {
			return nil, err
		}
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: outData}}, nil

	}

	r, err := wechat.MPOrder(orders.OrderNo, "购物", "商品消费", []model.OrdersGoods{}, m.User.OpenID, ip, orders.TotalPayMoney, model.OrdersTypeGoodsPackage)
	if err != nil {
		return result.NewError(err), nil //&result.JsonResult{Data: &result.ActionResult{Code: Success, Message: Message, Data: Result}}, nil
	}

	outData, err := wechat.GetWXAConfig(r.PrepayId)
	if err != nil {
		return nil, err
	}

	err = dao.UpdateByPrimaryKey(db.Orm(), entity.OrdersPackage, orders.ID, map[string]interface{}{"PrepayID": r.PrepayId})
	if err != nil {
		log.Println(err)
	}
	//outData["OrdersID"] = strconv.Itoa(int(orders.ID))
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: outData}}, nil

}
