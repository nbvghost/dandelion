package view

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type IndexRequest struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
}

type IndexReply struct {
	extends.ViewBase
	//MenusData     module.MenusData
	//ContentConfig *model.ContentConfig
	//Organization  *model.Organization
	SiteData               serviceargument.SiteData[*model.Content]
	ShowAtHomeList         []*model.ContentItem
	StickyTopGoodsTypeList []model.GoodsType
}

func (m *IndexRequest) Render(ctx constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &IndexReply{
		ViewBase: extends.ViewBase{},
	}

	contentItem := repository.ContentItemDao.GetContentItemOfIndex(db.GetDB(ctx), m.Organization.ID)
	reply.SiteData = service.Site.GetContentTypeByUri(ctx, m.Organization.ID, contentItem.Uri, "", 0)

	reply.ShowAtHomeList = repository.ContentItemDao.FindContentItemByShowAtHome(ctx, m.Organization.ID)
	reply.StickyTopGoodsTypeList = service.Goods.GoodsType.StickyTopGoodsTypeList(db.GetDB(ctx), m.Organization.ID)

	/*reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := m.ContentService.GetTitle(db.GetDB(ctx), m.Organization.ID)
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
