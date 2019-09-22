package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/glog"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Catch1688Service struct {
	dao.BaseDao
	Voucher      VoucherService
	Goods        GoodsService
	Organization OrganizationService
	Admin        AdminService
	URLS         []GGoodsType
}

type GGoodsType struct {
	GoodsType      string
	GoodsTypeChild string
	URL            string
}

func (service *Catch1688Service) Run() {

	service.URLS = make([]GGoodsType, 0)

	service.URLS = append(service.URLS, GGoodsType{GoodsType: "童装母婴玩具", GoodsTypeChild: "婴儿睡袋", URL: "https://detail.1688.com/offer/540490600398.html?sk=consign"})

	service.Catch()
}
func (service *Catch1688Service) GetInfo(content string) map[string]interface{} {

	ttttt := strings.ReplaceAll(content, "'", `"`)
	glog.Trace(ttttt)
	sdfdsfsd := make(map[string]interface{})
	json.Unmarshal([]byte(ttttt), &sdfdsfsd)
	glog.Trace(sdfdsfsd)
	//吴可国国国
	//s := &service.Catch1688Service{}
	//glog.Trace(s.GetInfo())
	//s.Run()

	return sdfdsfsd
}
func (service *Catch1688Service) getGoodsType(content string) map[string]interface{} {


	//reg := regexp.MustCompile(`iDetailData.registerRenderData\(([\s\S]*)[\)^;]*`)
	reg:= regexp.MustCompile(`<script.*?>([\s\S]*?)<\/script>`)

	resss:=reg.FindAllStringSubmatch(content, -1)
	glog.Trace(len(resss))

	for x:=0;x<len(resss);x++{
		glog.Trace(resss[x][1])

		
		/*reg := regexp.MustCompile(`iDetailData.registerRenderData\(({[\s\S]+})\);`)
		if reg.MatchString(resss[x][1]){
			ress := reg.FindAllStringSubmatch(resss[x][1],-1)
			fmt.Println(len(ress))
			fmt.Println(ress)
		}*/


	}




	ressss:=reg.FindAllString(content, -1)
	glog.Trace(len(ressss))




	dfd := make(map[string]interface{})

	tesd := resss[1][0]
	tesd = strings.ReplaceAll(tesd, "categoryId", `"categoryId"`)
	tesd = strings.ReplaceAll(tesd, "categoryName", `"categoryName"`)
	tesd = strings.ReplaceAll(tesd, "sellerInf", `"sellerInf"`)
	tesd = strings.ReplaceAll(tesd, "userId", `"userId"`)
	tesd = strings.ReplaceAll(tesd, "loginId", `"loginId"`)
	tesd = strings.ReplaceAll(tesd, "categoryList", `"categoryList"`)

	glog.Trace(tesd)

	json.Unmarshal([]byte(tesd), &dfd)
	glog.Trace(dfd)

	return dfd
}
func (service *Catch1688Service) Catch() {

	for i := 0; i < len(service.URLS); i++ {

		item := service.URLS[i]

		res, err := http.Get("https://detail.1688.com/offer/540490600398.html?sk=consign")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		//b, err := ioutil.ReadAll(res.Body)

		reader := transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())
		//b, err := ioutil.ReadAll(reader)

		//fmt.Println(err)
		//fmt.Println(string(b))
		//fmt.Println(utf8.DecodeLastRune(b))

		//doc, err := goquery.NewDocument("https://detail.1688.com/offer/540490600398.html?sk=consign")
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Fatal(err)
		}

		haveAdmin := service.Admin.FindAdminByAccount(dao.Orm(), "admin")

		gt, gtc := service.Goods.AddGoodsTypeByNameByChild(item.GoodsType, item.GoodsTypeChild)

		goods := &dao.Goods{}
		goods.GoodsTypeID = gt.ID
		goods.GoodsTypeChildID = gtc.ID
		goods.OID = haveAdmin.OID

		docHtml, err := doc.Html()
		fmt.Println(err)

		glog.Trace(service.getGoodsType(docHtml))

		doc.Find("h1.d-title").Each(func(i int, s *goquery.Selection) {
			title := s.Text()
			title = strings.ReplaceAll(title, "批发", "")

			fmt.Println(title)
			goods.Title = title
		})

		doc.Find(".obj-content .price-content>span.price-num").Each(func(i int, s *goquery.Selection) {

			fmt.Println(s.Text())
		})

		doc.Find("#mod-detail-attributes .obj-content table tbody tr td").Each(func(i int, s *goquery.Selection) {

			fmt.Println(s.Text())
		})
		totlStock := uint64(0)
		specifications := make([]dao.Specification, 0)
		doc.Find(".table-sku tbody tr").Each(func(i int, s *goquery.Selection) {

			spe := dao.Specification{}
			spe.Label, _ = s.Find("td.name span.image").Attr("title")
			spe.Num = 1
			spe.Weight = 250
			_stroc, _ := strconv.ParseUint(s.Find("td.count span em.value").Text(), 10, 64)
			spe.Stock = uint(_stroc)
			_price, _ := strconv.ParseUint(s.Find("td.price span em.value").Text(), 10, 64)

			spe.CostPrice = _price

			totlStock = totlStock + spe.CostPrice

			specifications = append(specifications, spe)

			fmt.Println(totlStock)
			fmt.Println("totlStock")
			fmt.Println(s.Find("td.name span.image").Attr("title"))
			fmt.Println(s.Find("td.price span em.value").Text())
			fmt.Println(s.Find("td.count span em.value").Text())
			//fmt.Println(s.Html())
		})

		fmt.Println(specifications)

		doc.Find("#desc-lazyload-container").Each(func(i int, s *goquery.Selection) {

			fmt.Printf(s.Attr("data-tfs-url"))
			u, ise := s.Attr("data-tfs-url")
			if ise {
				res, err := http.Get(u)
				fmt.Println(err)
				//b, err := ioutil.ReadAll(res.Body)
				reader := transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())
				b, err := ioutil.ReadAll(reader)
				fmt.Println(err)
				//fmt.Println(string(b))
				te := strings.TrimSpace(string(b))
				//fmt.Println(te)

				reg := regexp.MustCompile(`^var offer_details={"content":"(.+)"};`)
				ress := reg.FindStringSubmatch(te)
				fmt.Println(len(ress))
				fmt.Println(ress[1])

				red := strings.NewReader(ress[1])
				imgsDoc, err := goquery.NewDocumentFromReader(red)
				fmt.Println(err)

				imgsDoc.Find("img").Each(func(i int, s *goquery.Selection) {
					fmt.Println(s.Attr("src"))
				})

			}

		})

	}

}
