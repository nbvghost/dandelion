package order

import (
	"dandelion/app/play"
	"dandelion/app/service/dao"
	"strconv"

	"errors"
	"strings"

	"dandelion/app/util"

	"dandelion/app/service"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type OrderController struct {
	gweb.BaseController
	ShoppingCart    service.ShoppingCartService
	Orders          service.OrdersService
	OrdersGoods     service.OrdersGoodsService
	Wx              service.WxService
	ExpressTemplate service.ExpressTemplateService
	Organization    service.OrganizationService
}

func (controller *OrderController) Apply() {

	controller.AddHandler(gweb.POSMethod("/add", controller.ordersAddAction))
	controller.AddHandler(gweb.POSMethod("/buy", controller.ordersBuyAction))
	controller.AddHandler(gweb.GETMethod("/cart/list", controller.ordersCartListAction))
	controller.AddHandler(gweb.POSMethod("/cart/delete", controller.ordersCartDeleteAction))
	controller.AddHandler(gweb.POSMethod("/cart/change", controller.ordersCartChangeAction))
	controller.AddHandler(gweb.POSMethod("/confirm/list", controller.ordersConfirmListAction))
	controller.AddHandler(gweb.POSMethod("/createOrders", controller.createOrdersAction))
	controller.AddHandler(gweb.GETMethod("/wxpay/package", controller.ordersWxpayPackageAction))
	controller.AddHandler(gweb.GETMethod("/wxpay/alone", controller.ordersWxpayAloneAction))
	controller.AddHandler(gweb.GETMethod("/list", controller.ordersListAction))
	controller.AddHandler(gweb.GETMethod("/:ID/get", controller.ordersGetListAction))
	controller.AddHandler(gweb.PUTMethod("/change", controller.orderChangeAction))
	controller.AddHandler(gweb.PUTMethod("/express/info", controller.expressInfoAction))
}
func (controller *OrderController) ordersWxpayPackageAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	OrderNo := context.Request.URL.Query().Get("OrderNo")
	//OrderType := context.Request.URL.Query().Get("OrderType")

	WxConfig := controller.Wx.MiniProgram()
	ip := util.GetIP(context)

	//package
	orders := controller.Orders.GetOrdersPackageByOrderNo(OrderNo)
	if strings.EqualFold(orders.PrepayID, "") == false {

		outData := controller.Wx.GetWXAConfig(orders.PrepayID, WxConfig)
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: outData}}

	}

	Success, Message, result := controller.Wx.MPOrder(orders.OrderNo, "购物", "商品消费", []dao.OrdersGoods{}, user.OpenID, ip, orders.TotalPayMoney, play.OrdersType_GoodsPackage)
	if Success == false {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: Success, Message: Message, Data: result}}
	}

	outData := controller.Wx.GetWXAConfig(result.Prepay_id, WxConfig)

	err := controller.Orders.ChangeMap(dao.Orm(), orders.ID, &dao.OrdersPackage{}, map[string]interface{}{"PrepayID": result.Prepay_id})
	tool.CheckError(err)
	//outData["OrdersID"] = strconv.Itoa(int(orders.ID))
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: outData}}

}
func (controller *OrderController) ordersWxpayAloneAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	OrderNo := context.Request.URL.Query().Get("OrderNo")
	//OrderType := context.Request.URL.Query().Get("OrderType")

	WxConfig := controller.Wx.MiniProgram()
	ip := util.GetIP(context)

	//package
	orders := controller.Orders.GetOrdersByOrderNo(OrderNo)
	if strings.EqualFold(orders.PrepayID, "") == false {

		outData := controller.Wx.GetWXAConfig(orders.PrepayID, WxConfig)
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: outData}}

	}

	Success, Message, result := controller.Wx.MPOrder(orders.OrderNo, "购物", "商品消费", []dao.OrdersGoods{}, user.OpenID, ip, orders.PayMoney, play.OrdersType_Goods)
	if Success == false {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: Success, Message: Message, Data: result}}
	}

	outData := controller.Wx.GetWXAConfig(result.Prepay_id, WxConfig)

	err := controller.Orders.ChangeMap(dao.Orm(), orders.ID, &dao.Orders{}, map[string]interface{}{"PrepayID": result.Prepay_id})
	tool.CheckError(err)
	//outData["OrdersID"] = strconv.Itoa(int(orders.ID))
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: outData}}

}
func (controller *OrderController) expressInfoAction(context *gweb.Context) gweb.Result {
	//et := service.ExpressTemplateService{}
	//et.GetExpressInfo(4545458, "3957600136312", "韵达快递")
	context.Request.ParseForm()
	OrdersID, _ := strconv.ParseUint(context.Request.FormValue("OrdersID"), 10, 64)
	LogisticCode := context.Request.FormValue("LogisticCode")
	ShipperName := context.Request.FormValue("ShipperName")
	//LogisticCode, ShipperName
	result := controller.ExpressTemplate.GetExpressInfo(OrdersID, LogisticCode, ShipperName)

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: result}}
}

