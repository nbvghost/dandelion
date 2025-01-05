package template

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type List struct {
	Organization *model.Organization `mapping:""`
}

func (m *List) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.Orm()
	//var list []model.ExpressTemplate
	//err = m.ExpressTemplate.FindAllByOID(Orm, &list, m.Organization.ID)
	list := dao.Find(Orm, entity.ExpressTemplate).Where(`"OID"=?`, m.Organization.ID).List()
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", list)}, err
}
