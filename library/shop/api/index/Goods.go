package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Goods struct {
	Get struct {
		ID dao.PrimaryKey `uri:"ID"`
	} `method:"get"`
}

func (m *Goods) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	Orm := db.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	goodsInfo, err := service.Goods.Goods.GetGoods(Orm, ctx, dao.PrimaryKey(m.Get.ID))

	err = dao.UpdateByPrimaryKey(db.Orm(), entity.Goods, goodsInfo.Goods.ID, &model.Goods{CountView: goodsInfo.Goods.CountView + 1})
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: goodsInfo}}, err
}
