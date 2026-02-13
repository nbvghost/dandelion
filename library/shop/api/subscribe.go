package api

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Subscribe struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		Email string `form:"Email"`
	} `method:"Post"`
}

func (m *Subscribe) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Subscribe) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	subscribe := &model.Subscribe{}
	err := db.GetDB(ctx).Model(model.LeaveMessage{}).Where(map[string]any{"Email": m.Post.Email}).First(subscribe).Error
	if !subscribe.IsZero() {
		return &result.JsonResult{Data: result.ActionResult{}}, err
	}
	subscribe.OID = m.Organization.ID
	subscribe.Email = m.Post.Email
	err = db.GetDB(ctx).Model(model.Subscribe{}).Create(subscribe).Error
	return &result.JsonResult{Data: result.ActionResult{}}, err
}
