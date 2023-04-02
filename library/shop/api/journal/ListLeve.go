package journal

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/journal"
	"github.com/nbvghost/dandelion/service/user"
)

type ListLeve struct {
	UserService    user.UserService
	JournalService journal.JournalService
	User           *model.User `mapping:""`
}

func (controller *ListLeve) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (controller *ListLeve) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	user := controller.User

	leve1UserIDs := controller.UserService.Leve1(user.ID)
	leve1 := controller.JournalService.ListUserJournalLeveBrokerage(user.ID, leve1UserIDs)

	leve2UserIDs := controller.UserService.Leve2(leve1UserIDs)
	leve2 := controller.JournalService.ListUserJournalLeveBrokerage(user.ID, leve2UserIDs)

	leve3UserIDs := controller.UserService.Leve3(leve2UserIDs)
	leve3 := controller.JournalService.ListUserJournalLeveBrokerage(user.ID, leve3UserIDs)

	leve4UserIDs := controller.UserService.Leve4(leve3UserIDs)
	leve4 := controller.JournalService.ListUserJournalLeveBrokerage(user.ID, leve4UserIDs)

	leve5UserIDs := controller.UserService.Leve5(leve4UserIDs)
	leve5 := controller.JournalService.ListUserJournalLeveBrokerage(user.ID, leve5UserIDs)

	leve6UserIDs := controller.UserService.Leve6(leve5UserIDs)
	leve6 := controller.JournalService.ListUserJournalLeveBrokerage(user.ID, leve6UserIDs)

	results := make(map[string]interface{})

	fmt.Println(leve1, leve2, leve3, leve4, leve5, leve6)

	results["leve1"] = leve1
	results["leve2"] = leve2
	results["leve3"] = leve3
	results["leve4"] = leve4
	results["leve5"] = leve5
	results["leve6"] = leve6

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}, nil

}
