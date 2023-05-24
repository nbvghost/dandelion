package index

import (
	"log"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/journal"
)

type ShareScore struct {
	JournalService journal.JournalService
	//User *model.User `mapping:""`
}

func (m *ShareScore) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	Share := 50 //config.Config.Share

	err := m.JournalService.AddScoreJournal(singleton.Orm(),
		ctx.UID(),
		"转发与分享送积分", "转发与分享",
		play.ScoreJournal_Type_Share, int64(Share), extends.KV{})
	if err != nil {
		log.Println(err)
	}

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: nil}}, nil
}
