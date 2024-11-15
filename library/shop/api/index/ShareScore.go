package index

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"log"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
)

type ShareScore struct {
	//User *model.User `mapping:""`
}

func (m *ShareScore) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	//Share := 50 //config.Config.Share

	err := service.Journal.AddScoreJournal(db.Orm(), ctx.UID(), "转发与分享送积分", "转发与分享", model.ScoreJournal_Type_Share, 50)
	if err != nil {
		log.Println(err)
	}

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: nil}}, nil
}
