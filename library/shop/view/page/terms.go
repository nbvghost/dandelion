package page

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
)

type TermsRequest struct {
	Organization *model.Organization `mapping:""`
}
type TermsReply struct {
	extends.ViewBase
}

func (m *TermsRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &TermsReply{
		ViewBase: extends.ViewBase{
			Name: "page/terms",
		},
	}
	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", "terms & conditions"), siteName, "", "")
		return nil
	}
	return reply, nil
}
