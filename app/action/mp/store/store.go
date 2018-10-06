package store

import (
	"dandelion/app/play"
	"dandelion/app/service/dao"
	"strconv"

	"dandelion/app/service"

	"dandelion/app/util"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type StoreController struct {
	gweb.BaseController
	Store        service.StoreService
	StoreStock   service.StoreStockService
	Wx           service.WxService
	Orders       service.OrdersService
	Journal      service.JournalService
	Verification service.VerificationService
	CardItem     service.CardItemService
	Transfers    service.TransfersService
}

func (controller *StoreController) Apply() {

	controller.AddHandler(gweb.GETMethod("/location/list", controller.storeLocationListAction))
	controller.AddHandler(gweb.GETMethod("/get", controller.storeGetAction))
	controller.AddHandler(gweb.GETMethod("/get/:StoreID", controller.storeGetIDAction))
	controller.AddHandler(gweb.POSMethod("/supply", controller.supplyAction))
	controller.AddHandler(gweb.POSMethod("/verification", controller.verificationAction))
	controller.AddHandler(gweb.GETMethod("/verification/get/:VerificationNo", controller.verificationGetByVerificationNoAction))
	controller.AddHandler(gweb.GETMethod("/list/stock", controller.listStockAction))
	controller.AddHandler(gweb.GETMethod("/list/stock/goods/specification/:GoodsID", controller.listStockSpecificationsAction))

	controller.AddHandler(gweb.POSMethod("/journal/list", controller.journalListAction))
	controller.AddHandler(gweb.POSMethod("/transfers", controller.transfersAction))
	controller.AddHandler(gweb.POSMethod("/add/star", controller.addStarAction))

}
func (controller *StoreController) addStarAction(context *gweb.Context) gweb.Result {

	context.Request.ParseForm()
	StoreID, _ := strconv.ParseUint(context.Request.FormValue("StoreID"), 10, 64)
	Num, _ := strconv.ParseUint(context.Request.FormValue("Num"), 10, 64)

	var store dao.Store
	controller.Store.Get(dao.Orm(), StoreID, &store)
	if Num > 5 {
		Num = 5
	}
	store.Stars = store.Stars + Num

	store.StarsCount = store.StarsCount + 1
	err := controller.Store.ChangeModel(dao.Orm(), store.ID, &dao.Store{Stars: store.Stars, StarsCount: store.StarsCount})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "评价成功", nil)}

}
func (controller *StoreController) transfersAction(context *gweb.Context) gweb.Result {
	store := context.Session.Attributes.Get(play.SessionStore).(*dao.Store)
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	context.Request.ParseForm()
	ReUserName := context.Request.FormValue("ReUserName")

	IP := util.GetIP(context)
	err := controller.Transfers.StoreTransfers(store.ID, user.ID, ReUserName, IP)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "提现申请成功，请查看到账通知结果", nil)}
}
func (controller *StoreController) verificationGetByVerificationNoAction(context *gweb.Context) gweb.Result {

	VerificationNo := context.PathParams["VerificationNo"]
	verification := controller.Verification.GetVerificationByVerificationNo(VerificationNo)

	var cardItem dao.CardItem
	controller.CardItem.Get(dao.Orm(), verification.CardItemID, &cardItem)

	if verification.ID == 0 {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "", Data: nil}}
	}

	results := make(map[string]interface{})
	results["CardItem"] = cardItem
	results["Verification"] = verification

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: results}}
}
func (controller *StoreController) journalListAction(context *gweb.Context) gweb.Result {
	store := context.Session.Attributes.Get(play.SessionStore).(*dao.Store)
	context.Request.ParseForm()
	//StoreID, _ := strconv.ParseUint(context.Request.FormValue("StoreID"), 10, 64)
	StartDate := context.Request.FormValue("StartDate")
	EndDate := context.Request.FormValue("EndDate")

	list := controller.Journal.StoreListJournal(store.ID, StartDate, EndDate)

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: list}}
}
func (controller *StoreController) supplyAction(context *gweb.Context) gweb.Result {

	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	if context.Session.Attributes.Get(play.SessionUser) == nil || context.Session.Attributes.Get(play.SessionStore) == nil {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "登陆信息已经失效，请重新登陆", Data: nil}}
	}

	store := context.Session.Attributes.Get(play.SessionStore).(*dao.Store)
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	context.Request.ParseForm()

	PayMoney, _ := strconv.ParseUint(context.Request.FormValue("PayMoney"), 10, 64)
	if PayMoney <= 0 {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "无效的金额", Data: nil}}
	}
	ip := util.GetIP(context)

	supply := dao.SupplyOrders{}
	supply.StoreID = store.ID
	supply.OrderNo = tool.UUID()
	supply.PayMoney = PayMoney
	supply.UserID = user.ID
	supply.Type = play.SupplyType_Store

	WxConfig := controller.Wx.MiniProgram()

	Success, Message, result := controller.Wx.Order(supply.OrderNo, "门店", "充值", "", user.OpenID, ip, PayMoney, play.OrdersType_Supply, WxConfig)
	if Success == false {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: Success, Message: Message, Data: result}}
	}

	controller.Orders.Add(dao.Orm(), &supply)

	//WxConfig := controller.Wx.MiniProgram()

	outData := controller.Wx.GetWXAConfig(result.Prepay_id, WxConfig)

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: outData}}
}
func (controller *StoreController) verificationAction(context *gweb.Context) gweb.Result {
	context.Request.ParseForm()
	//self
	Action := context.Request.FormValue("Action")
	switch Action {
	case "User":
		//核销卡卷
		store := context.Session.Attributes.Get(play.SessionStore).(*dao.Store)
		user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
		Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)
		VerificationNo := context.Request.FormValue("VerificationNo")
		//verification := controller.Verification.GetVerificationByVerificationNo(VerificationNo)
		tx := dao.Orm().Begin()
		err := controller.Verification.VerificationCardItem(tx, VerificationNo, uint(Quantity), user, store)
		if err != nil {
			tx.Rollback()
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
		} else {
			tx.Commit()
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "核销成功", nil)}
		}
		//fmt.Println(verification)
	case "Self":
		store := context.Session.Attributes.Get(play.SessionStore).(*dao.Store)
		StoreStockID, _ := strconv.ParseUint(context.Request.FormValue("StoreStockID"), 10, 64)
		Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)
		as := controller.StoreStock.VerificationSelf(store.ID, StoreStockID, Quantity)
		return &gweb.JsonResult{Data: as}

	}
	return &gweb.JsonResult{Data: &dao.ActionStatus{}}
}
func (controller *StoreController) listStockSpecificationsAction(context *gweb.Context) gweb.Result {
	//GoodsID
	GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)

	store := context.Session.Attributes.Get(play.SessionStore).(*dao.Store)

	list := controller.StoreStock.ListStoreSpecifications(store.ID, GoodsID)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: list}}

}
func (controller *StoreController) listStockAction(context *gweb.Context) gweb.Result {

	store := context.Session.Attributes.Get(play.SessionStore).(*dao.Store)

	list := controller.StoreStock.ListStoreStock(store.ID)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: list}}

}
func (controller *StoreController) storeGetIDAction(context *gweb.Context) gweb.Result {
	StoreID, _ := strconv.ParseUint(context.PathParams["StoreID"], 10, 64)
	var Store dao.Store
	controller.Store.Get(dao.Orm(), StoreID, &Store)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: Store}}
}
func (controller *StoreController) storeGetAction(context *gweb.Context) gweb.Result {

	if context.Session.Attributes.Get(play.SessionStore) == nil {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "没有权限访问门店", Data: nil}}
	} else {
		Store := context.Session.Attributes.Get(play.SessionStore).(*dao.Store)
		controller.Store.Get(dao.Orm(), Store.ID, Store)
		context.Session.Attributes.Put(play.SessionStore, Store)
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: Store}}
	}
}
func (controller *StoreController) storeLocationListAction(context *gweb.Context) gweb.Result {

	Latitude, err := strconv.ParseFloat(context.Request.URL.Query().Get("Latitude"), 64)
	tool.CheckError(err)
	Longitude, err := strconv.ParseFloat(context.Request.URL.Query().Get("Longitude"), 64)
	tool.CheckError(err)

	list := controller.Store.LocationList(Latitude, Longitude)

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: list}}

}
