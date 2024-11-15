package view

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type ContentsRequest struct {
	Organization *model.Organization `mapping:""`

	//TemplateName      string `uri:"TemplateName"`
	ContentItemUri    string `uri:"TypeUri"`
	ContentSubTypeUri string `uri:"SubTypeUri"`
	PageIndex         int    `form:"page"`
}
type ContentsReply struct {
	extends.ViewBase

	SiteData serviceargument.SiteData[*model.Content]

	SubTypeMap map[dao.PrimaryKey]string
}

func (m *ContentsRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &ContentsReply{
		SubTypeMap: map[dao.PrimaryKey]string{},
	}

	reply.SiteData = service.Site.GetContentTypeByUri(context, m.Organization.ID, m.ContentItemUri, m.ContentSubTypeUri, m.PageIndex)

	reply.Name = "contents/" + reply.SiteData.CurrentMenuData.Menus.TemplateName

	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", reply.SiteData.CurrentMenuData.Menus.Name), siteName, "", "")
		return nil
	}
	return reply, nil
}
