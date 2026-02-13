package admin

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type RoleList struct {
	Admin *model.Admin `mapping:""`
	Post  struct{}     `method:"Post"`
}

func (m *RoleList) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	admin := service.Admin.Service.GetAdmin(ctx, ctx.UID())
	if admin.Initiator {
		return result.NewData(map[string]any{"RoleList": []model.Role{}}), nil
	}
	return result.NewData(map[string]any{"RoleList": []model.Role{}}), nil
}
