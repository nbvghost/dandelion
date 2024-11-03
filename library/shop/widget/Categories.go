package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type Categories struct {
	CurrentMenuData serviceargument.CurrentMenuData `arg:""`
	Tags            []extends.Tag                   `arg:""`
}

func (m *Categories) Template() ([]byte, error) {
	return nil, nil
}

func (m *Categories) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{"CurrentMenuData": m.CurrentMenuData, "Tags": m.Tags}, nil
}
