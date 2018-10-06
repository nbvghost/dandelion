package user

import (
	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type UserController struct {
	gweb.BaseController
	Store     service.StoreService
	User      service.UserService
	Rank      service.RankService
	Transfers service.TransfersService
	Orders    service.OrdersService
	CardItem  service.CardItemService
	Wx        service.WxService
	Journal   service.JournalService
}

func (controller *UserController) Apply() {

	controller.AddHandler(gweb.GETMethod("/info", controller.userInfoAction))
	controller.AddHandler(gweb.POSMethod("/update", controller.updateAction))
	controller.AddHandler(gweb.GETMethod("/info/DaySign", controller.userInfoDaySignAction))
	controller.AddHandler(gweb.GETMethod("/growth/list/:Order", controller.userGrowthListAction))
	controller.AddHandler(gweb.GETMethod("/info/:UserID", controller.userInfoByUserIDAction))
	controller.AddHandler(gweb.POSMethod("/info/sharekey", controller.userShareKeyAction))
	controller.AddHandler(gweb.GETMethod("/poster/qrcode", controller.posterqrcodeAction))
	controller.AddHandler(gweb.POSMethod("/transfers", controller.transfersAction))

}
func (controller *UserController) posterqrcodeAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	//context.Request.ParseForm()
	//Page := context.Request.FormValue("Page")
	Page := context.Request.URL.Query().Get("Page")

	MyShareKey := util.EncodeShareKey(user.ID)

	wxconfig := controller.Wx.MiniProgram()

	access_token := controller.Wx.GetAccessToken(wxconfig)

	postData := make(map[string]interface{})
	//results := make(map[string]interface{})

	postData["scene"] = MyShareKey
	postData["page"] = Page
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
	//fmt.Println(string(b))
	defer resp.Body.Close()

	path := tool.WriteTempFile(b, "image/png")
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: path}}
	//return &gweb.ImageBytesResult{Data:b,ContentType:"image/png"}
	//imageString := "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)
	//results["QRCodeBase64"] = imageString
	//return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: results}}

}
func (controller *UserController) userShareKeyAction(context *gweb.Context) gweb.Result {

	context.Request.ParseForm()
	ShareKey := context.Request.FormValue("ShareKey")

	UserID, _ := strconv.ParseUint(tool.CipherDecrypterData(ShareKey), 10, 64)

	var user dao.User
	controller.User.Get(dao.Orm(), UserID, &user)
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: user}}
}
func (controller *UserController) transfersAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	context.Request.ParseForm()
	ReUserName := context.Request.FormValue("ReUserName")

	IP := util.GetIP(context)
	err := controller.Transfers.UserTransfers(user.ID, ReUserName, IP)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "提现申请成功，请查看到账通知结果", nil)}
}
func (controller *UserController) updateAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	context.Request.ParseForm()
	Name := context.Request.FormValue("Name")
	Age := context.Request.FormValue("Age")
	changeDataMap := make(map[string]interface{})
	if !strings.EqualFold(Name, "") {
		changeDataMap["Name"] = Name
	}
	if !strings.EqualFold(Age, "") {
		changeDataMap["Age"], _ = strconv.Atoi(Age)
	}
	//fmt.Println(user,Name,Age)
	err := controller.User.ChangeMap(dao.Orm(), user.ID, user, changeDataMap)
	//IP := util.GetIP(context)
	//err := controller.Transfers.UserTransfers(user.ID, ReUserName, IP)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (controller *UserController) userInfoByUserIDAction(context *gweb.Context) gweb.Result {
	UserID, _ := strconv.ParseUint(context.PathParams["UserID"], 10, 64)
	var user dao.User
	controller.User.Get(dao.Orm(), UserID, &user)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", user)}
}
func (controller *UserController) userGrowthListAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Order := context.PathParams["Order"]
	if strings.EqualFold(Order, "asc") {
		Order = "Growth asc"
	} else if strings.EqualFold(Order, "desc") {
		Order = "Growth desc"
	} else {
		Order = "Growth asc"
	}
	var users []dao.User
	err := controller.User.FindOrderWhereLength(dao.Orm(), Order, &users, 20, "OID=?", company.ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", users)}
}
func (controller *UserController) userInfoDaySignAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	userInfo := controller.User.GetUserInfo(user.ID)

	now := userInfo.DaySignTime
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	//d,err:=time.ParseDuration("24h")
	//tool.CheckError(err)

	fmt.Println(":", time.Now().Unix())
	fmt.Println(":", today.Unix())
	dayCount := float64(float64(time.Now().Unix()-today.Unix()) / 60 / 60 / 24) //天
	fmt.Println("天", dayCount)
	as := dao.ActionStatus{}
	if dayCount > 1 {
		//已经超过一天了，
		userInfo.DaySignTime = time.Now()
		userInfo.DaySignCount = 1
		as.Success = true
		as.Message = "打卡成功，您的打卡已经超过一天了，打卡重新累计"
		//return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("打卡成功，您的打卡已经超过一天了，打卡重新累计"), "OK", nil)}
	} else if dayCount <= 1 && dayCount >= 0 {
		//可以打卡
		userInfo.DaySignTime = time.Now()
		userInfo.DaySignCount = userInfo.DaySignCount + 1
		as.Success = true
		as.Message = "打卡成功"
	} else {
		//负数
		//已经打过卡了
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("您今天已经打卡了"), "OK", nil)}
	}

	if userInfo.DaySignCount <= 0 {
		userInfo.DaySignCount = 1
	}

	if context.Data["DaySign"] != nil {

		DaySign := context.Data["DaySign"].(map[string]interface{})
		score, have := DaySign[strconv.Itoa(userInfo.DaySignCount)]
		if have {

		} else {
			score, have = DaySign["max"]
			if !have {
				tool.Trace("打卡data.json数据没有设置DaySign.max字段值")
			}
		}
		err := controller.Journal.AddScoreJournal(dao.Orm(),
			user.ID,
			"签到送积分",
			userInfo.DaySignTime.String()+"/"+strconv.Itoa(int(score.(float64)))+"/"+strconv.Itoa(userInfo.DaySignCount),
			play.ScoreJournal_Type_DaySign, int64(score.(float64)), dao.KV{Key: "UserInfoID", Value: userInfo.ID})
		if err != nil {
			as.Success = false
			as.Message = err.Error()
		} else {
			controller.User.ChangeMap(dao.Orm(), userInfo.ID, &dao.UserInfo{}, map[string]interface{}{"DaySignTime": userInfo.DaySignTime, "DaySignCount": userInfo.DaySignCount})
		}
		return &gweb.JsonResult{Data: &as}

	} else {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("暂时无法打卡"), "OK", nil)}
	}

}
func (controller *UserController) userInfoAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	controller.User.Get(dao.Orm(), user.ID, &user)

	store := controller.Store.GetByPhone(user.Tel)

	leve1UserIDs := controller.User.Leve1(user.ID)
	leve2UserIDs := controller.User.Leve2(leve1UserIDs)

	results := make(map[string]interface{})
	results["Store"] = store
	results["User"] = user
	results["Leve1Count"] = len(leve1UserIDs)
	results["Leve2Count"] = len(leve2UserIDs)

	ranks := controller.Rank.FindDESC()
	for i, v := range ranks {

		if user.Growth >= v.GrowMaxValue {
			results["RankName"] = v.Title
			if i == 0 {
				results["RankMinValue"] = v.GrowMaxValue
				results["RankMaxValue"] = v.GrowMaxValue
			} else {
				results["RankMinValue"] = v.GrowMaxValue
				results["RankMaxValue"] = ranks[i-1].GrowMaxValue
			}
			break
		}
	}

	if _, su := results["RankName"]; su == false {

		if len(ranks) == 0 {

			results["RankName"] = ""
			results["RankMinValue"] = 0
			results["RankMaxValue"] = 100

		} else if len(ranks) >= 1 {

			rank := ranks[len(ranks)-1]
			results["RankName"] = rank.Title
			results["RankMinValue"] = 0
			results["RankMaxValue"] = rank.GrowMaxValue
		}

	}

	ACount := controller.Orders.ListOrdersStatusCount(user.ID, []string{"Order"})
	BCount := controller.Orders.ListOrdersStatusCount(user.ID, []string{"Pay"})
	CCount := controller.Orders.ListOrdersStatusCount(user.ID, []string{"Deliver"})
	DCount := controller.Orders.ListOrdersStatusCount(user.ID, []string{"OrderOk"})
	ECount := controller.CardItem.ListNewCount(user.ID)

	results["ACount"] = ACount
	results["BCount"] = BCount
	results["CCount"] = CCount
	results["DCount"] = DCount
	results["ECount"] = ECount

	context.Session.Attributes.Put(play.SessionStore, &store)

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: results}}
}
