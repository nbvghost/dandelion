package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
)

type Get struct {
	Store *model.Store `mapping:""`
}

func (g *Get) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: g.Store}}, nil

}
