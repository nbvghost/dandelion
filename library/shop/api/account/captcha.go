package account

import (
	"bytes"
	"github.com/nbvghost/captcha"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/tool/object"
	"strings"
	"time"
)

type Captcha struct {
}

func (m *Captcha) Handle(context constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(context)

	contextValue.Response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	contextValue.Response.Header().Set("Pragma", "no-cache")
	contextValue.Response.Header().Set("Expires", "0")
	contextValue.Response.Header().Set("Content-Type", "image/png")
	//contextValue.Response.Header().Set("Content-Type", "application/octet-stream")

	content := bytes.NewBuffer(nil)

	captchaDigits := captcha.RandomDigits(5)
	_, err := captcha.NewImage(captcha.New(), captchaDigits, captcha.StdWidth, captcha.StdHeight).WriteTo(content)
	if err != nil {
		return &result.NoneResult{}, nil
	}

	var captchaCodes []string
	for i := 0; i < len(captchaDigits); i++ {
		captchaCodes = append(captchaCodes, object.ParseString(captchaDigits[i]))
	}

	err = context.Redis().Set(context, "captcha:"+context.Token(), strings.Join(captchaCodes, ""), time.Minute*15)
	if err != nil {
		return nil, err
	}
	return &result.ImageBytesResult{
		Data:        content.Bytes(),
		ContentType: "image/png",
		Filename:    "",
	}, nil
}
