package wx

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
)

type Callback struct {
	Get struct {
		APPID string `uri:"APPID"`
	} `method:"Get"`
}

func (m *Callback) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return result.NewJsonResult(nil), nil
}
