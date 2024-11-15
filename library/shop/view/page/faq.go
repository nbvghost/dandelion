package page

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
)

type FaqRequest struct {
	Organization *model.Organization `mapping:""`
}
type FaqReply struct {
	extends.ViewBase
}

func (m *FaqRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &FaqReply{
		ViewBase: extends.ViewBase{
			Name: "page/faq",
		},
	}
	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", "faq"), siteName, "fag", "fag")
		return nil
	}
	return reply, nil
}
