package mp

import (
	"dandelion/app/play"
	"dandelion/app/service/dao"
	"dandelion/app/util"

	"github.com/nbvghost/glog"

	"dandelion/app/service"

	"dandelion/app/action/mp/store"
	"strconv"

	"dandelion/app/action/mp/order"

	"encoding/base64"

	"dandelion/app/action/mp/content"
	"dandelion/app/action/mp/journal"
	"dandelion/app/action/mp/user"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
	"github.com/skip2/go-qrcode"
)

type InterceptorMp struct {
	Organization service.OrganizationService
}

//Execute(Session *Session,Request *http.Request)(bool,Result)
func (controller InterceptorMp) Execute(context *gweb.Context) (bool, gweb.Result) {

	if context.Session.Attributes.Get(play.SessionUser) == nil {
		//context.Response.Header().Add("Login-Status", "0")
		//context.Response.Write([]byte(util.StructToJSON(&dao.ActionStatus{Success: false, Message: "没有登陆", Data: nil})))
		return false, &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "没有登陆", Data: nil}}
	} else {
		//return controller.Organization.ReadOrganization(context)
		return true, nil
	}
}

type Controller struct {
	gweb.BaseController
	User          service.UserService
	Goods         service.GoodsService
	Orders        service.OrdersService
	Store         service.StoreService
	FullCut       service.FullCutService
	ShoppingCart  service.ShoppingCartService
	ScoreGoods    service.ScoreGoodsService
	Rank          service.RankService
	CardItem      service.CardItemService
	Wx            service.WxService
	Verification  service.VerificationService
	TimeSell      service.TimeSellService
	Configuration service.ConfigurationService
	Journal       service.JournalService
}

func (controller *Controller) Apply() {
	controller.Interceptors.Add(&InterceptorMp{})

	controller.AddHandler(gweb.ALLMethod("index", controller.indexPage))

	controller.AddHandler(gweb.GETMethod("get_login_user", controller.getLoginUserAction))
	controller.AddHandler(gweb.POSMethod("get_login_user_phone", controller.getLoginUserPhoneAction))
	controller.AddHandler(gweb.GETMethod("goods_type/list", controller.goodsTypeListAction))
	controller.AddHandler(gweb.GETMethod("goods_type/child/{GoodsTypeID}/list", controller.goodsTypeChildListByGoodsTypeIDAction))

	controller.AddHandler(gweb.GETMethod("goods/child/{GoodsTypeID}/{GoodsTypeChildID}/list", controller.goodsChildByGoodsTypeIDAction))
	controller.AddHandler(gweb.GETMethod("goods/get/{ID}", controller.goodsByGoodsIDAction))
	controller.AddHandler(gweb.GETMethod("goods/hot/list", controller.goodsHotListAction))
	controller.AddHandler(gweb.GETMethod("goods/all/list", controller.goodsAllListAction))

	controller.AddHandler(gweb.GETMethod("score_goods/list", controller.scoreGoodsListAction))
	controller.AddHandler(gweb.GETMethod("score_goods/exchange/{ScoreGoodsID}", controller.scoreGoodsExchangeAction))
	//controller.AddHandler(gweb.GETMethod("fullcut/list", controller.fullcutListAction))

	controller.AddHandler(gweb.GETMethod("share/score", controller.shareScoreAction))

	controller.AddHandler(gweb.GETMethod("card/list", controller.cardListAction))
	controller.AddHandler(gweb.GETMethod("card/get/{CardItemID}", controller.cardGetAction))
	controller.AddHandler(gweb.GETMethod("verification/get/{VerificationNo}", controller.verificationGetByVerificationNoAction))

	controller.AddHandler(gweb.POSMethod("configuration/list", controller.configurationListAction))
	controller.AddHandler(gweb.GETMethod("read/share/key", controller.readShareKeyAction))

	OrderController := &order.OrderController{}
	OrderController.Interceptors = controller.Interceptors
	controller.AddSubController("/order/", OrderController)

	StoreController := &store.StoreController{}
	StoreController.Interceptors = controller.Interceptors
	controller.AddSubController("/store/", StoreController)

	JournalController := &journal.JournalController{}
	JournalController.Interceptors = controller.Interceptors
	controller.AddSubController("/journal/", JournalController)

	ContentController := &content.ContentController{}
	ContentController.Interceptors = controller.Interceptors
	controller.AddSubController("/content/", ContentController)

	UserController := &user.UserController{}
	UserController.Interceptors = controller.Interceptors
	controller.AddSubController("/user/", UserController)

}
func (controller *Controller) readShareKeyAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	ShareKey := context.Request.URL.Query().Get("ShareKey")
	UserID, ProductID := util.DecodeShareKey(ShareKey)

	result := make(map[string]interface{})
	result["UserID"] = UserID
	result["ProductID"] = ProductID
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: result}}
}
func (controller *Controller) configurationListAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var ks []uint64
	util.RequestBodyToJSON(context.Request.Body, &ks)
	list := controller.Configuration.GetConfigurations(0, ks)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}
}
func (controller *Controller) indexPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}

