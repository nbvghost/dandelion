package account

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/shop/embed"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"github.com/nbvghost/tool"
	"github.com/nbvghost/tool/object"
	"go.uber.org/zap"
	"html/template"
	"strings"
	"time"
)

type ResetPassword struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		Email   string
		Captcha string
	} `method:"Post"`
}

func (m *ResetPassword) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
func (m *ResetPassword) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	captcha, _ := ctx.Redis().Get(ctx, "captcha:"+ctx.Token())
	if !strings.EqualFold(captcha, m.Post.Captcha) {
		return nil, errors.New("incorrect CAPTCHA")
	}

	var err error

	u := service.User.GetByEmail(db.Orm(), m.Organization.ID, m.Post.Email)

	if u.IsZero() == false {
		err = m.sendEmail(ctx, u)
		if err != nil {
			ctx.Logger().With(zap.NamedError("Send-Email", err))
		}
	}
	ctx.Redis().Del(ctx, "captcha:"+ctx.Token())
	return result.NewData(map[string]any{"Email": m.Post.Email}), nil
}
func (m *ResetPassword) sendEmail(ctx constrain.IContext, u *model.User) error {
	contextValue := contexext.FromContext(ctx)
	var t *template.Template
	var err error
	t, err = template.New("").Parse(embed.RestPasswordEmailTemplate)
	if err != nil {
		return err
	}

	token := tool.UUID()
	err = ctx.Redis().Set(ctx, token, object.ParseString(u.ID), time.Minute*60)
	if err != nil {
		return err
	}

	host := util.GetHost(contextValue.Request)

	ossHost, err := ctx.GetDNSName(key.MicroServerOSS)
	if err != nil {
		return err
	}

	siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)

	tBuf := bytes.NewBuffer(nil)
	err = t.Execute(tBuf, map[string]any{
		"NickName":        u.Name,
		"RestPasswordURL": host + "/reset-password?token=" + token,
		"ContactURL":      host + "/page/contact-us",
		"SiteName":        siteName,
		"ToEmail":         u.Email,
		"LogoURL":         fmt.Sprintf("%s://%s/assets/%s/content/logo", util.GetScheme(contextValue.Request), ossHost, contextValue.DomainName),
	})
	if err != nil {
		return err
	}
	err = service.Network.Email.SendEmailTLS(m.Organization.ID, serviceargument.EmailContent{
		Subject:    "Reset your password",
		SenderName: siteName,
		ToEmails:   []string{u.Email},
		Content:    tBuf.String(),
	})
	if err != nil {
		return err
	}
	return nil
}
