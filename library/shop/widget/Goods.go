package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
)

type ShowType string

const (
	HorizontalShowType    ShowType = "h"
	VerticalSmallShowType ShowType = "h-s"
	VerticalShowType      ShowType = "v"
)

type Goods struct {
	ShowType ShowType     `arg:""`
	Goods    *model.Goods `arg:""`
}

func (m *Goods) Template() ([]byte, error) {
	return nil, nil
}

func (m *Goods) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{
		"Goods":    m.Goods,
		"ShowType": m.ShowType,
	}, nil
}
