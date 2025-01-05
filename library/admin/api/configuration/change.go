package configuration

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/library/result"
)

type Change struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		Configurations []*model.Configuration `body:""`
	} `method:"POST"`
}

func (m *Change) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Change) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//item := model.Configuration{}
	//util.RequestBodyToJSON(context.Request.Body, &item)
	tx := db.Orm().Begin()
	for i := range m.POST.Configurations {
		item := m.POST.Configurations[i]
		err = service.Configuration.ChangeConfiguration(tx, m.Organization.ID, item.K, item.V)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
