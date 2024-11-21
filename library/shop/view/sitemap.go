package view

import (
	"encoding/xml"
	"fmt"
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/tag"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/environments"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
	"strings"
	"time"
)

type SitemapRequest struct {
	Organization *model.Organization `mapping:""`
}
type SitemapReply struct {
	extends.ViewBase
	UrlSet UrlSet
}

func (m *SitemapReply) GetResult(context constrain.IContext, viewHandler constrain.IViewHandler) constrain.IResult {

	return &result.XMLResult{Data: m.UrlSet}
}

type Url struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"` //always,hourly,daily,weekly,	monthly,yearly,	never
	Priority   string `xml:"priority"`
}

//xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
//xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
//xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"

type UrlSet struct {
	XMLName        xml.Name `xml:"urlset"`
	Xsi            string   `xml:"xmlns:xsi,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	Urls           []Url    `xml:"url"`
}

func (m *SitemapRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	contextValue := contexext.FromContext(context)
	domainName := contextValue.DomainName
	if !environments.Release() {
		//domainName = fmt.Sprintf("dev.%s", domainName)
	}

	reply := &SitemapReply{}
	reply.UrlSet = UrlSet{
		Xsi:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
		Xmlns:          "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	var urls []Url

	scheme := util.GetScheme(contextValue.Request)

	var defaultLang = "en"
	var httpHost string
	if strings.EqualFold(contextValue.Lang, defaultLang) {
		httpHost = fmt.Sprintf("%s://%s", scheme, domainName)
	} else {
		httpHost = fmt.Sprintf("%s://%s.%s", scheme, contextValue.Lang, domainName)
	}

	var contentItemMap = map[dao.PrimaryKey]model.ContentItem{}
	{
		contentItemList := repository.ContentItemDao.ListContentItemByOID(m.Organization.ID)
		for i, v := range contentItemList {
			contentItemMap[v.ID] = contentItemList[i]

			urls = append(urls, Url{
				Loc:        fmt.Sprintf("%s/%s/%s", httpHost, v.Type, v.Uri),
				LastMod:    v.UpdatedAt.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.3",
			})
		}
	}
	var contentSubTypeMap = map[dao.PrimaryKey]model.ContentSubType{}
	{
		list := repository.ContentSubTypeDao.FindAllContentSubType(m.Organization.ID)
		for i, v := range list {
			contentSubTypeMap[v.ID] = list[i]
			var contentItem = contentItemMap[v.ContentItemID]

			urls = append(urls, Url{
				Loc:        fmt.Sprintf("%s/%s/%s/%s", httpHost, contentItem.Type, contentItem.Uri, v.Uri),
				LastMod:    v.UpdatedAt.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.3",
			})

		}
	}
	//var contentList []model.Content
	contentList := dao.Find(db.Orm(), &model.Content{}).Where(`"OID"=?`, m.Organization.ID).List()
	/*err = m.ContentService.FindAllByOID(db.Orm(), &contentList, m.Organization.ID)
	if err != nil {
		return nil, err
	}*/
	var contentTags pq.StringArray
	for i := range contentList {
		v := contentList[i].(*model.Content)

		var contentItem = contentItemMap[v.ContentItemID]
		var contentSubType = contentSubTypeMap[v.ContentSubTypeID]

		switch contentItem.Type {
		case model.ContentTypeContents:
			urls = append(urls, Url{
				Loc:        fmt.Sprintf("%s/content/detail/%s", httpHost, v.Uri),
				LastMod:    v.UpdatedAt.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.7",
			})
		case model.ContentTypeContent:
			urls = append(urls, Url{
				Loc:        fmt.Sprintf("%s/content/%s/%s", httpHost, contentItem.Uri, contentSubType.Uri),
				LastMod:    v.UpdatedAt.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.7",
			})
		case model.ContentTypeIndex:
			urls = append(urls, Url{
				Loc:        fmt.Sprintf("%s", httpHost),
				LastMod:    v.UpdatedAt.Format("2006-01-02"),
				ChangeFreq: "daily",
				Priority:   "1.0",
			})
		case model.ContentTypeGallery:
			//todo 相册暂时使用文章形式来显示其详细内容
			urls = append(urls, Url{
				Loc:        fmt.Sprintf("%s/content/detail/%s", httpHost, v.Uri),
				LastMod:    v.UpdatedAt.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.5",
			})
		case model.ContentTypeProducts:
			//产品走产品表
		case model.ContentTypeBlog:
			urls = append(urls, Url{
				Loc:        fmt.Sprintf("%s/blog/detail/%s", httpHost, v.Uri),
				LastMod:    v.UpdatedAt.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.7",
			})
		case model.ContentTypePage:
			//固定页面
		}
		contentTags = append(contentTags, v.Tags...)
	}

	//goods

	var goodsTags pq.StringArray

	/*var goodsList []model.Goods
	err = m.GoodsService.FindAllByOID(db.Orm(), &goodsList, m.Organization.ID)
	if err != nil {
		return nil, err
	}*/
	goodsTypeData := service.Goods.GoodsType.GetGoodsTypeData(m.Organization.ID)
	for i := range goodsTypeData.List {
		goodsType := goodsTypeData.List[i].Item
		urls = append(urls, Url{
			Loc:        fmt.Sprintf("%s/products/%s", httpHost, goodsType.Uri),
			LastMod:    goodsType.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.3",
		})
		for ii := range goodsTypeData.List[i].SubType {
			goodsSubType := goodsTypeData.List[i].SubType[ii].Item
			urls = append(urls, Url{
				Loc:        fmt.Sprintf("%s/products/%s/%s", httpHost, goodsType.Uri, goodsSubType.Uri),
				LastMod:    goodsSubType.UpdatedAt.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.3",
			})
		}
	}
	goodsList := dao.Find(db.Orm(), &model.Goods{}).Where(`"OID"=?`, m.Organization.ID).List()
	for i := range goodsList {
		v := goodsList[i].(*model.Goods)
		urls = append(urls, Url{
			Loc:        fmt.Sprintf("%s/product/detail/%d", httpHost, v.ID),
			LastMod:    v.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "1.0",
		})
		goodsTags = append(goodsTags, v.Tags...)
	}
	for _, v := range contentTags {
		urls = append(urls, Url{
			Loc:        fmt.Sprintf("%s/content/tag/%s", httpHost, tag.ToTagName(v).Uri),
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "1.0",
		})
	}
	for _, v := range goodsTags {
		urls = append(urls, Url{
			Loc:        fmt.Sprintf("%s/product/tag/%s", httpHost, tag.ToTagName(v).Uri),
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "1.0",
		})
	}
	reply.UrlSet.Urls = urls
	return reply, nil
}
