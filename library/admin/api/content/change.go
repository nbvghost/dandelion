package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type Change struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		ID         dao.PrimaryKey
		FieldName  string
		FieldValue any
	} `method:"POST"`
	Get struct{} `method:"Get"`
}

func (m *Change) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	//context.Request.ParseForm()
	//fieldName := context.Request.FormValue("FieldName")
	//fieldValue := context.Request.FormValue("FieldValue")
	if m.POST.ID == 0 {
		return nil, result.NewErrorText("设置失败")
	}

	err := repository.ContentDao.ChangeContentByField(ctx, m.POST.ID, m.POST.FieldName, m.POST.FieldValue)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "设置成功", nil)}, err
}

func (m *Change) Handle(context constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
