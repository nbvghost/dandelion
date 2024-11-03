package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/sqltype"
)

type LeaveMessageStyleType string

const (
	LeaveMessageStyleTypeCard   LeaveMessageStyleType = "card"
	LeaveMessageStyleTypeForm   LeaveMessageStyleType = "form"
	LeaveMessageStyleTypeButton LeaveMessageStyleType = "button"
)

type LeaveMessage struct {
	StyleType LeaveMessageStyleType `arg:""`
}

func (m *LeaveMessage) Render(ctx constrain.IContext) (map[string]any, error) {
	return map[string]any{
		"SocialTypeList": []sqltype.SocialType{
			sqltype.SocialTypeOther,
			sqltype.SocialTypeFacebook,
			sqltype.SocialTypeTwitter,
			sqltype.SocialTypeLinkedin,
			sqltype.SocialTypeInstagram,
			sqltype.SocialTypeSkype,
			sqltype.SocialTypeWechat,
			sqltype.SocialTypeQq,
			sqltype.SocialTypeEmail,
			sqltype.SocialTypePhone,
			sqltype.SocialTypeWhatsapp,
		},
		"StyleType": m.StyleType,
	}, nil
}
func (m *LeaveMessage) Template() ([]byte, error) {
	return LeaveMessageTemplate, nil
}
