package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type SearchPagination struct {
	Pagination serviceargument.Pagination[*model.FullTextSearch] `arg:""`
	Keyword    string                                            `arg:""`
	Type       string                                            `arg:""`
}

func (m *SearchPagination) Template() ([]byte, error) {
	return nil, nil
}

func (m *SearchPagination) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{"Pagination": m.Pagination, "Keyword": m.Keyword, "Type": m.Type}, nil
}
