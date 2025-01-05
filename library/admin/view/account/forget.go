package account

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
)

type Forget struct {
}
type ForgetView struct {
	extends.ViewBase
}

func (m *Forget) Render(ctx constrain.IContext) (constrain.IViewResult, error) {

	return &ForgetView{}, nil
}
