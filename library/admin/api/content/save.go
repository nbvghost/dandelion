package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Save struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		*model.Content
	} `method:"POST"`
}

func (m *Save) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Save) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	isAdd := m.POST.Content.IsZero()
	err = service.Content.SaveContent(ctx, m.Organization.ID, m.POST.Content)
	as := result.NewData(m.POST.Content)
	if isAdd {
		as.Message = "添加成功"
	} else {
		as.Message = "修改成功"
	}
	return &result.JsonResult{Data: as}, err

}
