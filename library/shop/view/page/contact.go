package page

import (
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
)

type ContactRequest struct {
	Organization *model.Organization `mapping:""`
}
type ContactReply struct {
	extends.ViewBase
	Organization  *model.Organization
	ContentConfig model.ContentConfig
}

func (m *ContactRequest) Render(ctx constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &ContactReply{
		ViewBase: extends.ViewBase{
			Name: "page/contact",
		},
	}

	reply.Organization = m.Organization
	reply.ContentConfig = repository.ContentConfigDao.GetContentConfig(db.GetDB(ctx), m.Organization.ID)
	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.GetDB(ctx), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", "contact us"), siteName, "contact", fmt.Sprintf("E-mail:%s,Phone:%s", m.Organization.Email, m.Organization.Telephone))
		return nil
	}
	return reply, nil
}
