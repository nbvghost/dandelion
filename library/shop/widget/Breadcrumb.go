package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
)

type Breadcrumb struct {
	Organization *model.Organization `mapping:""`
	Navigations  []extends.Menus     `arg:""`
	Align        string              `arg:""`
}

func (m *Breadcrumb) Template() ([]byte, error) {
	return nil, nil
}

func (m *Breadcrumb) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{
		"Align": m.Align,
		"List":  m.Navigations,
	}, nil
}
