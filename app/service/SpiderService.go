package service

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/app/service/dao"

	"github.com/PuerkitoBio/goquery"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/tool"
	"golang.org/x/net/html"
)

type SpiderService struct {
	dao.BaseDao
	File    FileService
	Content ContentService
	Admin   AdminService
	OID     uint64
}

func init() {
	//go SpiderService{}.StartSpider()
	//SpiderService{}.ReadWeiXinArticle("https://mp.weixin.qq.com/s/Z5s02hxVJ2MbFnYbMLGGkw")
	//SpiderService{}.ReadWeiXinArticle("https://mp.weixin.qq.com/s?src=11&timestamp=1532599519&ver=1022&signature=s8FEg9-9SjADeW1PmUmWrGS8yCY1dBBFZ8Jh1Zhx06BBkfPM21KkPvZPQyYEt2i7AS-UALDxRvq-SAS9T68EYkM5bXHDdb3YI91I3s8Cn6wXcQ27wI*XMtuj8ulx*d5O&new=1")
}
func (spider SpiderService) StartSpider() {

	admin := spider.Admin.FindAdminByAccount(dao.Orm(), "admin")
	spider.OID = admin.OID
	//美女
	urlList := [][]string{
		{"http://weixin.sogou.com/weixin?type=2&s_from=input&query=%E7%BE%8E%E5%A5%B3", "美女"},
		{"http://weixin.sogou.com/weixin?type=2&s_from=input&query=%E6%80%A7%E6%84%9F", "美女"},
		{"http://weixin.sogou.com/weixin?type=2&ie=utf8&query=%E6%9D%A8%E5%B9%82%E6%80%A7%E6%84%9F%E8%A7%86%E9%A2%91&s_from=bottom_hint", "美女"},
		{"http://weixin.sogou.com/weixin?type=2&ie=utf8&query=%E6%80%A7%E6%84%9F%E7%B2%89%E8%89%B2%E5%A5%B6%E7%BD%A9&s_from=up_hint", "美女"},
	}
	go spider.WeixinQuerySogou(urlList)

	if true {
		//return
	}

	//搞笑=pc_1
	urls := []string{"pc_1/pc_1.html", "pc_1/1.html", "pc_1/2.html", "pc_1/3.html", "pc_1/4.html", "pc_1/5.html", "pc_1/6.html", "pc_1/7.html", "pc_1/8.html", "pc_1/9.html", "pc_1/10.html", "pc_1/11.html", "pc_1/12.html", "pc_1/13.html", "pc_1/14.html"}
	go spider.WeixinSogou(urls, "搞笑")

	//养生堂=pc_2
	urls = []string{"pc_2/pc_2.html", "pc_2/1.html", "pc_2/2.html", "pc_2/3.html", "pc_2/4.html", "pc_2/5.html", "pc_2/6.html", "pc_2/7.html", "pc_2/8.html", "pc_2/9.html", "pc_2/10.html", "pc_2/11.html", "pc_2/12.html", "pc_2/13.html", "pc_2/14.html"}
	go spider.WeixinSogou(urls, "养生堂")

	//私房话=pc_3
	urls = []string{"pc_3/pc_3.html", "pc_3/1.html", "pc_3/2.html", "pc_3/3.html", "pc_3/4.html", "pc_3/5.html", "pc_3/6.html", "pc_3/7.html", "pc_3/8.html", "pc_3/9.html", "pc_3/10.html", "pc_3/11.html", "pc_3/12.html", "pc_3/13.html", "pc_3/14.html"}
	go spider.WeixinSogou(urls, "私房话")

	//八卦精=pc_4
	urls = []string{"pc_4/pc_4.html", "pc_4/1.html", "pc_4/2.html", "pc_4/3.html", "pc_4/4.html", "pc_4/5.html", "pc_4/6.html", "pc_4/7.html", "pc_4/8.html", "pc_4/9.html", "pc_4/10.html", "pc_4/11.html", "pc_4/12.html", "pc_4/13.html", "pc_4/14.html"}
	go spider.WeixinSogou(urls, "八卦精")

	//科技咖=pc_5
	urls = []string{"pc_5/pc_5.html", "pc_5/1.html", "pc_5/2.html", "pc_5/3.html", "pc_5/4.html", "pc_5/5.html", "pc_5/6.html", "pc_5/7.html", "pc_5/8.html", "pc_5/9.html", "pc_5/10.html", "pc_5/11.html", "pc_5/12.html", "pc_5/13.html", "pc_5/14.html"}
	go spider.WeixinSogou(urls, "科技咖")

	//财经迷=pc_6
	//urls=[]string{"pc_6/pc_6.html","pc_6/1.html","pc_6/2.html","pc_6/3.html","pc_6/4.html","pc_6/5.html","pc_6/6.html","pc_6/7.html","pc_6/8.html","pc_6/9.html","pc_6/10.html","pc_6/11.html","pc_6/12.html","pc_6/13.html","pc_6/14.html"}
	//go spider.WeixinSogou(urls,"财经迷")

	//汽车控=pc_7
	//urls=[]string{"pc_7/pc_7.html","pc_7/1.html","pc_7/2.html","pc_7/3.html","pc_7/4.html","pc_7/5.html","pc_7/6.html","pc_7/7.html","pc_7/8.html","pc_7/9.html","pc_7/10.html","pc_7/11.html","pc_7/12.html","pc_7/13.html","pc_7/14.html"}
	//go spider.WeixinSogou(urls,"汽车控")

	//生活家=pc_8
	urls = []string{"pc_8/pc_8.html", "pc_8/1.html", "pc_8/2.html", "pc_8/3.html", "pc_8/4.html", "pc_8/5.html", "pc_8/6.html", "pc_8/7.html", "pc_8/8.html", "pc_8/9.html", "pc_8/10.html", "pc_8/11.html", "pc_8/12.html", "pc_8/13.html", "pc_8/14.html"}
	go spider.WeixinSogou(urls, "生活家")

	//时尚圈=pc_9
	//urls=[]string{"pc_9/pc_9.html","pc_9/1.html","pc_9/2.html","pc_9/3.html","pc_9/4.html","pc_9/5.html","pc_9/6.html","pc_9/7.html","pc_9/8.html","pc_9/9.html","pc_9/10.html","pc_9/11.html","pc_9/12.html","pc_9/13.html","pc_9/14.html"}
	//go spider.WeixinSogou(urls,"时尚圈")

	//育儿=pc_10
	urls = []string{"pc_10/pc_10.html", "pc_10/1.html", "pc_10/2.html", "pc_10/3.html", "pc_10/4.html", "pc_10/5.html", "pc_10/6.html", "pc_10/7.html", "pc_10/8.html", "pc_10/9.html", "pc_10/10.html", "pc_10/11.html", "pc_10/12.html", "pc_10/13.html", "pc_10/14.html"}
	go spider.WeixinSogou(urls, "育儿")

	//旅游=pc_11
	urls = []string{"pc_11/pc_11.html", "pc_11/1.html", "pc_11/2.html", "pc_11/3.html", "pc_11/4.html", "pc_11/5.html", "pc_11/6.html", "pc_11/7.html", "pc_11/8.html", "pc_11/9.html", "pc_11/10.html", "pc_11/11.html", "pc_11/12.html", "pc_11/13.html", "pc_11/14.html"}
	go spider.WeixinSogou(urls, "旅游")

	//职场=pc_12
	urls = []string{"pc_12/pc_12.html", "pc_12/1.html", "pc_12/2.html", "pc_12/3.html", "pc_12/4.html", "pc_12/5.html", "pc_12/6.html", "pc_12/7.html", "pc_12/8.html", "pc_12/9.html", "pc_12/10.html", "pc_12/11.html", "pc_12/12.html", "pc_12/13.html", "pc_12/14.html"}
	go spider.WeixinSogou(urls, "职场")

	//美食=pc_13
	urls = []string{"pc_13/pc_13.html", "pc_13/1.html", "pc_13/2.html", "pc_13/3.html", "pc_13/4.html", "pc_13/5.html", "pc_13/6.html", "pc_13/7.html", "pc_13/8.html", "pc_13/9.html", "pc_13/10.html", "pc_13/11.html", "pc_13/12.html", "pc_13/13.html", "pc_13/14.html"}
	go spider.WeixinSogou(urls, "美食")

	//历史=pc_14
	//urls=[]string{"pc_14/pc_14.html","pc_14/1.html","pc_14/2.html","pc_14/3.html","pc_14/4.html","pc_14/5.html","pc_14/6.html","pc_14/7.html","pc_14/8.html","pc_14/9.html","pc_14/10.html","pc_14/11.html","pc_14/12.html","pc_14/13.html","pc_14/14.html"}
	//go spider.WeixinSogou(urls,"历史")

	//教育=pc_15
	urls = []string{"pc_15/pc_15.html", "pc_15/1.html", "pc_15/2.html", "pc_15/3.html", "pc_15/4.html", "pc_15/5.html", "pc_15/6.html", "pc_15/7.html", "pc_15/8.html", "pc_15/9.html", "pc_15/10.html", "pc_15/11.html", "pc_15/12.html", "pc_15/13.html", "pc_15/14.html"}
	go spider.WeixinSogou(urls, "教育")

	//星座=pc_16
	//urls=[]string{"pc_16/pc_16.html","pc_16/1.html","pc_16/2.html","pc_16/3.html","pc_16/4.html","pc_16/5.html","pc_16/6.html","pc_16/7.html","pc_16/8.html","pc_16/9.html","pc_16/10.html","pc_16/11.html","pc_16/12.html","pc_16/13.html","pc_16/14.html"}
	//go spider.WeixinSogou(urls,"星座")

	//体育=pc_17
	//urls=[]string{"pc_17/pc_17.html","pc_17/1.html","pc_17/2.html","pc_17/3.html","pc_17/4.html","pc_17/5.html","pc_17/6.html","pc_17/7.html","pc_17/8.html","pc_17/9.html","pc_17/10.html","pc_17/11.html","pc_17/12.html","pc_17/13.html","pc_17/14.html"}
	//go spider.WeixinSogou(urls,"体育")

	//军事=pc_18
	//urls=[]string{"pc_18/pc_18.html","pc_18/1.html","pc_18/2.html","pc_18/3.html","pc_18/4.html","pc_18/5.html","pc_18/6.html","pc_18/7.html","pc_18/8.html","pc_18/9.html","pc_18/10.html","pc_18/11.html","pc_18/12.html","pc_18/13.html","pc_18/14.html"}
	//go spider.WeixinSogou(urls,"军事")

	//游戏=pc_19
	//urls=[]string{"pc_19/pc_19.html","pc_19/1.html","pc_19/2.html","pc_19/3.html","pc_19/4.html","pc_19/5.html","pc_19/6.html","pc_19/7.html","pc_19/8.html","pc_19/9.html","pc_19/10.html","pc_19/11.html","pc_19/12.html","pc_19/13.html","pc_19/14.html"}
	//go spider.WeixinSogou(urls,"游戏")

	//萌宠=pc_20
	//urls=[]string{"pc_20/pc_20.html","pc_20/1.html","pc_20/2.html","pc_20/3.html","pc_20/4.html","pc_20/5.html","pc_20/6.html","pc_20/7.html","pc_20/8.html","pc_20/9.html","pc_20/10.html","pc_20/11.html","pc_20/12.html","pc_20/13.html","pc_20/14.html"}
	//go spider.WeixinSogou(urls,"萌宠")

}
func (spider SpiderService) WeixinQuerySogou(urls [][]string) {

	for {

		for _, value := range urls {

			client := http.Client{}
			request, err := http.NewRequest("GET", value[0], nil)

			//request.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
			/*request.Header.Set("Accept-Encoding","gzip, deflate")
			request.Header.Set("Accept-Language","zh-CN,zh;q=0.9")
			request.Header.Set("Cache-Control","no-cache")
			request.Header.Set("Connection","keep-alive")
			request.Header.Set("Host","weixin.sogou.com")
			request.Header.Set("Pragma","no-cache")
			request.Header.Set("Upgrade-Insecure-Requests","1")
			request.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")*/

			glog.Error(err)
			response, err := client.Do(request)
			glog.Error(err)

			/*reader, err := gzip.NewReader(response.Body)

			b,err:=ioutil.ReadAll(reader)
			fmt.Println(err)
			fmt.Println(string(b))*/

			if err == nil {
				spider.GetArticleQueryDataAndAdd(response.Body, value[1])
			}

			time.Sleep(1 * time.Minute)
		}

		time.Sleep(60 * time.Minute)
	}

}
func (spider SpiderService) GetArticleQueryDataAndAdd(body io.ReadCloser, ContentSubTypeName string) {

	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(doc.Html())

	doc.Find("#main .news-box .news-list li").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		//ss,ee:=s.Html()

		//fmt.Printf("Review %s\n",ss)

		//fmt.Println(s.Find(".txt-box .img-d a img").Nodes[0])

		imgs, _ := s.Find(".img-box a img").Attr("src")
		//imgsurl := spider.File.DownNetWriteAliyunOSS("http:" + imgs)
		imgsurl := "/file/load?path=" + spider.File.DownNetImage("http:"+imgs)
		//fmt.Println(ds)
		//fmt.Println(imgs)

		info := s.Find(".txt-box h3 a")
		var link string
		if info != nil {
			link, _ = info.Attr("href")
		}

		//fmt.Println(ds)
		//fmt.Println(link)
		content := spider.ReadWeiXinArticle(link)
		if !strings.EqualFold(content, "") {

			_title := info.Text()
			title := strings.Replace(_title, " ", "", -1)
			//fmt.Println(title)

			des := s.Find(".txt-box p.txt-info")
			desTxt := strings.Replace(des.Text(), " ", "", -1)

			account := s.Find(".txt-box div.s-p")
			auth := account.Find("a.account")
			//fmt.Println(auth.Text())

			var createTime time.Time
			timedate, df := account.Attr("t")
			if df {
				ts, _ := strconv.ParseInt(timedate, 10, 64)
				createTime = time.Unix(ts, 0)
			}
			//fmt.Println(timedate)
			//fmt.Println(df)
			spider.Content.AddSpiderArticle(spider.OID, "热点文摘", ContentSubTypeName, auth.Text(), title, link, desTxt, imgsurl, content, createTime)

		}

	})

}
func (spider SpiderService) GetArticleDataAndAdd(body io.ReadCloser, ContentSubTypeName string) {

	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		//ss,ee:=s.Html()

		//fmt.Printf("Review %s\n",ss)
		//fmt.Printf("Review %s\n",ee)

		imgs, _ := s.Find(".img-box a img").Attr("src")
		//imgsurl := spider.File.DownNetWriteAliyunOSS("http:" + imgs)
		imgsurl := "/file/load?path=" + spider.File.DownNetImage("http:"+imgs)
		//fmt.Println(ds)
		//fmt.Println(imgs)

		info := s.Find(".txt-box h3 a")
		var link string
		if info != nil {
			link, _ = info.Attr("href")
		}

		//fmt.Println(ds)
		//fmt.Println(link)
		content := spider.ReadWeiXinArticle(link)
		if !strings.EqualFold(content, "") {

			_title := info.Text()
			title := strings.Replace(_title, " ", "", -1)
			//fmt.Println(title)

			des := s.Find(".txt-box p.txt-info")
			desTxt := strings.Replace(des.Text(), " ", "", -1)

			account := s.Find(".txt-box div.s-p")
			auth := account.Find("a.account")
			//fmt.Println(auth.Text())

			var createTime time.Time
			timedate, df := account.Find("span").Attr("t")
			if df {
				ts, _ := strconv.ParseInt(timedate, 10, 64)
				createTime = time.Unix(ts, 0)
			}
			//fmt.Println(timedate)
			//fmt.Println(df)
			spider.Content.AddSpiderArticle(spider.OID, "热点文摘", ContentSubTypeName, auth.Text(), title, link, desTxt, imgsurl, content, createTime)

		}

	})

}

