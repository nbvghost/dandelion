package blog

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type DetailRequest struct {
	Organization *model.Organization `mapping:""`
	ContentUri   string              `uri:"ContentID"`
}
type DetailReply struct {
	extends.ViewBase
	SiteData serviceargument.SiteData[*model.Content]
	//MenusData  module.MenusData
	//SubTypeMap map[dao.PrimaryKey]string
	//Content    *model.Content
	//LeftRight [2]model.Content
	//Tags      []extends.Tag
}

func (m *DetailRequest) Render(context constrain.IContext) (constrain.IViewResult, error) {
	reply := &DetailReply{
		ViewBase: extends.ViewBase{
			Name: "blog/detail",
		},
		//SubTypeMap: map[dao.PrimaryKey]string{},
	}
	c := repository.ContentDao.GetContentByUri(m.Organization.ID, m.ContentUri)
	//reply.Content = c

	contentItem := repository.ContentItemDao.GetContentItemByID(c.ContentItemID)
	contentSubType := repository.ContentSubTypeDao.GetContentSubTypeByID(c.ContentSubTypeID)

	reply.SiteData = service.Site.GetContentTypeByUri(context, m.Organization.ID, contentItem.Uri, contentSubType.Uri, 0)
	reply.SiteData.Item = c
	//reply.MenusData = module.NewMenusData(contentItem, contentSubType)

	//menusData := m.ContentService.FindShowMenus(m.Organization.ID)
	/*for _, v := range menusData.List {
		if v.ID == c.ContentItemID {
			reply.MenusData.Menus = v
			for _, sv := range v.List {
				reply.SubTypeMap[sv.ID] = sv.Name
				for _, ssv := range sv.List {
					reply.SubTypeMap[ssv.ID] = ssv.Name
				}
			}
			break
		}
	}
	reply.Tags = tag.ToTagsUri(c.Tags)*/
	//reply.LeftRight = m.ContentService.FindContentListForLeftRight(c.ContentItemID, c.ContentSubTypeID, c.ID, c.CreatedAt)
	/*reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
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
