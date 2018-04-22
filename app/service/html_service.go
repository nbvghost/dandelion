package service

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nbvghost/gweb/tool"
	"golang.org/x/net/html"
)

type HtmlService struct{}

func (self HtmlService) ReadWeiXinArticle(url string) string {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return ""
	}
	content := doc.Find("#js_content")
	imgs := content.Find("img").Nodes
	for _, value := range imgs {

		for index, vatt := range value.Attr {

			if strings.EqualFold(vatt.Key, "data-src") {
				//fmt.Println(vatt.Val)
				//vatt.Val = DownNetImage(vatt.Val)
				value.Attr[index].Key = "data-src"
				value.Attr[index].Val = "/file/load?path=" + File.DownNetImage(vatt.Val)

				att := html.Attribute{}
				att.Key = "src"
				att.Val = value.Attr[index].Val
				value.Attr = append(value.Attr, att)
				break
			}

		}

	}

	_html, err := content.Html()
	tool.CheckError(err)
	//fmt.Println(strings.TrimSpace(_html))
	return strings.TrimSpace(_html)

}
