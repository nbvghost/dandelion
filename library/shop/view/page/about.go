package page

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
)

type AboutRequest struct {
	Organization *model.Organization `mapping:""`
}
type AboutReply struct {
	extends.ViewBase
	Organization model.Organization
}

func (m *AboutRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &AboutReply{
		ViewBase: extends.ViewBase{
			Name: "page/about",
		},
	}
	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", "about"), siteName, "about", m.Organization.Introduction)
		return nil
	}
	reply.Organization = *m.Organization
	return reply, nil
}
