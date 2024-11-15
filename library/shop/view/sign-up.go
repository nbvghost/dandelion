package view

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"net/url"
	"strings"
)

type SignUpRequest struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
}

type SignUpReply struct {
	extends.ViewBase
	//MenusData     module.MenusData
	//ContentConfig *model.ContentConfig
	//Organization  *model.Organization
	SiteData serviceargument.SiteData[*model.Content]
	Redirect string
}

func (m *SignUpReply) GetResult(context constrain.IContext, viewHandler constrain.IViewHandler) constrain.IResult {
	if strings.Contains(m.Redirect, "/sign-in") || strings.Contains(m.Redirect, "/sign-up") {
		m.Redirect = ""
		return nil
	}
	{
		//如果已经带有redirect，则不要求在跳转
		contextValue := contexext.FromContext(context)
		redirect := contextValue.Request.URL.Query().Get("redirect")
		if len(redirect) > 0 {
			return nil
		}
	}
	if len(m.Redirect) == 0 {
		return nil
	}
	params := &url.Values{}
	params.Set("redirect", m.Redirect)
	return &result.RedirectToUrlResult{Url: "/sign-up?" + params.Encode()}
}

func (m *SignUpRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	contextValue := contexext.FromContext(context)
	reply := &SignUpReply{
		ViewBase: extends.ViewBase{},
		Redirect: "",
	}

	redirect := contextValue.Request.URL.Query().Get("redirect")
	if len(redirect) == 0 {
		redirect = contextValue.Request.Header.Get("Referer")
	}
	reply.Redirect = redirect
	contentItem := repository.ContentItemDao.GetContentItemOfIndex(db.Orm(), m.Organization.ID)
	reply.SiteData = service.Site.GetContentTypeByUri(context, m.Organization.ID, contentItem.Uri, "", 0)

	/*reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := m.ContentService.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(contentItem.Name, siteName, m.Organization.Introduction)
		photos := m.Organization.Photos
		if len(photos) > 0 {
			photo, err := ossurl.CreateUrl(context, photos[0])
			if err != nil {
				return err
			}
			meta.SetOGImage(photo, 0, 0, m.Organization.Introduction, "")
		}
		return nil
	}*/
	return reply, nil
}
