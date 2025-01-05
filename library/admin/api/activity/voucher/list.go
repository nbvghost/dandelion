package voucher

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/library/db"
)

type List struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		Datatables *model.Datatables `body:""`
	} `method:"Post"`
}

func (m *List) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}
func (m *List) HandleGet(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}
func (m *List) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.Orm()
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.Activity.Voucher.DatatablesListOrder(Orm, m.Post.Datatables, &[]model.Voucher{}, m.Organization.ID, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, err
}
