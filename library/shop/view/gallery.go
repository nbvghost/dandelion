package view

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type GalleryRequest struct {
	Organization *model.Organization `mapping:""`

	ContentItemUri    string `uri:"TypeID"`
	ContentSubTypeUri string `uri:"SubTypeID"`
	PageIndex         int    `form:"page"`
}
type GalleryReply struct {
	extends.ViewBase
	SiteData serviceargument.SiteData[*model.Content]
}

func (m *GalleryRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {

	reply := &GalleryReply{
		ViewBase: extends.ViewBase{
			Name: "gallery",
		},
	}

	reply.SiteData = service.Site.GetContentTypeByUri(context, m.Organization.ID, m.ContentItemUri, m.ContentSubTypeUri, m.PageIndex)

	//contentItem, contentSubType := m.ContentService.GetContentTypeByUri(m.Organization.ID, m.ContentItemUri, m.ContentSubTypeUri)

	//reply.MenusData = module.NewMenusData(contentItem, contentSubType)
	//menusData := m.ContentService.FindShowMenus(m.Organization.ID)
	/*for _, v := range menusData.List {
		if v.ID == reply.MenusData.TypeID {
			reply.MenusData.Menus = v
			break
		}
	}*/

	/*pageIndex, pageSize, total, list, err := m.ContentService.PaginationContent(m.Organization.ID, reply.MenusData.TypeID, reply.MenusData.SubTypeID, m.Page)
	if err != nil {
		return nil, err
	}*/

	//reply.Pagination = module.NewContentPagination(pageIndex, pageSize, total, list)

	/*reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := m.ContentService.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", reply.ContentData.CurrentMenuData.Menus.Name), siteName, "")
		return nil
	}*/
	return reply, nil
}
