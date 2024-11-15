package page

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type PageRequest struct {
	Organization  *model.Organization `mapping:""`
	ContentConfig model.ContentConfig `mapping:""`
	Uri           string              `uri:"Uri"`
	PageIndex     int                 `form:"page"`
}
type PageReply struct {
	extends.ViewBase
	SiteData   serviceargument.SiteData[*model.Content]
	Pagination serviceargument.Pagination[*model.ContentItem]
}
type Config struct {
	TemplateName string
	Type         string
}

func (m *PageRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &PageReply{}

	reply.SiteData = service.Site.GetContentTypeByUri(context, m.Organization.ID, m.Uri, "", 0)
	//reply.ContentData = m.ContentService.GetContentTypeByUri(context, m.Organization.ID, m.Uri, "", 0)
	//contentItem, contentSubType := m.ContentService.GetContentTypeByUri(m.Organization.ID, m.Uri, "")
	//reply.MenusData = module.NewMenusData(contentItem, contentSubType)
	reply.Name = fmt.Sprintf("%s/%s", "page", reply.SiteData.ContentItem.TemplateName)

	var c Config
	json.Unmarshal([]byte(reply.SiteData.ContentItem.Config), &c)
	total, list := repository.ContentItemDao.FindContentItemByTypeTemplate(m.Organization.ID, c.Type, c.TemplateName, m.PageIndex)
	pagination := serviceargument.NewPagination[*model.ContentItem](m.PageIndex, 20, int(total), list)
	reply.Pagination = pagination
	//<option value="contents">文章列表</option>
	//<option value="content">独立文章</option>
	//<option value="gallery">画廊</option>

	//reply.ContentConfig = m.ContentConfig
	/*menusData := m.ContentService.FindShowMenus(m.Organization.ID)
	for _, v := range menusData.List {
		if v.ID == contentItem.ID {
			reply.MenusData.Menus = v
			break
		}
	}*/

	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", reply.SiteData.ContentItem.Name), siteName, reply.SiteData.CurrentMenuData.Menus.Name, reply.SiteData.CurrentMenuData.Menus.Introduction)
		return nil
	}
	return reply, nil
}
