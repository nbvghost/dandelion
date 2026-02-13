package widget

import (
	"context"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/repository"
)

type GalleryBlock struct {
	Organization *model.Organization `mapping:""`
}

func (g *GalleryBlock) Template() ([]byte, error) {
	return nil, nil
}

func (g *GalleryBlock) Render(ctx constrain.IContext) (map[string]any, error) {
	//org := ctx.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)

	outData := make(map[string]any)

	ContentItemList, ContentList := g.galleryBlock(ctx, g.Organization.ID, 8)
	outData["ContentItemList"] = ContentItemList
	outData["ContentList"] = ContentList

	return outData, nil
}
func (g *GalleryBlock) galleryBlock(ctx context.Context, OID dao.PrimaryKey, num int) ([]model.ContentItem, []model.Content) {
	contentItemList := repository.ContentItemDao.FindContentItemByType(ctx, model.ContentTypeGallery, OID)

	contentItemIDList := make([]dao.PrimaryKey, 0)
	for _, item := range contentItemList {
		contentItemIDList = append(contentItemIDList, item.ID)
	}
	contentList := repository.ContentDao.FindContentByIDAndNum(ctx, contentItemIDList, num)
	return contentItemList, contentList
}
