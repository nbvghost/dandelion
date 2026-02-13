package page

import (
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type SubPageRequest struct {
	Organization  *model.Organization `mapping:""`
	ContentConfig model.ContentConfig `mapping:""`
	Sub           string              `uri:"Sub"`
	Name          string              `uri:"Name"`
}
type SubPageReply struct {
	extends.ViewBase

	SiteData serviceargument.SiteData[*model.Content]
}

func (m *SubPageRequest) Render(ctx constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &SubPageReply{}

	reply.SiteData = service.Site.GetContentTypeByUri(ctx, m.Organization.ID, "", "", 0)
	//reply.ContentData = m.ContentService.GetContentTypeByUri(context, m.Organization.ID, m.Uri, "", 0)
	//contentItem, contentSubType := m.ContentService.GetContentTypeByUri(m.Organization.ID, m.Uri, "")
	//reply.MenusData = module.NewMenusData(contentItem, contentSubType)
	reply.Name = fmt.Sprintf("%s/%s/%s", "page", m.Sub, m.Name)

	//reply.ContentConfig = m.ContentConfig
	/*menusData := m.ContentService.FindShowMenus(m.Organization.ID)
	for _, v := range menusData.List {
		if v.ID == contentItem.ID {
			reply.MenusData.Menus = v
			break
		}
	}*/

	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.GetDB(ctx), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", reply.SiteData.ContentItem.Name), siteName, "", "")
		return nil
	}
	return reply, nil
}
