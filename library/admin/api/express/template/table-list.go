package template

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"

	"github.com/nbvghost/dandelion/library/db"
)

type TableList struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
	} `method:"POST"`
}

func (m *TableList) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *TableList) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {

	Orm := db.GetDB(ctx)
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)
	//draw, recordsTotal, recordsFiltered, list := m.ExpressTemplate.DatatablesListOrder(Orm, m.POST.Datatables, &[]model.ExpressTemplate{}, m.Organization.ID, "")
	list := dao.Find(Orm, &model.ExpressTemplate{}).Where(`"OID"=?`, m.Organization.ID).List()
	return &result.JsonResult{Data: map[string]interface{}{"data": list}}, nil
}
