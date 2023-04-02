package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/activity"
)

type CardList struct {
	CardItemService activity.CardItemService
	User            *model.User `mapping:""`
}

func (m *CardList) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	list := m.CardItemService.FindByUserID(m.User.ID)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: list}}, nil
}
