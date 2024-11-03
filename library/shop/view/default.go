package view

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"

)

type DefaultRequest struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
}

type DefaultReply struct {
	extends.ViewBase
	ContentConfig *model.ContentConfig
	Organization  *model.Organization
}

func (m *DefaultRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &DefaultReply{
		ContentConfig: m.ContentConfig,
		Organization:  m.Organization,
	}

	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase("home", siteName, "", m.Organization.Introduction)
		photos := m.Organization.Photos
		if len(photos) > 0 {
			photo, err := oss.ReadUrl(context, photos[0])
			if err != nil {
				return err
			}
			meta.SetOGImage(photo, 0, 0, m.Organization.Introduction, "")
		}
		return nil
	}
	return reply, nil
}