func (controller *OrderController) orderChangeAction(context *gweb.Context) gweb.Result {
	context.Request.ParseForm()
	Action := context.Request.FormValue("Action")
	switch Action {
	case "RefundInfo":
		OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		ShipName := context.Request.FormValue("ShipName")
		ShipNo := context.Request.FormValue("ShipNo")
		err, info := controller.Orders.RefundInfo(OrdersGoodsID, ShipName, ShipNo)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, info, nil)}
	case "AskRefund":
		OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		RefundInfoJson := context.Request.FormValue("RefundInfo")
		var RefundInfo dao.RefundInfo
		util.JSONToStruct(RefundInfoJson, &RefundInfo)
		err, info := controller.Orders.AskRefund(OrdersGoodsID, RefundInfo)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, info, nil)}
	case "TakeDeliver":
		ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		err, info := controller.Orders.TakeDeliver(ID)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, info, nil)}
	case "Cancel":
		ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		err, info := controller.Orders.Cancel(ID)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, info, nil)}

	}

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("无法操作"), "OK", nil)}

}
func (controller *OrderController) ordersListAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Status := context.Request.URL.Query().Get("Status")

	var StatusList []string
	if !strings.EqualFold(Status, "") {
		StatusList = strings.Split(Status, ",")
	}

	list, _ := controller.Orders.ListOrders(user.ID, 0, 0, StatusList, 0, 0)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: list}}
	//fullcuts := controller.FullCut.FindOrderByAmountASC(service.Orm)
	//return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "", fullcuts)}
}
func (controller *OrderController) ordersGetListAction(context *gweb.Context) gweb.Result {

	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)

	pack := struct {
		Orders          dao.Orders
		OrdersGoodsList []dao.OrdersGoods
	}{}
	pack.Orders = controller.Orders.GetOrdersByID(ID)
	pack.OrdersGoodsList, _ = controller.OrdersGoods.FindByOrdersID(dao.Orm(), pack.Orders.ID)

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: pack}}
}
func (controller *OrderController) createOrdersAction(context *gweb.Context) gweb.Result {

	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	ogs := context.Session.Attributes.Get(play.SessionConfirmOrders).(*[]dao.OrdersGoods)
	context.Request.ParseForm()

	_TotalPrice, _ := strconv.ParseUint(context.Request.FormValue("TotalPrice"), 10, 64)
	//_FullCutPrice, _ := strconv.ParseUint(context.Request.FormValue("FullCutPrice"), 10, 64)
	//_GoodsPrice, _ := strconv.ParseUint(context.Request.FormValue("GoodsPrice"), 10, 64)
	//_ExpressPrice, _ := strconv.ParseUint(context.Request.FormValue("ExpressPrice"), 10, 64)
	PostType, _ := strconv.ParseInt(context.Request.FormValue("PostType"), 10, 64)
	AddressTxt := context.Request.FormValue("Address")
	address := dao.Address{}
	util.JSONToStruct(AddressTxt, &address)

	Error, results, TotalPrice := controller.Orders.AnalyseOrdersGoodsList(user.ID, address, int(PostType), *ogs)

	if _TotalPrice == TotalPrice && Error == nil {
		//controller.Orders.AddOrdersPackage(tool.UUID(),)

		if len(results) > 1 {
			orderList := make([]dao.Orders, 0)

			err, op := controller.Orders.AddOrdersPackage(TotalPrice, user.ID)
			if err != nil {
				return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", nil)}
			}
			for _, value := range results {

				oggs := value["OrdersGoodsInfos"].([]dao.OrdersGoodsInfo)
				//result["OrdersGoodsInfos"]=oggs
				FavouredPrice := value["FavouredPrice"].(uint64)
				FullCutAll := value["FullCutAll"].(uint64)
				GoodsPrice := value["GoodsPrice"].(uint64)
				ExpressPrice := value["ExpressPrice"].(uint64)

				organization := value["Organization"].(dao.Organization)

				PayMoney := GoodsPrice - FullCutAll + ExpressPrice //支付价格已经包含了 满减，限时抢购的扣去的部分  - _FullCutPrice-FavouredPrice

				orders := dao.Orders{}
				orders.OrderNo = tool.UUID()
				orders.UserID = user.ID
				orders.OID = organization.ID
				orders.OrdersPackageNo = op.OrderNo
				//PayMoney = 100

				orders.PayMoney = PayMoney
				orders.PostType = int(PostType)
				orders.Status = play.OS_Order
				orders.Address = util.StructToJSON(address)
				orders.DiscountMoney = uint(FullCutAll + FavouredPrice)
				orders.GoodsMoney = uint(GoodsPrice)
				orders.ExpressMoney = uint(ExpressPrice)

				err := controller.Orders.AddOrders(&orders, oggs)
				if err != nil {
					return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", nil)}
				}
				orderList = append(orderList, orders)
			}

			OutResult := make(map[string]interface{})
			OutResult["OrderNo"] = op.OrderNo
			OutResult["OrderCount"] = len(orderList)
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", OutResult)}
		} else {

			orderList := make([]dao.Orders, 0)

			for _, value := range results {

				oggs := value["OrdersGoodsInfos"].([]dao.OrdersGoodsInfo)
				//result["OrdersGoodsInfos"]=oggs
				FavouredPrice := value["FavouredPrice"].(uint64)
				FullCutAll := value["FullCutAll"].(uint64)
				GoodsPrice := value["GoodsPrice"].(uint64)
				ExpressPrice := value["ExpressPrice"].(uint64)

				organization := value["Organization"].(dao.Organization)

				PayMoney := GoodsPrice - FullCutAll + ExpressPrice //支付价格已经包含了 满减，限时抢购的扣去的部分  - _FullCutPrice-FavouredPrice

				orders := dao.Orders{}
				orders.OrderNo = tool.UUID()
				orders.UserID = user.ID
				orders.OID = organization.ID
				//orders.OrdersPackageNo = op.OrderNo
				//PayMoney = 100

				orders.PayMoney = PayMoney
				orders.PostType = int(PostType)
				orders.Status = play.OS_Order
				orders.Address = util.StructToJSON(address)
				orders.DiscountMoney = uint(FullCutAll + FavouredPrice)
				orders.GoodsMoney = uint(GoodsPrice)
				orders.ExpressMoney = uint(ExpressPrice)

				err := controller.Orders.AddOrders(&orders, oggs)
				if err != nil {
					return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", nil)}
				}
				orderList = append(orderList, orders)
			}

			OutResult := make(map[string]interface{})
			OutResult["OrderNo"] = orderList[0].OrderNo
			OutResult["OrderCount"] = len(orderList)
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", OutResult)}

		}

	} else {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: Error.Error(), Data: nil}}
	}
}

