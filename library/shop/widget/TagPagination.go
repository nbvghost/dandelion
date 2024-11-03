package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type TagPagination[T serviceargument.ListType] struct {
	Pagination serviceargument.Pagination[T] `arg:""`
	Tag        extends.Tag                   `arg:""`
	Order      string                        `arg:""`
}

func (m *TagPagination[T]) Template() ([]byte, error) {
	return nil, nil
}

func (m *TagPagination[T]) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{"Pagination": m.Pagination, "Tag": m.Tag, "Order": m.Order}, nil
}