func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func (controller *Controller) rootPage(context *gweb.Context) gweb.Result {

	return &gweb.RedirectToUrlResult{"index"}
}

func (controller *Controller) verificationGetByVerificationNoAction(context *gweb.Context) gweb.Result {

	VerificationNo := context.PathParams["VerificationNo"]
	verification := controller.Verification.GetVerificationByVerificationNo(VerificationNo)

	if verification.StoreID > 0 && verification.Quantity > 0 {
		var store dao.Store
		controller.Store.Get(dao.Orm(), verification.StoreID, &store)
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: store}}
	} else {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "", Data: nil}}
	}

}
func (controller *Controller) cardGetAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	CardItemID, _ := strconv.ParseUint(context.PathParams["CardItemID"], 10, 64)
	var cardItem dao.CardItem
	controller.CardItem.Get(dao.Orm(), CardItemID, &cardItem)

	controller.Verification.DeleteWhere(dao.Orm(), &dao.Verification{}, "UserID=? and CardItemID=? and StoreID=? and Quantity=?", user.ID, cardItem.ID, 0, 0)

	verification := dao.Verification{}
	verification.CardItemID = cardItem.ID
	verification.UserID = user.ID
	verification.VerificationNo = tool.UUID()
	verification.Name, verification.Label = cardItem.GetNameLabel(dao.Orm())

	controller.Verification.Add(dao.Orm(), &verification)

	results := make(map[string]interface{})

	if false {
		//wxconfig := controller.Wx.MiniProgram(user.OID)
		//postData := make(map[string]interface{})

		/*access_token := controller.Wx.GetAccessToken(wxconfig.ID)
		postData["scene"] = verification.VerificationNo
		//postData["page"] = "pages/store_verification/store_verification"
		postData["width"] = 430
		postData["auto_color"] = true

		body := strings.NewReader(util.StructToJSON(postData))
		//postData := url.Values{}
		//postData.Add("scene","sdfsd")
		resp, err := http.Post("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token="+access_token, "application/json", body)
		if err != nil {
			return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: err.Error(), Data: nil}}
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: err.Error(), Data: nil}}
		}

		defer resp.Body.Close()

		imageString := "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)

		results["QRCodeBase64"] = imageString*/
	} else {

		png, _ := qrcode.Encode(verification.VerificationNo, qrcode.Low, 256)
		imageString := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
		results["QRCodeBase64"] = imageString

	}

	results["Verification"] = verification
	results["CardItem"] = cardItem
	results["ExpireTime"] = cardItem.ExpireTime
	results["HasQuantity"] = cardItem.Quantity - cardItem.UseQuantity

	if (cardItem.Quantity - cardItem.UseQuantity) <= 0 {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "数量不足，无法核销", Data: nil}}
	} else {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: results}}
	}

}
func (controller *Controller) shareScoreAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	if Share, have := context.Data["Share"]; have {
		err := controller.Journal.AddScoreJournal(dao.Orm(),
			user.ID,
			"转发与分享送积分", "转发与分享",
			play.ScoreJournal_Type_Share, int64(Share.(float64)), dao.KV{})
		glog.Error(err)

	}
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: nil}}
}
func (controller *Controller) cardListAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	list := controller.CardItem.FindByUserID(user.ID)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: list}}
}
func (controller *Controller) scoreGoodsExchangeAction(context *gweb.Context) gweb.Result {
	//ScoreGoodsID
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	ScoreGoodsID, _ := strconv.ParseUint(context.PathParams["ScoreGoodsID"], 10, 64)
	err := controller.ScoreGoods.Exchange(user.ID, ScoreGoodsID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "兑换成功", nil)}
}
func (controller *Controller) scoreGoodsListAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	list := controller.ScoreGoods.ListScoreGoods()
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: list}}
}

