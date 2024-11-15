package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type ContentItemPagination[T serviceargument.ListType] struct {
	Pagination      serviceargument.Pagination[T]   `arg:""`
	CurrentMenuData serviceargument.CurrentMenuData `arg:""`
}

func (m *ContentItemPagination[T]) Template() ([]byte, error) {
	return nil, nil
}

func (m *ContentItemPagination[T]) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{"Pagination": m.Pagination, "CurrentMenuData": m.CurrentMenuData}, nil
}
