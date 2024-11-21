package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/shop/domain/module"
	"github.com/nbvghost/dandelion/service/serviceargument"

)

type ContentPagination[T module.ListType] struct {
	Pagination      serviceargument.Pagination[T]   `arg:""`
	CurrentMenuData serviceargument.CurrentMenuData `arg:""`
}

func (m *ContentPagination[T]) Template() ([]byte, error) {
	return nil, nil
}

func (m *ContentPagination[T]) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{"Pagination": m.Pagination, "CurrentMenuData": m.CurrentMenuData}, nil
}
