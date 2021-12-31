package result

import (
	"encoding/json"
	"github.com/nbvghost/dandelion/library/context"
)

var _ Result = (*JsonResult)(nil)

type Result interface {
	Apply(context.IContext) ([]byte, error)
}

type JsonResult struct {
	Data interface{}
}

func (r *JsonResult) Apply(context context.IContext) ([]byte, error) {
	b, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}
	return b, nil
}
