package funcmap

import (
	"github.com/nbvghost/dandelion/library/context"
)

type IWidget interface {
	Render(ctx context.IContext) (map[string]interface{}, error)
}
