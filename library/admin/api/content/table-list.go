package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type TableList struct {
	POST struct {
		*model.Datatables
	} `method:"POST"`
}

func (m *TableList) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *TableList) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.Orm()
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.Content.DatatablesListOrder(Orm, m.POST.Datatables, &[]model.Content{}, 0, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, nil
}
