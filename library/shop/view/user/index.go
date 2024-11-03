package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type IndexRequest struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
}

type IndexReply struct {
	extends.ViewBase
	SiteData serviceargument.SiteData[*model.Content]
}

func (m *IndexRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &IndexReply{
		ViewBase: extends.ViewBase{},
	}
	reply.SiteData = service.GetSiteData[*model.Content](context, m.Organization.ID)
	return reply, nil
}
