package timesell

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type GoodsHashList struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		Hash       string            `uri:"Hash"`
		Datatables *model.Datatables `body:""`
	} `method:"post"`
}

func (m *GoodsHashList) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *GoodsHashList) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	//Hash := context.PathParams["Hash"]
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//dts := &model.Datatables{}
	//dts.Draw = 10
	//dts.Length = play.Paging
	//util.RequestBodyToJSON(context.Request.Body, dts)
	GoodsIDs := []uint{}
	db.GetDB(ctx).Model(&model.TimeSellGoods{}).Where("TimeSellHash=? and OID=?", m.Post.Hash, m.Organization.ID).Pluck("GoodsID", &GoodsIDs)
	if len(GoodsIDs) == 0 {
		GoodsIDs = []uint{0}
	}
	m.Post.Datatables.InIDs = GoodsIDs
	draw, recordsTotal, recordsFiltered, list := service.Activity.TimeSell.DatatablesListOrder(db.GetDB(ctx), m.Post.Datatables, &[]model.Goods{}, m.Organization.ID, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, err

	/*Hash := context.PathParams["Hash"]
	company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.GetDB(ctx)
	dts := &model.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	//dts.Groupbys = make([]string, 0)
	//dts.Groupbys = append(dts.Groupbys, "Hash")

	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]model.TimeSellGoods{}, company.ID, "TimeSellHash=?", Hash)
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}*/

	/*Hash := context.PathParams["Hash"]

	list := GlobalService.Goods.FindGoodsByTimeSellHash(Hash)
	//var content_item model.ExpressTemplate
	//err := controller.ExpressTemplate.Get(service.Orm, ID, &content_item)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "", list)}*/
	//2002
}
