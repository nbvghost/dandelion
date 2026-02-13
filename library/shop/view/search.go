package view

import (
	"fmt"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type SearchRequest struct {
	Organization *model.Organization      `mapping:""`
	Keyword      string                   `form:"keyword"`
	Type         model.FullTextSearchType `form:"type"`
	PageIndex    int                      `form:"page"`
}
type SearchReply struct {
	extends.ViewBase
	//Pagination module.Pagination[*model.FullTextSearch]
	Keyword string
	Type    string

	SiteData serviceargument.SiteData[*model.FullTextSearch]
}

func (m *SearchRequest) Render(ctx constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &SearchReply{}
	reply.SiteData = service.GetSiteData[*model.FullTextSearch](ctx, m.Organization.ID)

	var isBlur bool

	if strings.Contains(m.Keyword, "%") {
		isBlur = true
	}

	db := db.GetDB(ctx).Model(model.FullTextSearch{})

	switch m.Type {
	case model.FullTextSearchTypeContent:
		db.Where(`"Type"=?`, model.FullTextSearchTypeContent)
	case model.FullTextSearchTypeProducts:
		db.Where(`"Type"=?`, model.FullTextSearchTypeProducts)
	default:

	}
	if isBlur {
		kw := strings.ReplaceAll(m.Keyword, "%", "")
		db.Select(fmt.Sprintf(`"ID","CreatedAt","UpdatedAt","OID","TID","Uri","ContentItemID",ts_headline('english',"Title",to_tsquery('english','%s:*'),'HighlightAll=true') as "Title",ts_headline('english',"Content",to_tsquery('english','%s:*'),'MaxFragments=4') as "Content","Picture","Type"`, kw, kw))
		db.Where(fmt.Sprintf(`"Title" like '%s' or "Content" like '%s'`, m.Keyword, m.Keyword))
	} else {
		db.Select(fmt.Sprintf(`"ID","CreatedAt","UpdatedAt","OID","TID","Uri","ContentItemID",ts_headline('english',"Title",websearch_to_tsquery('english','%s'),'HighlightAll=true') as "Title",ts_headline('english',"Content",websearch_to_tsquery('english','%s'),'MaxFragments=4') as "Content","Picture","Type"`, m.Keyword, m.Keyword))
		db.Where(fmt.Sprintf(`"Index" @@ websearch_to_tsquery('english','%s')`, m.Keyword))
	}

	db.Where(`"OID"=?`, m.Organization.ID)

	var total int64
	db.Count(&total)

	var typeNameMap = make(map[dao.PrimaryKey]extends.Menus)
	listContentItem := repository.ContentItemDao.ListContentItemByOID(ctx, m.Organization.ID)
	for i := range listContentItem {
		item := listContentItem[i]
		typeNameMap[item.ID] = extends.NewMenusByContentItem(&item)
	}
	reply.SiteData.TypeNameMap = typeNameMap

	var pageSize = 20

	var list []*model.FullTextSearch
	err = db.Limit(pageSize).Offset(m.PageIndex * pageSize).Find(&list).Error
	if err != nil {
		return nil, err
	}

	reply.SiteData.Pagination = serviceargument.NewPagination(m.PageIndex, pageSize, int(total), list)

	reply.Keyword = m.Keyword
	reply.Type = string(m.Type)

	/*reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := m.ContentService.GetTitle(db.GetDB(ctx), m.Organization.ID)
		meta.SetBase(fmt.Sprintf("search results for %s", reply.Keyword), siteName, "search")
		return nil
	}*/
	return reply, nil
}
