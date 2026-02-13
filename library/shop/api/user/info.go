package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Info struct {
	User *model.User `mapping:""`

	Get struct {
	} `method:"Get"`
	Put struct {
		ChangeAllowAssistance bool
		AllowAssistance       bool
		Subscribe             bool
		ChangeSubscribe       bool
	} `method:"Put"`
}

func (m *Info) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {

	store := service.Company.Store.GetByPhone(ctx, m.User.Phone)
	leve1UserIDs := service.User.Leve1(ctx, m.User.ID)
	leve2UserIDs := service.User.Leve2(ctx, leve1UserIDs)

	results := make(map[string]interface{})
	results["Store"] = store
	results["User"] = m.User
	results["Leve1Count"] = len(leve1UserIDs)
	results["Leve2Count"] = len(leve2UserIDs)

	ranks := service.Activity.Rank.FindDESC(ctx)
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

	ACount := service.Order.Orders.ListOrdersStatusCount(ctx, m.User.ID, []string{"Order"})
	BCount := service.Order.Orders.ListOrdersStatusCount(ctx, m.User.ID, []string{"Pay"})
	CCount := service.Order.Orders.ListOrdersStatusCount(ctx, m.User.ID, []string{"Deliver"})
	DCount := service.Order.Orders.ListOrdersStatusCount(ctx, m.User.ID, []string{"OrderOk"})
	ECount := service.Activity.CardItem.ListNewCount(ctx, m.User.ID)

	results["ACount"] = ACount
	results["BCount"] = BCount
	results["CCount"] = CCount
	results["DCount"] = DCount
	results["ECount"] = ECount

	//context.Session.Attributes.Put(gweb.AttributesKey(string(play.SessionStore)), &store)

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}, nil

}

func (m *Info) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	userInfo := service.User.GetUserInfo(context, context.UID())
	if m.Put.ChangeAllowAssistance {
		userInfo.SetAllowAssistance(m.Put.AllowAssistance)
		//changeMap["AllowAssistance"] = m.Put.AllowAssistance
	}
	if m.Put.ChangeSubscribe {
		//changeMap["Subscribe"] = m.Put.Subscribe
		userInfo.SetSubscribe(m.Put.Subscribe)
	}

	err = userInfo.Update(db.GetDB(context))
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("OK"), nil
}
