package account

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
)

type Login struct {
}
type LoginView struct {
	extends.ViewBase
}

func (m *Login) Render(ctx constrain.IContext) (constrain.IViewResult, error) {

	return &LoginView{}, nil
}
