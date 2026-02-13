package view

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type BlogRequest struct {
	Organization *model.Organization `mapping:""`

	ContentItemUri    string `uri:"TypeID"`
	ContentSubTypeUri string `uri:"SubTypeID"`
	PageIndex         int    `form:"page"`
}
type BlogReply struct {
	extends.ViewBase
	SiteData     serviceargument.SiteData[*model.Content]
	SubTypeMap   map[dao.PrimaryKey]extends.Menus
	TypeID       dao.PrimaryKey
	SubTypeID    dao.PrimaryKey
	SubSubTypeID dao.PrimaryKey
}

func (m *BlogRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &BlogReply{
		ViewBase: extends.ViewBase{
			Name: "blog",
		},
		SubTypeMap: map[dao.PrimaryKey]extends.Menus{},
	}

	reply.SiteData = service.Site.GetContentTypeByUri(context, m.Organization.ID, m.ContentItemUri, m.ContentSubTypeUri, 0)

	//contentItem, contentSubType := m.ContentService.GetContentTypeByUri(m.Organization.ID, m.ContentItemUri, m.ContentSubTypeUri)
	//reply.MenusData = module.NewMenusData(contentItem, contentSubType)
	//menusData := m.ContentService.FindShowMenus(m.Organization.ID)
	/*for _, v := range menusData.List {
		if v.ID == reply.MenusData.TypeID {
			reply.MenusData.Menus = v
			for _, sv := range v.List {
				reply.SubTypeMap[sv.ID] = sv
				for _, ssv := range sv.List {
					reply.SubTypeMap[ssv.ID] = ssv
				}
			}
			break
		}
	}

	reply.TypeID = contentItem.ID
	if contentSubType.ParentContentSubTypeID == 0 {
		reply.SubTypeID = contentSubType.ID
		reply.SubSubTypeID = 0
	} else {
		reply.SubTypeID = contentSubType.ParentContentSubTypeID
		reply.SubSubTypeID = contentSubType.ID
	}*/

	/*reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := m.ContentService.GetTitle(db.GetDB(ctx), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", reply.ContentData.CurrentMenuData.Menus.Name), siteName, "")
		return nil
	}*/
	return reply, nil
}
