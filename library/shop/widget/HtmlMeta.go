package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
)

type HtmlMeta struct {
	Organization *model.Organization `mapping:""`

	HtmlMeta extends.HtmlMeta `arg:""`
}

func (m *HtmlMeta) Template() ([]byte, error) {
	return nil, nil
}

func (m *HtmlMeta) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{
		"HtmlMeta": m.HtmlMeta,
	}, nil
}
