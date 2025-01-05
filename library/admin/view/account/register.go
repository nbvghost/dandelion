package account

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
)

type Register struct {
}
type RegisterView struct {
	extends.ViewBase
}

func (m *Register) Render(ctx constrain.IContext) (constrain.IViewResult, error) {

	return &RegisterView{}, nil
}
