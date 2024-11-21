package content

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
	PageIndex    int                 `form:"page"`
	Order        string              `form:"order"`
}
type TagReply struct {
	extends.ViewBase

	Tag      extends.Tag
	SiteData serviceargument.SiteData[*model.Content]

	Order string
}

func (m *TagRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &TagReply{
		ViewBase: extends.ViewBase{
			Name: "content/tag",
		},
		//TypeMap:    map[dao.PrimaryKey]extends.Menus{},
		//SubTypeMap: map[dao.PrimaryKey]extends.Menus{},
	}

	reply.SiteData = service.GetSiteData[*model.Content](context, m.Organization.ID)

	reply.Tag = tag.ToTagsName([]extends.Tag{{Uri: m.Tag}})[0]

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
	case "like":
		orderMethod = append(orderMethod, dao.Sort{
			ColumnName: "CountLike",
			Method:     dao.OrderMethodDESC,
		})
	default:
		orderMethod = append(orderMethod, dao.Sort{
			ColumnName: "CountView",
			Method:     dao.OrderMethodDESC,
		})
		reply.Order = "trending"
	}

	tags, err := service.Content.FindContentTags(m.Organization.ID)
	if err != nil {
		return nil, err
	}
	reply.SiteData.Tags = tags

	pageIndex, pageSize, total, list, err := service.Content.FindContentByTag(m.Organization.ID, reply.Tag, m.PageIndex, orderMethod...)
	reply.SiteData.Pagination = serviceargument.NewPagination(pageIndex, pageSize, int(total), list)

	//listContentItem := m.ContentService.ListContentItemByOID(m.Organization.ID)
	//listContentSubType := m.ContentService.FindAllContentSubType(m.Organization.ID)

	/*for i := range listContentItem {
		item := listContentItem[i]
		reply.TypeMap[item.ID] = extends.NewMenusByContentItem(&item)
		for ii := range listContentSubType {
			subItem := listContentSubType[ii]
			if item.ID == subItem.ContentItemID {
				reply.SubTypeMap[subItem.ID] = extends.NewMenusByContentSubType(&item, &subItem)
			}
		}
	}*/

	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("%s", reply.Tag.Name), siteName, "", "")
		return nil
	}
	return reply, err
}
