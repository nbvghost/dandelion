package session

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Index struct {
	Get struct{} `method:"get"`
}

func (m *Index) Handle(context constrain.IContext) (constrain.IResult, error) {
	var err error

	var cartCount uint
	if context.UID() > 0 {
		cartCount, err = service.Order.ShoppingCart.FindShoppingCartListCount(context.UID())
		if err != nil {
			return nil, err
		}
	}
	return result.NewData(map[string]any{
		"CartCount": cartCount,
	}), nil
}
