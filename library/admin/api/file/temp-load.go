package file

import (
	"github.com/nbvghost/dandelion/constrain"
)

type TempLoad struct {
}

func (m *TempLoad) HandleGet(context constrain.IContext) (r constrain.IResult, err error) {
	//return gweb.FileTempLoadAction(context), err
	return nil, err
}

func (m *TempLoad) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}
