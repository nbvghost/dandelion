package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Config struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		FieldName  string `form:"FieldName"`
		FieldValue string `form:"FieldValue"`
	} `method:"POST"`
	Get struct{} `method:"Get"`
}

func (m *Config) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	//context.Request.ParseForm()
	//fieldName := context.Request.FormValue("FieldName")
	//fieldValue := context.Request.FormValue("FieldValue")

	err = service.Content.ChangeContentConfig(ctx, m.Organization.ID, m.POST.FieldName, m.POST.FieldValue)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "设置成功", nil)}, err
}

func (m *Config) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	contentConfig := repository.ContentConfigDao.GetContentConfig(db.GetDB(ctx), m.Organization.ID)
	return &result.JsonResult{Data: &result.ActionResult{
		Code:    result.Success,
		Message: "OK",
		Data:    contentConfig,
	}}, nil
}
