package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/activity"
)

type ScoreGoodsList struct {
	ScoreGoodsService activity.ScoreGoodsService
}

func (m *ScoreGoodsList) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)
	list := m.ScoreGoodsService.ListScoreGoods()
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: list}}, nil
}
