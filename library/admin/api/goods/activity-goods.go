package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/library/db"
)

type ActivityGoods struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		Datatables *model.Datatables `body:""`
	} `method:"POST"`
}

func (m *ActivityGoods) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *ActivityGoods) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {

	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//Hash := context.Request.URL.Query().Get("Hash")
	//dts := &model.Datatables{}
	//dts.Draw = 10
	//dts.Length = play.Paging
	//util.RequestBodyToJSON(context.Request.Body, dts)
	var TimeSellGoodsIDs []uint
	db.GetDB(ctx).Model(&model.TimeSellGoods{}).Where("OID=?", m.Organization.ID).Pluck("GoodsID", &TimeSellGoodsIDs)
	var CollageGoodsIDs []uint
	db.GetDB(ctx).Model(&model.CollageGoods{}).Where("OID=?", m.Organization.ID).Pluck("GoodsID", &CollageGoodsIDs)
	activityGoods := make([]uint, 0)
	activityGoods = append(activityGoods, TimeSellGoodsIDs...)
	activityGoods = append(activityGoods, CollageGoodsIDs...)
	m.POST.Datatables.NotIDs = activityGoods
	draw, recordsTotal, recordsFiltered, list := service.Goods.Goods.DatatablesListOrder(db.GetDB(ctx), m.POST.Datatables, &[]model.Goods{}, m.Organization.ID, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, err
}