//http://weixin.sogou.com/
func (spider SpiderService) WeixinSogou(urls []string, ContentSubTypeName string) {

	/*for _,value:=range urls{
		client:=http.Client{}
		request,err:=http.NewRequest("GET","http://weixin.sogou.com/pcindex/pc/"+value,nil)
		//Referer: http://weixin.sogou.com/
		//User-Agent: Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36
		request.Header.Set("Referer","http://weixin.sogou.com/")
		request.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
		request.Header.Set("X-Requested-With","XMLHttpRequest")
		glog.Error(err)
		response,err:=client.Do(request)
		glog.Error(err)
		if err==nil{
			spider.GetArticleDataAndAdd(response.Body,ContentSubTypeName)
		}
	}*/

	for {

		client := http.Client{}
		request, err := http.NewRequest("GET", "http://weixin.sogou.com/pcindex/pc/"+urls[0], nil)
		request.Header.Set("Referer", "http://weixin.sogou.com/")
		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
		request.Header.Set("X-Requested-With", "XMLHttpRequest")
		request.Header.Set("Content-Type", "text/html;utf-8")
		glog.Error(err)
		response, err := client.Do(request)
		glog.Error(err)
		if err == nil {
			spider.GetArticleDataAndAdd(response.Body, ContentSubTypeName)
		}
		time.Sleep(30 * time.Minute)
	}
}
func (spider SpiderService) ReadWeiXinArticle(url string) string {

	err, _, b := tool.RequestByHeader(url, "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1", "")

	//doc, err := goquery.NewDocument(url)
	if err != nil {
		glog.Error(err)
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(b))
	if err != nil {
		glog.Error(err)
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
				//value.Attr[index].Val = spider.File.DownNetWriteAliyunOSS(vatt.Val)
				value.Attr[index].Val = "/file/load?path=" + spider.File.DownNetImage(vatt.Val)
				att := html.Attribute{}
				att.Key = "src"
				att.Val = value.Attr[index].Val
				value.Attr = append(value.Attr, att)
				break
			}

		}

	}

	videos := content.Find(".video_iframe").Nodes
	for _, value := range videos {

		var isw = 0
		var ish = 0
		for index, vatt := range value.Attr {

			if strings.EqualFold(vatt.Key, "data-src") {
				//fmt.Println(vatt.Val)
				//vatt.Val = DownNetImage(vatt.Val)
				value.Attr[index].Key = "src"
				value.Attr[index].Val = vatt.Val
				/*att := html.Attribute{}
				att.Key = "src"
				att.Val = value.Attr[index].Val
				value.Attr = append(value.Attr, att)*/
			}
			if strings.EqualFold(vatt.Key, "width") {
				value.Attr[index].Val = "500"
				isw = 1
			}
			if strings.EqualFold(vatt.Key, "height") {
				value.Attr[index].Val = "375"
				ish = 1
			}
			if strings.EqualFold(vatt.Key, "frameborder") {
				value.Attr[index].Val = "0"
			}
			if strings.EqualFold(vatt.Key, "allow") {
				value.Attr[index].Val = "autoplay;fullscreen"
			}
			if strings.EqualFold(vatt.Key, "allowfullscreen") {
				value.Attr[index].Val = "true"
			}

		}
		if ish == 0 {
			att := html.Attribute{}
			att.Key = "height"
			att.Val = "375"
			value.Attr = append(value.Attr, att)
		}
		if isw == 0 {
			att := html.Attribute{}
			att.Key = "width"
			att.Val = "500"
			value.Attr = append(value.Attr, att)
		}
	}

	_html, err := content.Html()
	glog.Error(err)
	//fmt.Println(strings.TrimSpace(_html))
	return strings.TrimSpace(_html)
}
