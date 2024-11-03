package product

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/domain/tag"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"regexp"
)

type DetailRequest struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
	GoodsID       dao.PrimaryKey       `uri:"GoodsID"`
}
type DetailReply struct {
	extends.ViewBase
	GoodsInfo extends.GoodsMix
	LikeList  [][]model.Goods
	//MenusData module.MenusData
	SiteData serviceargument.SiteData[*extends.GoodsDetail]
}

func (m *DetailRequest) Render(context constrain.IContext) (constrain.IViewResult, error) {
	reply := &DetailReply{
		ViewBase: extends.ViewBase{
			Name: "product/detail",
		},
	}

	var err error
	var goodsInfo *extends.GoodsMix
	goodsInfo, err = service.Goods.Goods.GetGoods(db.Orm(), context, m.GoodsID)
	if err != nil {
		return nil, err
	}
	reply.GoodsInfo = *goodsInfo

	reply.SiteData = service.Site.GoodsDetail(context, m.Organization.ID, goodsInfo.GoodsType.Uri, goodsInfo.GoodsTypeChild.Uri)
	reply.SiteData.Tags = tag.ToTagsUri(goodsInfo.Goods.Tags)

	likeList := service.Goods.GoodsType.ListGoodsByType(m.Organization.ID, goodsInfo.Goods.GoodsTypeID, goodsInfo.Goods.GoodsTypeChildID)
	pi := 0
	for i := range likeList {
		if len(reply.LikeList) == 0 {
			reply.LikeList = append(reply.LikeList, []model.Goods{likeList[i]})
			continue
		}
		if len(reply.LikeList[pi]) < 4 {
			reply.LikeList[pi] = append(reply.LikeList[pi], likeList[i])
		} else {
			reply.LikeList = append(reply.LikeList, []model.Goods{likeList[i]})
			pi++
		}
	}
	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", reply.GoodsInfo.Goods.Title), siteName, s.ReplaceAllString(reply.GoodsInfo.Goods.Summary, ","), reply.GoodsInfo.Goods.Introduce)
		if len(reply.GoodsInfo.Goods.Images) > 0 {
			imgUrl, err := oss.ReadUrl(context, reply.GoodsInfo.Goods.Images[0])
			if err != nil {
				return err
			}
			meta.SetOGImage(imgUrl, 0, 0, reply.GoodsInfo.Goods.Title, "")
		}
		meta.SetProduct(reply.GoodsInfo.GoodsType.Name, reply.GoodsInfo.Goods.CreatedAt, reply.GoodsInfo.Goods.UpdatedAt, reply.GoodsInfo.Goods.Tags...)
		return nil
	}
	return reply, nil
}

var s = regexp.MustCompile(`\s+`)
