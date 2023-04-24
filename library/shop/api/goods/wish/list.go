package wish

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
)

type List struct {
	User *model.User `mapping:""`
}

func (m *List) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	return nil, nil
}
