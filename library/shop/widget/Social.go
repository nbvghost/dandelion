package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
)

type Social struct {
	Organization  model.Organization  `arg:""`
	ContentConfig model.ContentConfig `arg:""`
}

func (m *Social) Template() ([]byte, error) {
	return nil, nil
}

func (m *Social) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{
		"ContentConfig": m.ContentConfig,
		"Organization":  m.Organization,
	}, nil
}
