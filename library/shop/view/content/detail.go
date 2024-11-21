package content

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type DetailRequest struct {
	Organization *model.Organization `mapping:""`
	ContentUri   string              `uri:"ContentUri"`
	//ContentID      dao.PrimaryKey      `uri:"ContentID"`
}
type DetailReply struct {
	extends.ViewBase
	//MenusData module.MenusData
	//Content   *model.Content
	//LeftRight [2]model.Content
	SiteData serviceargument.SiteData[*model.Content]
}

func (m *DetailRequest) Render(context constrain.IContext) (constrain.IViewResult, error) {
	reply := &DetailReply{}

	var c *model.Content

	/*if m.ContentID > 0 {
		c = m.ContentService.GetContentByID(m.ContentID)
	} else if len(m.ContentUri) > 0 {
		c = m.ContentService.GetContentByUri(m.Organization.ID, m.ContentUri)
	} else {
		return nil, errors.New("没有找到内容")
	}*/

	c = repository.ContentDao.GetContentByUri(m.Organization.ID, m.ContentUri)

	contentItem := repository.ContentItemDao.GetContentItemByID(c.ContentItemID)
	contentSubType := repository.ContentSubTypeDao.GetContentSubTypeByID(c.ContentSubTypeID)

	reply.Name = string(contentItem.Type) + "/detail/" + contentItem.TemplateName

	reply.SiteData = service.Site.GetContentTypeByUri(context, m.Organization.ID, contentItem.Uri, contentSubType.Uri, 0)
	reply.SiteData.Item = c
	if len(reply.SiteData.Pagination.List) > 0 {
		modelContent := reply.SiteData.Pagination.List[0]
		reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
			siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
			meta.SetBase(fmt.Sprintf("%s | %s", modelContent.Title, reply.SiteData.CurrentMenuData.Menus.Name), siteName, modelContent.Keywords, modelContent.Description)

			imgUrl, _ := oss.ReadUrl(context, modelContent.Picture)
			meta.SetOGImage(imgUrl, 0, 0, modelContent.Title, "")

			for _, v := range modelContent.Images {
				imgUrl, _ = oss.ReadUrl(context, v)
				meta.SetOGImage(imgUrl, 0, 0, modelContent.Title, "")
			}

			meta.SetArticle(reply.SiteData.CurrentMenuData.Menus.Name, modelContent.Author, modelContent.CreatedAt, modelContent.UpdatedAt, modelContent.Tags...)
			return nil
		}
	}

	/*reply.Content = c

	contentItem, contentSubType := m.ContentService.GetContentTypeByID(m.Organization.ID, c.ContentItemID, c.ContentSubTypeID)
	reply.MenusData = module.NewMenusData(contentItem, contentSubType)

	menusData := m.ContentService.FindShowMenus(m.Organization.ID)

	for _, v := range menusData.List {
		if v.ID == c.ContentItemID {
			reply.MenusData.Menus = v
			break
		}
	}
	reply.LeftRight = m.ContentService.FindContentListForLeftRight(c.ContentItemID, c.ContentSubTypeID, c.ID, c.CreatedAt)
	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := m.ContentService.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s | %s", reply.Content.Title, reply.MenusData.Menus.Name), siteName, reply.Content.Summary)
		imgUrl, err := ossurl.CreateUrl(context, reply.Content.Picture)
		if err != nil {
			return err
		}
		meta.SetOGImage(imgUrl, 0, 0, reply.Content.Title, "")
		meta.SetArticle(reply.MenusData.Menus.Name, reply.Content.Author, reply.Content.CreatedAt, reply.Content.UpdatedAt, reply.Content.Tags...)
		return nil
	}*/
	return reply, nil
}
