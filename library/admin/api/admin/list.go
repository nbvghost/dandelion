package admin

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type List struct {
	Admin *model.Admin `mapping:""`
	Post  struct {
		*model.Datatables
	} `method:"Post"`
}

func (m *List) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *List) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	//admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.Admin.Service.DatatablesListOrder(db.GetDB(ctx), m.Post.Datatables, &[]model.Admin{}, m.Admin.OID, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, nil
}