/*func (controller *Controller) fullcutListAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	fullcuts := controller.FullCut.FindOrderByAmountASC(Orm)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "", fullcuts)}
}*/

func (controller *Controller) goodsHotListAction(context *gweb.Context) gweb.Result {
	index, _ := strconv.Atoi(context.Request.URL.Query().Get("index"))
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: controller.Goods.GoodsList(user.ID, "CountSale desc", index, "Hide=?", 0)}}

	//return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: controller.Goods.HotList()}}
}
func (controller *Controller) goodsAllListAction(context *gweb.Context) gweb.Result {

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: controller.Goods.AllList()}}
}
func (controller *Controller) goodsByGoodsIDAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	goodsInfo := controller.Goods.GetGoods(Orm, ID)



	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: goodsInfo}}
}
func (controller *Controller) goodsTypeListAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: controller.Goods.ListAllGoodsType()}}
}
func (controller *Controller) goodsTypeChildListByGoodsTypeIDAction(context *gweb.Context) gweb.Result {
	GoodsTypeID, _ := strconv.ParseUint(context.PathParams["GoodsTypeID"], 10, 64)
	results := controller.Goods.ListGoodsTypeChild(GoodsTypeID)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: results}}
}
func (controller *Controller) goodsChildByGoodsTypeIDAction(context *gweb.Context) gweb.Result {
	GoodsTypeID, _ := strconv.ParseUint(context.PathParams["GoodsTypeID"], 10, 64)
	GoodsTypeChildID, _ := strconv.ParseUint(context.PathParams["GoodsTypeChildID"], 10, 64)
	results := controller.Goods.ListGoodsChildByGoodsTypeID(GoodsTypeID, GoodsTypeChildID)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: results}}
}
func (controller *Controller) getLoginUserPhoneAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	userInfo := make(map[string]interface{})
	util.RequestBodyToJSON(context.Request.Body, &userInfo)

	iv := userInfo["iv"].(string)
	encryptedData := userInfo["encryptedData"].(string)

	SessionKey := context.Session.Attributes.Get(play.SessionMiniProgramKey).(string)

	su, con := controller.Wx.Decrypt(encryptedData, SessionKey, iv)

	if su {
		phoneInfo := make(map[string]interface{})
		util.JSONToStruct(con, &phoneInfo)
		user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
		user.Tel = phoneInfo["phoneNumber"].(string)
		controller.User.ChangeModel(Orm, user.ID, &dao.User{Tel: user.Tel})
		context.Session.Attributes.Put(play.SessionUser, user)
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "绑定成功", Data: user}}
	} else {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "绑定失败", Data: nil}}
	}
}
func (controller *Controller) getLoginUserAction(context *gweb.Context) gweb.Result {
	if context.Session.Attributes.Get(play.SessionUser) == nil {
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: false, Message: "", Data: nil}}
	} else {

		user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

		controller.User.Get(dao.Orm(), user.ID, &user)

		results := make(map[string]interface{})
		results["User"] = user
		results["MyShareKey"] = util.EncodeShareKey(user.ID, 0) //tool.Hashids{}.Encode(user.ID) //tool.CipherEncrypterData(strconv.Itoa(int(user.ID)))
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: results}}
	}
}
