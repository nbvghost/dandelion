package journal

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type ListLeve struct {
	User *model.User `mapping:""`
}

func (controller *ListLeve) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (controller *ListLeve) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//u := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	u := controller.User

	leve1UserIDs := service.User.Leve1(ctx, u.ID)
	leve1 := service.Journal.ListUserJournalLeveBrokerage(ctx, u.ID, leve1UserIDs)

	leve2UserIDs := service.User.Leve2(ctx, leve1UserIDs)
	leve2 := service.Journal.ListUserJournalLeveBrokerage(ctx, u.ID, leve2UserIDs)

	leve3UserIDs := service.User.Leve3(ctx, leve2UserIDs)
	leve3 := service.Journal.ListUserJournalLeveBrokerage(ctx, u.ID, leve3UserIDs)

	leve4UserIDs := service.User.Leve4(ctx, leve3UserIDs)
	leve4 := service.Journal.ListUserJournalLeveBrokerage(ctx, u.ID, leve4UserIDs)

	leve5UserIDs := service.User.Leve5(ctx, leve4UserIDs)
	leve5 := service.Journal.ListUserJournalLeveBrokerage(ctx, u.ID, leve5UserIDs)

	leve6UserIDs := service.User.Leve6(ctx, leve5UserIDs)
	leve6 := service.Journal.ListUserJournalLeveBrokerage(ctx, u.ID, leve6UserIDs)

	results := make(map[string]interface{})

	results["leve1"] = leve1
	results["leve2"] = leve2
	results["leve3"] = leve3
	results["leve4"] = leve4
	results["leve5"] = leve5
	results["leve6"] = leve6

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}, nil

}
