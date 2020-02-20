package service

import (
	"bytes"
	"encoding/json"
	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
	"html"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/glog"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Catch1688Service struct {
	dao.BaseDao
	Voucher         VoucherService
	Goods           GoodsService
	Organization    OrganizationService
	ExpressTemplate ExpressTemplateService
	Admin           AdminService
	//URLS         []string
}

func (service *Catch1688Service) Run() {

	//service.URLS = append(service.URLS, GGoodsType{URL:"https://detail.1688.com/offer/565805587556.html?sk=consign"})
	//service.URLS = append(service.URLS, GGoodsType{URL:"https://detail.1688.com/offer/594644186430.html?sk=consign"})
	//service.URLS = append(service.URLS, GGoodsType{URL:"https://detail.1688.com/offer/586995720260.html?sk=consign"})
	//service.URLS = append(service.URLS, GGoodsType{URL:"https://detail.1688.com/offer/588862561819.html?sk=consign"})
	//service.URLS = append(service.URLS, GGoodsType{URL:"https://detail.1688.com/offer/590544720505.html?sk=consign"})
	//service.URLS = append(service.URLS, GGoodsType{URL:"https://detail.1688.com/offer/589160608401.html?sk=consign"})
	//service.URLS = append(service.URLS, GGoodsType{URL:"https://detail.1688.com/offer/591799200224.html?sk=consign"})

	type URLModel struct {
		Catch []string
	}

	go func() {
		for {

			URLS := URLModel{}
			util.JSONToStruct(conf.JsonText, &URLS)

			//URLS = append(URLS, "https://detail.1688.com/offer/562482031336.html?sk=consign")
			for index := range URLS.Catch {
				service.URLCatch(URLS.Catch[index])
				time.Sleep(60 * time.Second)
			}

			time.Sleep(60 * time.Second)



		}
	}()

	go func() {

		for{
			list,err:=os.Open("1688")
			if glog.Error(err)==false{

				fl,err:=list.Readdir(-1)
				if glog.Error(err)==false{

					for i:=range fl{

						f,err:=os.Open("1688/"+fl[i].Name())
						if glog.Error(err)==false{

							b,err:=ioutil.ReadAll(f)
							if glog.Error(err)==false{
								service.Catch(string(b),fl[i].Name(),false)
							}
						}

					}

				}

			}

			time.Sleep(60 * time.Second)
		}


	}()
}

