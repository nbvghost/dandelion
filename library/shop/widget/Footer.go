package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
)

type Footer struct {
	Organization  model.Organization  `arg:""`
	ContentConfig model.ContentConfig `arg:""`
	AllMenusData  extends.MenusData   `arg:""`
	PageMenus     []extends.Menus     `arg:""`
}

func (m *Footer) Template() ([]byte, error) {
	return nil, nil
}

func (m *Footer) Render(ctx constrain.IContext) (map[string]any, error) {

	return map[string]any{
		"Menus":         m.AllMenusData,
		"Pages":         m.PageMenus,
		"ContentConfig": m.ContentConfig,
		"Organization":  m.Organization,
	}, nil
}
