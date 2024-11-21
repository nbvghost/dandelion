package widget

import (
	_ "embed"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
)

type Menus struct {
	ContentConfig model.ContentConfig `arg:""`
	MenusData     extends.MenusData   `arg:""`
	CurrentMenu   extends.Menus       `arg:""`
}

func (m *Menus) Template() ([]byte, error) {
	return nil, nil
}

func (m *Menus) Render(ctx constrain.IContext) (map[string]any, error) {

	return map[string]any{
		"ContentConfig": m.ContentConfig,
		"MenusData":     m.MenusData,
		"CurrentMenu":   m.CurrentMenu,
	}, nil
}
