package order

import (
	"log"
	"strings"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
)

type WXPayAlone struct {
	User *model.User `mapping:""`
	//WechatConfig *model.WechatConfig `mapping:""`
	Get struct {
		OrderNo string `form:"OrderNo"`
	} `method:"get"`
}

func (m *WXPayAlone) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//OrderNo := context.Request.URL.Query().Get("OrderNo")
	//OrderType := context.Request.URL.Query().Get("OrderType")
	contextValue := contexext.FromContext(ctx)
	//WxConfig := m.WechatConfig
	ip := util.GetIP(contextValue.Request)

	wechat := service.Payment.NewWechat(ctx, m.User.OID)

	//package
	orders := repository.OrdersDao.GetOrdersByOrderNo(ctx, m.Get.OrderNo)
	if strings.EqualFold(orders.PrepayID, "") == false {

		outData, err := wechat.GetWXAConfig(orders.PrepayID)
		if err != nil {
			return nil, err
		}
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: outData}}, nil
	}

	r, err := wechat.MPOrder(orders.OrderNo, "购物", "商品消费", []model.OrdersGoods{}, m.User.OpenID, ip, orders.PayMoney, model.OrdersTypeGoods)
	if err != nil {
		return result.NewError(err), nil //&result.JsonResult{Data: &result.ActionResult{Code: Success, Message: err.Error(), Data: Result}}, nil
	}

	outData, err := wechat.GetWXAConfig(r.PrepayId)
	if err != nil {
		return nil, err
	}

	err = dao.UpdateByPrimaryKey(db.GetDB(ctx), entity.Orders, orders.ID, map[string]interface{}{"PrepayID": r.PrepayId})
	if err != nil {
		log.Println(err)
	}
	//outData["OrdersID"] = strconv.Itoa(int(orders.ID))
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: outData}}, nil

}
