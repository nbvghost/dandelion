package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/user"
)

type Info struct {
	UserService     user.UserService
	StoreService    company.StoreService
	RankService     activity.RankService
	OrdersService   order.OrdersService
	CardItemService activity.CardItemService
	User            *model.User `mapping:""`

	Get struct {
	} `method:"Get"`
	Put struct {
		ChangeAllowAssistance bool
		AllowAssistance       bool
		Subscribe             bool
		ChangeSubscribe       bool
	} `method:"Put"`
}

func (m *Info) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	store := m.StoreService.GetByPhone(m.User.Phone)

	leve1UserIDs := m.UserService.Leve1(m.User.ID)
	leve2UserIDs := m.UserService.Leve2(leve1UserIDs)

	results := make(map[string]interface{})
	results["Store"] = store
	results["User"] = m.User
	results["Leve1Count"] = len(leve1UserIDs)
	results["Leve2Count"] = len(leve2UserIDs)

	ranks := m.RankService.FindDESC()
	for i, v := range ranks {

		if m.User.Growth >= v.GrowMaxValue {
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

	ACount := m.OrdersService.ListOrdersStatusCount(m.User.ID, []string{"Order"})
	BCount := m.OrdersService.ListOrdersStatusCount(m.User.ID, []string{"Pay"})
	CCount := m.OrdersService.ListOrdersStatusCount(m.User.ID, []string{"Deliver"})
	DCount := m.OrdersService.ListOrdersStatusCount(m.User.ID, []string{"OrderOk"})
	ECount := m.CardItemService.ListNewCount(m.User.ID)

	results["ACount"] = ACount
	results["BCount"] = BCount
	results["CCount"] = CCount
	results["DCount"] = DCount
	results["ECount"] = ECount

	//context.Session.Attributes.Put(gweb.AttributesKey(string(play.SessionStore)), &store)

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}, nil

}

func (m *Info) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	userInfo := m.UserService.GetUserInfo(context.UID())
	if m.Put.ChangeAllowAssistance {
		userInfo.SetAllowAssistance(m.Put.AllowAssistance)
		//changeMap["AllowAssistance"] = m.Put.AllowAssistance
	}
	if m.Put.ChangeSubscribe {
		//changeMap["Subscribe"] = m.Put.Subscribe
		userInfo.SetSubscribe(m.Put.Subscribe)
	}

	err = userInfo.Update(db.Orm())
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("OK"), nil
}
