package view

import (
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type ContentRequest struct {
	Organization *model.Organization `mapping:""`

	ContentItemUri    string `uri:"TypeUri"`
	ContentSubTypeUri string `uri:"SubTypeUri"`
}
type ContentReply struct {
	extends.ViewBase
	SiteData serviceargument.SiteData[*model.Content]
}

func (m *ContentRequest) Render(ctx constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &ContentReply{
		ViewBase: extends.ViewBase{
			Name: "content",
		},
	}

	reply.SiteData = service.Site.GetContentTypeByUri(ctx, m.Organization.ID, m.ContentItemUri, m.ContentSubTypeUri, 0)

	if len(reply.SiteData.Pagination.List) > 0 {
		modelContent := reply.SiteData.Pagination.List[0]
		reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
			siteName := service.Content.GetTitle(db.GetDB(ctx), m.Organization.ID)
			meta.SetBase(fmt.Sprintf("%s | %s", modelContent.Title, reply.SiteData.CurrentMenuData.Menus.Name), siteName, modelContent.Keywords, modelContent.Description)
			imgUrl, _ := oss.ReadUrl(ctx, modelContent.Picture)
			meta.SetOGImage(imgUrl, 0, 0, modelContent.Title, "")

			for _, v := range modelContent.Images {
				imgUrl, _ = oss.ReadUrl(ctx, v)
				meta.SetOGImage(imgUrl, 0, 0, modelContent.Title, "")
			}

			meta.SetArticle(reply.SiteData.CurrentMenuData.Menus.Name, modelContent.Author, modelContent.CreatedAt, modelContent.UpdatedAt, modelContent.Tags...)
			return nil
		}
	}

	/*reply.MenusData = module.NewMenusData(contentItem, contentSubType)

	menusData := m.ContentService.FindShowMenus(m.Organization.ID)

	for _, v := range menusData.List {
		if v.ID == contentItem.ID {
			reply.MenusData.Menus = v
			break
		}
	}*/

	/*c := m.ContentService.GetContentByContentItemIDAndContentSubTypeID(contentItem.ID, contentSubType.ID)
	if c.IsZero() && len(reply.MenusData.Menus.List) > 0 {
		for _, v := range reply.MenusData.Menus.List {
			if v.ID == contentSubType.ID {
				if len(v.List) > 0 {
					c = m.ContentService.GetContentByContentItemIDAndContentSubTypeID(contentItem.ID, v.List[0].ID)
					break
				}
			}
		}
	}*/

	/*var modelContent = reply.ContentData.Content()

	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := m.ContentService.GetTitle(db.GetDB(ctx), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s | %s", modelContent.Title, reply.ContentData.MenusData.Menus.Name), siteName, modelContent.Summary)
		imgUrl, err := ossurl.CreateUrl(context, modelContent.Picture)
		if err != nil {
			return err
		}
		meta.SetOGImage(imgUrl, 0, 0, modelContent.Title, "")
		meta.SetArticle(reply.ContentData.MenusData.Menus.Name, modelContent.Author, modelContent.CreatedAt, modelContent.UpdatedAt, modelContent.Tags...)
		return nil
	}*/
	return reply, nil
}
