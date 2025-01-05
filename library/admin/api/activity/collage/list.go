package collage

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
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

func (m *List) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.Orm()
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)

	m.Post.Datatables.Groupbys = make([]string, 0)
	m.Post.Datatables.Groupbys = append(m.Post.Datatables.Groupbys, "Hash")

	draw, recordsTotal, recordsFiltered, list := service.Activity.Collage.DatatablesListOrder(Orm, m.Post.Datatables, &[]model.Collage{}, m.Organization.ID, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, err
}
