package file

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
)

type Captcha struct {
}

func (m *Captcha) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	buf := util.CreateCaptchaCodeBytes(string(play.SessionCaptcha))
	return &result.ImageBytesResult{Data: buf}, err
}