func (controller *OrderController) ordersConfirmListAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	var ogs []dao.OrdersGoods
	if context.Session.Attributes.Get(play.SessionConfirmOrders) == nil {
		ogs = make([]dao.OrdersGoods, 0)
	} else {
		ogs = *(context.Session.Attributes.Get(play.SessionConfirmOrders)).(*[]dao.OrdersGoods)
	}
	context.Request.ParseForm()

	PostType, _ := strconv.ParseInt(context.Request.FormValue("PostType"), 10, 64)
	AddressTxt := context.Request.FormValue("Address")
	address := dao.Address{}
	util.JSONToStruct(AddressTxt, &address)

	Error, results, _ := controller.Orders.AnalyseOrdersGoodsList(user.ID, address, int(PostType), ogs)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(Error, "OK", results)}
}
func (controller *OrderController) ordersCartChangeAction(context *gweb.Context) gweb.Result {
	context.Request.ParseForm()
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	GSID, _ := strconv.ParseUint(context.Request.FormValue("GSID"), 10, 64)
	Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)

	err := controller.ShoppingCart.UpdateByUserIDAndID(user.ID, GSID, uint(Quantity))
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", nil)}
}
func (controller *OrderController) ordersCartDeleteAction(context *gweb.Context) gweb.Result {
	context.Request.ParseForm()

	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	_ShoppingCartIDs := context.Request.FormValue("GSIDs")
	ShoppingCartIDs := strings.Split(_ShoppingCartIDs, ",")

	IDs := make([]uint64, 0)
	for _, value := range ShoppingCartIDs {
		ID, _ := strconv.ParseUint(value, 10, 64)
		IDs = append(IDs, ID)
	}
	err := controller.ShoppingCart.DeleteListByIDs(user.ID, IDs)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", nil)}
}
func (controller *OrderController) ordersCartListAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	err, list, _ := controller.ShoppingCart.FindShoppingCartListDetails(user.ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", list)}
}
func (controller *OrderController) ordersBuyAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	context.Request.ParseForm()
	_GSIDs := context.Request.FormValue("GSIDs")
	GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	SpecificationID, _ := strconv.ParseUint(context.Request.FormValue("SpecificationID"), 10, 64)
	Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)

	if !strings.EqualFold(_GSIDs, "") && GoodsID == 0 && SpecificationID == 0 && Quantity == 0 {
		GSIDs := strings.Split(_GSIDs, ",")
		if len(GSIDs) > 0 {
			GSIDsList := make([]uint64, 0)
			for _, value := range GSIDs {
				ID, _ := strconv.ParseUint(value, 10, 64)
				GSIDsList = append(GSIDsList, ID)
			}
			err := controller.Orders.AddCartOrdersByShoppingCartIDs(context.Session, user.ID, GSIDsList)
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
		} else {
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("没有相关ID"), "", nil)}
		}
	} else {
		if GoodsID != 0 && SpecificationID != 0 && Quantity != 0 {
			err := controller.Orders.BuyOrders(context.Session, user.ID, GoodsID, SpecificationID, uint(Quantity))
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "立即购买", nil)}
		} else {
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("订单数据出错"), "", nil)}
		}
	}
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("数据出错"), "", nil)}
}
func (controller *OrderController) ordersAddAction(context *gweb.Context) gweb.Result {

	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	context.Request.ParseForm()
	GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	SpecificationID, _ := strconv.ParseUint(context.Request.FormValue("SpecificationID"), 10, 64)
	Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)

	err := controller.Orders.AddCartOrders(user.ID, GoodsID, SpecificationID, uint(Quantity))
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "已添加到购物车", nil)}
}
