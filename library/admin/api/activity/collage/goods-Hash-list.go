package collage

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type GoodsHashList struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		Datatables *model.Datatables `body:""`
		Hash       string            `form:"Hash"`
	} `method:"POST"`
}

func (m *GoodsHashList) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *GoodsHashList) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {

	//Hash := context.PathParams["Hash"]

	//list := GlobalService.Goods.FindGoodsByCollageHash(Hash)
	//var content_item model.ExpressTemplate
	//err := controller.ExpressTemplate.Get(service.Orm, ID, &content_item)
	//return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "", list)}
	//2002

	//Hash := context.PathParams["Hash"]
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//dts := &model.Datatables{}
	//dts.Draw = 10
	//dts.Length = play.Paging
	//util.RequestBodyToJSON(context.Request.Body, dts)
	var GoodsIDs []uint
	db.Orm().Model(&model.CollageGoods{}).Where("CollageHash=? and OID=?", m.POST.Hash, m.Organization.ID).Pluck("GoodsID", &GoodsIDs)
	if len(GoodsIDs) == 0 {
		GoodsIDs = []uint{0}
	}
	m.POST.Datatables.InIDs = GoodsIDs
	draw, recordsTotal, recordsFiltered, list := service.Goods.Goods.DatatablesListOrder(db.Orm(), m.POST.Datatables, &[]model.Goods{}, m.Organization.ID, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, err
}
