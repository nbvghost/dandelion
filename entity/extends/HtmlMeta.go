package extends

import (
	"strings"
	"time"
)

type HtmlMetaOGType string

const (
	HtmlMetaOGTypeArticle = "article"
	HtmlMetaOGTypeProduct = "product"
	HtmlMetaOGTypeWebsite = "website"
)

type HtmlMeta struct {
	OG struct {
		Type HtmlMetaOGType
	}
	Properties []map[string]interface{}

	//<meta property = "og:site_name" content = "Colby Fayock" />
	//<meta property = "og:title"    content = "Anyone Can Map! Inspiration and an	introduction to the world of mapping - Colby Fayock"    />
	//<meta property = "og:description"    content = "Chef Gusteau wld of…"    />
	//<meta property = "og:url"    content = "https://www.colbyfayock.com/2020apping/"    />
	//<meta property = "og:type" content = "article" />

	//<meta property = "article:publisher" content = "https://www.colbyfayock.com" />
	//<meta property = "article:section" content = "Coding" />
	//<meta property = "article:tag" content = "Coding" />
	//<meta property = "og:image"    content = "https://res.cloudinary.com/1.1"    />
	//<meta property = "og:image:secure_url"    content = "https://res.cloudi20of%20mapping/blog-social-card-1.1"/>
	//<meta property = "og:image:width" content = "1280" />
	//<meta property = "og:image:height" content = "640" />
	//<meta property = "twitter:card" content = "summary_large_image" />
	//<meta property = "twitter:image"    content = "https://res.cloudinary.com/fay/he%20world%20of%20mapping/blog-social-card-1.1"    />
	//<meta property = "twitter:site" content = "@colbyfayock" />

	//<meta property="og:title" content="Open Graph protocol">
	//<meta property="og:type" content="website">
	//<meta property="og:url" content="https://ogp.me/">
	//<meta property="og:image" content="https://ogp.me/logo.png">
	//<meta property="og:image:type" content="image/png">
	//<meta property="og:image:width" content="300">
	//<meta property="og:image:height" content="300">
	//<meta property="og:image:alt" content="The Open Graph logo">
	//<meta property="og:description" content="The Open Graph protocol enables any web page to become a rich object in a social graph.">
}

func (m *HtmlMeta) SetOGType(Type HtmlMetaOGType) *HtmlMeta {
	m.OG.Type = Type
	return m
}
func (m *HtmlMeta) SetOGImage(imgUrl string, width int, height int, alt string, contentType string) *HtmlMeta {
	property := map[string]interface{}{}
	property["og:image"] = imgUrl
	if len(imgUrl) >= 5 {
		if strings.Contains(strings.ToLower(imgUrl[:5]), "https") {
			property["og:image:secure_url"] = imgUrl
		}
	}
	if width > 0 {
		property["og:image:width"] = width
	}
	if height > 0 {
		property["og:image:height"] = height
	}
	if len(contentType) > 0 {
		property["og:image:type"] = contentType
	}
	if len(alt) > 0 {
		property["og:image:alt"] = alt
	}

	property["twitter:image"] = imgUrl

	m.Properties = append(m.Properties, property)
	return m
}
func (m *HtmlMeta) setLang(lang string) *HtmlMeta {
	property := map[string]interface{}{}
	property["og:locale"] = lang
	m.Properties = append(m.Properties, property)
	return m
}
func (m *HtmlMeta) setRequestURl(url string) *HtmlMeta {
	property := map[string]interface{}{}
	property["og:url"] = url
	m.Properties = append(m.Properties, property)
	return m
}
func (m *HtmlMeta) SetBase(title string, siteName string, keywords string, description string) *HtmlMeta {
	property := map[string]interface{}{}
	if len([]rune(description)) > 150 {
		description = string([]rune(description)[0:150])
	}
	property["description"] = description
	property["og:title"] = title + " | " + siteName
	property["og:description"] = description
	property["og:site_name"] = siteName
	property["twitter:card"] = "summary_large_image"
	property["twitter:title"] = title + " | " + siteName
	property["twitter:description"] = description

	property["keywords"] = keywords
	property["description"] = description

	m.Properties = append(m.Properties, property)

	//keywords    string `gorm:"column:Keywords"`    //
	//description string `gorm:"column:Description"` //

	return m
}

/*
article:published_time-日期时间 - 文章首次发表的时间。
article:modified_time-日期时间 - 文章最后一次更改的时间。
article:expiration_time- datetime - 当文章过期后。
article:author-配置文件 数组- 文章的作者。
article:section-字符串- 高级部分名称。例如技术
article:tag-字符串 数组- 与本文相关的标签词。
*/
func (m *HtmlMeta) SetArticle(section string, author string, publishedTime, modifiedTime time.Time, tags ...string) *HtmlMeta {
	m.OG.Type = HtmlMetaOGTypeWebsite
	property := map[string]interface{}{}
	property["article:published_time"] = publishedTime.Format("2006-01-02T15:04:05Z")
	property["article:modified_time"] = modifiedTime.Format("2006-01-02T15:04:05Z")
	property["article:author"] = author
	property["article:section"] = section
	for _, v := range tags {
		property["article:tag"] = v
	}
	m.Properties = append(m.Properties, property)
	return m
}
func (m *HtmlMeta) SetProduct(section string, publishedTime, modifiedTime time.Time, tags ...string) *HtmlMeta {
	m.OG.Type = HtmlMetaOGTypeProduct
	property := map[string]interface{}{}
	property["product:published_time"] = publishedTime.Format("2006-01-02T15:04:05Z")
	property["product:modified_time"] = modifiedTime.Format("2006-01-02T15:04:05Z")
	property["product:section"] = section
	for _, v := range tags {
		property["product:tag"] = v
	}
	m.Properties = append(m.Properties, property)
	return m
}
func NewHtmlMeta(lang, requestUrl string) *HtmlMeta {
	n := &HtmlMeta{}
	n.Properties = []map[string]interface{}{}
	n.OG.Type = HtmlMetaOGTypeWebsite
	n.setLang(lang).setRequestURl(requestUrl)
	return n
}
