package journal

import (
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/dao"

	"github.com/nbvghost/gweb"

	"fmt"
)

type JournalController struct {
	gweb.BaseController
	Journal service.JournalService
	User    service.UserService
}

func (controller *JournalController) Apply() {

	controller.AddHandler(gweb.GETMethod("/list/leve", controller.listLeveAction))
	controller.AddHandler(gweb.GETMethod("/list/journal", controller.listJournalAction))

}
func (controller *JournalController) listLeveAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)

	leve1UserIDs := controller.User.Leve1(user.ID)
	leve1 := controller.Journal.ListUserJournalLeveBrokerage(user.ID, leve1UserIDs)

	leve2UserIDs := controller.User.Leve2(leve1UserIDs)
	leve2 := controller.Journal.ListUserJournalLeveBrokerage(user.ID, leve2UserIDs)

	leve3UserIDs := controller.User.Leve3(leve2UserIDs)
	leve3 := controller.Journal.ListUserJournalLeveBrokerage(user.ID, leve3UserIDs)

	leve4UserIDs := controller.User.Leve4(leve3UserIDs)
	leve4 := controller.Journal.ListUserJournalLeveBrokerage(user.ID, leve4UserIDs)

	leve5UserIDs := controller.User.Leve5(leve4UserIDs)
	leve5 := controller.Journal.ListUserJournalLeveBrokerage(user.ID, leve5UserIDs)

	leve6UserIDs := controller.User.Leve6(leve5UserIDs)
	leve6 := controller.Journal.ListUserJournalLeveBrokerage(user.ID, leve6UserIDs)

	results := make(map[string]interface{})

	fmt.Println(leve1, leve2, leve3, leve4, leve5, leve6)

	results["leve1"] = leve1
	results["leve2"] = leve2
	results["leve3"] = leve3
	results["leve4"] = leve4
	results["leve5"] = leve5
	results["leve6"] = leve6

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "", Data: results}}
}
func (controller *JournalController) listJournalAction(context *gweb.Context) gweb.Result {
	user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	///CarID, _ := strconv.ParseUint(context.PathParams["CarID"], 10, 64)
	var list []dao.UserJournal
	Orm := dao.Orm()
	controller.User.FindOrderWhere(Orm, "CreatedAt desc", &list, &dao.UserJournal{UserID: user.ID})
	//err := controller.Car.FindWhere(dao.Orm(), &list, &dao.CarRecord{CarID: CarID})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}
}
