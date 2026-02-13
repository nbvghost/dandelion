package collage

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type GoodsGoodsID struct {
	Organization *model.Organization `mapping:""`
	Delete       struct {
		GoodsID uint `uri:"GoodsID"`
	} `method:"Delete"`
}

func (m *GoodsGoodsID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *GoodsGoodsID) HandleDelete(ctx constrain.IContext) (r constrain.IResult, err error) {

	/*Orm := db.GetDB(ctx)
	ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	list := GlobalService.Goods.DeleteTimeSellGoods(Orm, ID, company.ID)

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "删除成功", list)}
	*/

	Orm := db.GetDB(ctx)
	//ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["GoodsID"])
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	list := service.Goods.Goods.DeleteCollageGoods(ctx, Orm, dao.PrimaryKey(m.Delete.GoodsID), m.Organization.ID)

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "删除成功", list)}, err
}
