package page

import (
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
)

type PrivacyRequest struct {
	Organization *model.Organization `mapping:""`
}
type PrivacyReply struct {
	extends.ViewBase
}

func (m *PrivacyRequest) Render(ctx constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &PrivacyReply{
		ViewBase: extends.ViewBase{
			Name: "page/privacy",
		},
	}
	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.GetDB(ctx), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", "privacy & policy"), siteName, "", "")
		return nil
	}
	return reply, nil
}
