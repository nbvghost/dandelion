package admin

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Admin struct {
	Admin *model.Admin `mapping:""`
	Post  struct {
		*model.Admin
	} `method:"Post"`
}

func (m *Admin) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Admin) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {

	//todo
	err = service.Admin.Service.AddItem(ctx, m.Admin.OID, m.Post.Admin)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}, err
}