func (service *Catch1688Service) getGoodsType(content string) (map[string]interface{}, map[string]interface{}) {

	goodsType := make(map[string]interface{})
	goodsInfo := make(map[string]interface{})

	//reg := regexp.MustCompile(`iDetailData.registerRenderData\(([\s\S]*)[\)^;]*`)
	reg := regexp.MustCompile(`<script.*?>([\s\S]*?)<\/script>`)

	resss := reg.FindAllStringSubmatch(content, -1)
	glog.Trace(len(resss))

	for x := 0; x < len(resss); x++ {

		if strings.EqualFold(resss[x][1], "") == false {
			glog.Trace("----------------")

			reg := regexp.MustCompile(`\siDetailData.registerRenderData\(({[\s\S]+})\);\s+$`)
			if reg.MatchString(resss[x][1]) {
				//glog.Trace(resss[x][1])

				ress := reg.FindAllStringSubmatch(resss[x][1], -1)
				//fmt.Println(len(ress))
				//fmt.Println(ress[0][1])

				tesd := ress[0][1]
				tesd = strings.ReplaceAll(tesd, "categoryId", `"categoryId"`)
				tesd = strings.ReplaceAll(tesd, "categoryName", `"categoryName"`)
				tesd = strings.ReplaceAll(tesd, "sellerInf", `"sellerInf"`)
				tesd = strings.ReplaceAll(tesd, "userId", `"userId"`)
				tesd = strings.ReplaceAll(tesd, "loginId", `"loginId"`)
				tesd = strings.ReplaceAll(tesd, "categoryList", `"categoryList"`)

				//glog.Trace(tesd)

				//return goodsType,goodsInfo
				json.Unmarshal([]byte(tesd), &goodsType)
				glog.Trace(goodsType)

			}
			///------------------------info--------------
			reg = regexp.MustCompile(`\s"sku"\s:\s({[\s\S]+})\s+};\s+iDetailData`)

			if reg.MatchString(resss[x][1]) {
				ress := reg.FindAllStringSubmatch(resss[x][1], -1)

				glog.Trace(ress[0][1])

				err := json.Unmarshal([]byte(ress[0][1]), &goodsInfo)
				glog.Error(err)
				glog.Trace(goodsInfo)
			}

			glog.Trace("----------------")

		}

	}

	return goodsType, goodsInfo
}
func (service *Catch1688Service) Catch(CatchContent,Mark string,isGbk bool) {
	//b, err := ioutil.ReadAll(res.Body)

	addPriceRotia := 0.3
	brokerageRotia := 0.1
	expresstePrice := 500


	haveAdmin := service.Admin.FindAdminByAccount(dao.Orm(), "admin")
	//邮件模板
	express := service.ExpressTemplate.GetExpressTemplateByOID(haveAdmin.OID)
	if express.ID == 0 {
		glog.Trace("没有创建快递模板无法添加产品")
		return
	}

	//item := service.URLS[i]

	_havg := service.Goods.FindGoodsLikeMark(Mark)
	if _havg.ID > 0 {
		return
	}


	body:=bytes.NewBufferString(CatchContent)

	var doc *goquery.Document
	var err error
	if isGbk{
		reader := transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
		doc, err = goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Fatal(err)
		}
	}else{
		doc, err = goquery.NewDocumentFromReader(body)
		if err != nil {
			log.Fatal(err)
		}
	}



	goods := &dao.Goods{}

	doc.Find("h1.d-title").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		title = strings.ReplaceAll(title, "批发", "")
		//fmt.Println(title)
		goods.Title = title
	})

	_goods := service.Goods.FindGoodsByTitle(goods.Title)
	if _goods.ID != 0 {
		return
	}
	goods.Mark = Mark
	goods.OID = haveAdmin.OID
	goods.ExpressTemplateID = express.ID

	docHtml, err := doc.Html()
	glog.Error(err)

	goodsType, goodsInfo := service.getGoodsType(docHtml)
	glog.Trace(goodsType, goodsInfo)

	totlStock := uint64(0)
	minPrice := math.MaxFloat64

	specifications := make([]dao.Specification, 0)

	if goodsInfo["skuMap"] != nil {
		skuMap := goodsInfo["skuMap"].(map[string]interface{})

		priceRange := float64(0)
		if goodsInfo["priceRange"] != nil {
			priceRanges := goodsInfo["priceRange"].([]interface{})

			if priceRanges != nil {
				priceRange = priceRanges[0].([]interface{})[1].(float64)
			}
		}

		for key, value := range skuMap {
			//k:=key.(string)
			v := value.(map[string]interface{})

			_price := float64(0)
			if v["price"] == nil {
				_price = priceRange
			} else {
				_price, _ = strconv.ParseFloat(v["price"].(string), 64)
			}

			_canBookCount := v["canBookCount"].(float64)
			if _canBookCount < 100 {
				continue
			}
			_price = _price * 100

			costPrice := _price
			brokerage := _price * brokerageRotia
			_price = _price + (_price * addPriceRotia) + brokerage + float64(expresstePrice)

			minPrice = math.Min(minPrice, _price)

			specification := dao.Specification{}
			specification.Label = html.UnescapeString(key)
			specification.MarketPrice = uint64(_price)
			specification.Brokerage = uint64(brokerage)
			specification.CostPrice = uint64(costPrice)
			specification.Stock = uint(_canBookCount)
			specification.Num = 1
			specification.Weight = 200
			specifications = append(specifications, specification)

			totlStock = totlStock + uint64(specification.Stock)
		}
	} else {

		////price-text price-num
		priceNum := doc.Find("div.price-content span.price-text.price-num").Text()
		minPrice, _ = strconv.ParseFloat(priceNum, 64)
		minPrice = minPrice * 100
		totlStock = 99999

		costPrice := minPrice
		brokerage := minPrice * brokerageRotia
		minPrice = minPrice + (minPrice * addPriceRotia) + brokerage + float64(expresstePrice)

		minPrice = math.Min(minPrice, minPrice)

		specification := dao.Specification{}
		specification.Label = html.UnescapeString(goods.Title)
		specification.MarketPrice = uint64(minPrice)
		specification.Brokerage = uint64(brokerage)
		specification.CostPrice = uint64(costPrice)
		specification.Stock = uint(totlStock)
		specification.Num = 1
		specification.Weight = 200
		specifications = append(specifications, specification)

	}

	if len(specifications) == 0 {
		return
	}

	goods.Price = uint64(minPrice)
	goods.Stock = uint(totlStock)

	categoryList := goodsType["categoryList"].([]interface{})

	categoryA := categoryList[0].(map[string]interface{})
	categoryB := categoryList[1].(map[string]interface{})

	gt, gtc := service.Goods.AddGoodsTypeByNameByChild(categoryA["name"].(string), categoryB["name"].(string))
	goods.GoodsTypeID = gt.ID
	goods.GoodsTypeChildID = gtc.ID

	imageList, exist := doc.Find("div.mod-detail-version2018-gallery").Attr("data-gallery-image-list")
	if exist {
		imageLists := strings.Split(imageList, ",")
		glog.Trace(imageLists)

		images := make([]string, 0)
		for i := 0; i < len(imageLists); i++ {
			if strings.EqualFold(imageLists[i], "") == false {
				imgPath := tool.DownloadInternetImage(imageLists[i], "", "")
				glog.Trace(imgPath)
				if strings.EqualFold(imgPath, "") == false {
					images = append(images, "//"+conf.Config.Domain+"/file/load?path="+imgPath)
					time.Sleep(200 * time.Millisecond)
				}

			}

		}

		goods.Images = util.StructToJSON(images)
	}else {
		images := make([]string, 0)
		doc.Find("div.mod-detail-gallery li.tab-trigger").Each(func(i int, selection *goquery.Selection) {

			imgPath,exist:=selection.Attr("data-imgs")
			if exist{
				//glog.Trace(imgPath)
				imgPathMap:=make(map[string]interface{})
				util.JSONToStruct(imgPath,&imgPathMap)
				if imgPathMap["original"]!=nil{
					imgPath := tool.DownloadInternetImage(imgPathMap["original"].(string), "", "")
					images = append(images, "//"+conf.Config.Domain+"/file/load?path="+imgPath)
					time.Sleep(200 * time.Millisecond)
				}
			}

		})
		goods.Images = util.StructToJSON(images)

	}

	/*doc.Find(".obj-content .price-content>span.price-num").Each(func(i int, s *goquery.Selection) {

		fmt.Println(s.Text())
	})*/

	//[{"Name":"gkhjkhg","Value":"jkghj"}]
	attributes := make([]map[string]interface{}, 0)
	ai := 0
	ak := ""
	doc.Find("#mod-detail-attributes .obj-content table tbody tr td").Each(func(i int, s *goquery.Selection) {
		ai++
		if ai%2 == 0 {
			attributes = append(attributes, map[string]interface{}{"Name": ak, "Value": s.Text()})
		} else {
			ak = s.Text()
		}

	})
	//fmt.Println(attributes)

	goods.Params = util.StructToJSON(attributes)

	/*doc.Find(".table-sku tbody tr").Each(func(i int, s *goquery.Selection) {

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
	})*/

	//fmt.Println(specifications)

	doc.Find("#desc-lazyload-container").Each(func(i int, s *goquery.Selection) {

		//fmt.Printf(s.Attr("data-tfs-url"))
		u, ise := s.Attr("data-tfs-url")
		if ise {
			res, err := http.Get(u)
			glog.Error(err)
			//b, err := ioutil.ReadAll(res.Body)
			reader := transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())
			b, err := ioutil.ReadAll(reader)
			glog.Error(err)
			//fmt.Println(string(b))
			te := strings.TrimSpace(string(b))
			//fmt.Println(te)

			reg := regexp.MustCompile(`^var offer_details={"content":"(.+)"};`)
			ress := reg.FindStringSubmatch(te)
			//fmt.Println(len(ress))
			//fmt.Println(ress[1])

			red := strings.NewReader(ress[1])
			imgsDoc, err := goquery.NewDocumentFromReader(red)
			glog.Error(err)

			images := make([]string, 0)
			imgsDoc.Find("img").Each(func(i int, s *goquery.Selection) {
				pimgUrl, exits := s.Attr("src")
				if exits && strings.EqualFold(pimgUrl, "") == false {
					pimgUrl = strings.ReplaceAll(pimgUrl, `\`, "")
					pimgUrl = strings.ReplaceAll(pimgUrl, `"`, "")
					imgPath := tool.DownloadInternetImage(pimgUrl, "", "")

					if strings.EqualFold(imgPath, "") == false {
						images = append(images, "//"+conf.Config.Domain+"/file/load?path="+imgPath)
						time.Sleep(200 * time.Millisecond)
					}

				}
			})

			goods.Pictures = util.StructToJSON(images)

		}

	})

	goods.Videos = "[]"

	service.Goods.Add(dao.Orm(), goods)
	for s := 0; s < len(specifications); s++ {
		specifications[s].GoodsID = goods.ID
		service.Goods.Add(dao.Orm(), &(specifications[s]))
	}
}
func (service *Catch1688Service) URLCatch(URL string) {

	res, err := http.Get(URL)
	if err != nil {
		glog.Trace(err)
		time.Sleep(time.Hour * 3)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	b,err:=ioutil.ReadAll(res.Body)
	glog.Trace(err)
	service.Catch(string(b),URL,true)


}
