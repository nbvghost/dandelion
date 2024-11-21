package product

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/tag"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type TagRequest struct {
	Organization *model.Organization `mapping:""`
	Tag          string              `uri:"Tag"`
	PageIndex    int                 `uri:"PageIndex"`
	Order        string              `form:"order"`
}
type TagReply struct {
	extends.ViewBase
	Tag extends.Tag
	//Pagination module.Pagination[*model.Goods]
	Order string
	//TypeMap    map[dao.PrimaryKey]extends.Menus
	//Tags       []extends.Tag
	TagContent *model.Content

	SiteData serviceargument.SiteData[*model.Goods]
}

func (m *TagRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &TagReply{
		ViewBase: extends.ViewBase{
			Name: "product/tag",
		},
		//TypeMap: map[dao.PrimaryKey]extends.Menus{},
	}

	/*menusData := m.ContentService.FindShowMenus(m.Organization.ID)
	for _, v := range menusData.List {
		if v.Type == model.ContentTypeProducts {
			reply.MenusData.Menus = v
			break
		}
	}*/
	reply.SiteData = service.GetSiteData[*model.Goods](context, m.Organization.ID)

	ts := tag.ToTagsName([]extends.Tag{{Uri: m.Tag}})
	reply.Tag = ts[0]

	reply.Order = m.Order

	var orderMethod []dao.Sort
	switch m.Order {
	case "trending":
		orderMethod = append(orderMethod, dao.Sort{
			ColumnName: "CountView",
			Method:     dao.OrderMethodDESC,
		})
	case "latest":
		orderMethod = append(orderMethod, dao.Sort{
			ColumnName: "CreatedAt",
			Method:     dao.OrderMethodDESC,
		})
	case "sale":
		orderMethod = append(orderMethod, dao.Sort{
			ColumnName: "CountSale",
			Method:     dao.OrderMethodDESC,
		})
	default:
		orderMethod = append(orderMethod, dao.Sort{
			ColumnName: "CountView",
			Method:     dao.OrderMethodDESC,
		})
		reply.Order = "trending"
	}
	pageIndex, pageSize, total, list, err := service.Goods.Tag.FindGoodsByTag(m.Organization.ID, reply.Tag, m.PageIndex, orderMethod...)
	reply.SiteData.Pagination = serviceargument.NewPagination(pageIndex, pageSize, int(total), list)

	reply.TagContent = service.Content.GetByTitle(db.Orm(), m.Organization.ID, reply.Tag.Name)

	/*listContentItem := m.GoodsTypeService.ListGoodsByOID(m.Organization.ID)
	for _, v := range listContentItem {
		reply.TypeMap[v.ID] = extends.Menus{
			ID:           v.ID,
			Name:         v.Name,
			TemplateName: "product",
			Type:         model.ContentTypeProducts,
		}
	}*/

	tags, err := service.Goods.Tag.FindGoodsTags(m.Organization.ID)
	if err != nil {
		return nil, err
	}
	reply.SiteData.Tags = tags
	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", reply.Tag.Name), siteName, "", "")
		return nil
	}
	return reply, err
}
