package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type TypeList struct {
}

func (m *TypeList) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: repository.ContentTypeDao.ListContentType(ctx)}}, nil
}
