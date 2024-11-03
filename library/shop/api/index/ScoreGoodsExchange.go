package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type ScoreGoodsExchange struct {
	User *model.User `mapping:""`
	Get  struct {
		ScoreGoodsID dao.PrimaryKey `uri:"ScoreGoodsID"`
	} `method:"get"`
}

func (m *ScoreGoodsExchange) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//ScoreGoodsID
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//ScoreGoodsID, _ := strconv.ParseUint(context.PathParams["ScoreGoodsID"], 10, 64)
	//ScoreGoodsID := object.ParseUint(context.PathParams["ScoreGoodsID"])
	err := service.Activity.ScoreGoods.Exchange(m.User, dao.PrimaryKey(m.Get.ScoreGoodsID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "兑换成功", nil)}, nil
}
