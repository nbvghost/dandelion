package api

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
)

type Heartbeat struct{}

func (m *Heartbeat) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return &result.JsonResult{Data: result.NewSuccess("OK")}, err
}
