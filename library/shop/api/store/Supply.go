package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/wechat"

	"github.com/nbvghost/tool"
)

type Supply struct {
	OrdersService order.OrdersService
	StoreService  company.StoreService
	Store         *model.Store        `mapping:""`
	User          *model.User         `mapping:""`
	WechatConfig  *model.WechatConfig `mapping:""`
	Wx            wechat.WxService
	Post          struct {
		PayMoney uint `form:"PayMoney"`
	} `method:"Post"`
}

func (m *Supply) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(context)
	//PayMoney, _ := strconv.ParseUint(context.Request.FormValue("PayMoney"), 10, 64)

	if m.Post.PayMoney <= 0 {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: "无效的金额", Data: nil}}, nil
	}
	ip := util.GetIP(contextValue.Request)

	supply := model.SupplyOrders{}
	supply.StoreID = m.Store.ID
	supply.OrderNo = tool.UUID()
	supply.PayMoney = m.Post.PayMoney
	supply.UserID = m.User.ID
	supply.Type = play.SupplyType_Store

	WxConfig := m.WechatConfig

	Success, Message, Result := m.Wx.Order(context, supply.OrderNo, "门店", "充值", "", m.User.OpenID, ip, m.Post.PayMoney, play.OrdersType_Supply, WxConfig)
	if Success != result.Success {
		return &result.JsonResult{Data: &result.ActionResult{Code: Success, Message: Message, Data: Result}}, nil
	}

	err := dao.Create(singleton.Orm(), &supply)
	if err != nil {
		return nil, err
	}

	//WxConfig := controller.Wx.MiniProgram()

	outData, err := m.Wx.GetWXAConfig(*Result.PrepayId, WxConfig)
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: outData}}, nil
}

func (m *Supply) Handle(context constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")

}
