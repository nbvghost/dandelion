package catch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/activity"
	"github.com/nbvghost/dandelion/service/internal/admin"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/express"
	"github.com/nbvghost/dandelion/service/internal/goods"
	"html"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/util"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

var vm = otto.New()

type Catch1688Service struct {
	model.BaseDao
	Voucher          activity.VoucherService
	Goods            goods.GoodsService
	GoodsTypeService goods.GoodsTypeService
	Organization     company.OrganizationService
	ExpressTemplate  express.ExpressTemplateService
	Admin            admin.AdminService
	//URLS         []string
}

var scriptRegexp = regexp.MustCompile(`<script.*?>([\s\S]*?)<\/script>`)
var titleRegexp = regexp.MustCompile(`<title>([\s\S]*?)</title>`)

type NameValue struct {
	Name  string
	Value string
}

func (m *Catch1688Service) Api(catchDir string) error {
	http.HandleFunc("/push-html", func(writer http.ResponseWriter, request *http.Request) {

		//writer.Header().Set("Access-Control-Allow-Origin", strings.TrimRight(request.Referer(), "/"))
		writer.Header().Set("Access-Control-Allow-Origin", "https://detail.1688.com")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,FromURL")        //todo
		writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT,OPTIONS,GET") //todo
		writer.Header().Set("Access-Control-Allow-Credentials", "true")                    //todo
		if request.Method == http.MethodOptions {
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(404)
			writer.Write([]byte(err.Error()))
		} else {
			resss := titleRegexp.FindStringSubmatch(string(body))
			if len(resss) < 2 {
				writer.WriteHeader(404)
				writer.Write([]byte("读取标题错误"))
				return
			}
			title := resss[1]

			os.MkdirAll(fmt.Sprintf("%s/%s", catchDir, title), os.ModePerm)
			err := os.WriteFile(fmt.Sprintf("%s/%s/%s", catchDir, title, "content.html"), body, os.ModePerm)
			if err != nil {
				writer.WriteHeader(404)
				writer.Write([]byte(err.Error()))
				return
			}
			err = os.WriteFile(fmt.Sprintf("%s/%s/%s", catchDir, title, "url.txt"), []byte(request.Header.Get("FromURL")), os.ModePerm)
			if err != nil {
				writer.WriteHeader(404)
				writer.Write([]byte(err.Error()))
				return
			}
			log.Println(len(resss))
		}

	})
	err := http.ListenAndServe(":8080", http.DefaultServeMux)
	if err != nil {
		return err
	}
	return nil
}
func (m *Catch1688Service) readGoods(dir string) error {
	//log.Println(dir.Name(),subDirs[i].Name(),"content.html")
	contentHtml, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, "content.html"))
	if err != nil {
		return err
	}
	//log.Println(contentHtml,err)

	document, err := goquery.NewDocumentFromReader(bytes.NewReader(contentHtml))
	if err != nil {
		return err
	}

	images := make([]string, 0)
	document.Find(".content-detail").Find(".desc-img-loaded").Each(func(i int, selection *goquery.Selection) {
		if len(selection.Nodes) > 0 {
			if selection.Nodes[0].Parent.Data != "a" {
				if value, has := selection.Attr("data-lazyload-src"); has {
					images = append(images, value)
				}
			}
		}
	})
	document.Find(".content-detail").Find(".desc-img-no-load").Each(func(i int, selection *goquery.Selection) {
		if len(selection.Nodes) > 0 {
			if selection.Nodes[0].Parent.Data != "a" {
				if value, has := selection.Attr("data-lazyload-src"); has {
					images = append(images, value)
				}
			}
		}
	})

	for imageIndex, imageSrc := range images {
		_, fileName := filepath.Split(imageSrc)
		saveFile := fmt.Sprintf("%s/image/%d-%s", dir, 1000+imageIndex, strings.Split(fileName, "?")[0])
		os.MkdirAll(fmt.Sprintf("%s/image", dir), os.ModePerm)
		if fi, err := os.Stat(saveFile); err == nil && fi.Size() > 0 {
			continue
		}
		log.Println(imageSrc)
		response, err := http.Get(imageSrc)
		if err != nil {
			log.Println(err)
			continue
		}
		defer response.Body.Close()
		imgBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			continue
		}
		err = os.WriteFile(saveFile, imgBody, os.ModePerm)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	func() {
		//offer-attr-item
		nameValueList := make([]NameValue, 0)
		document.Find(".offer-attr-item").Each(func(i int, selection *goquery.Selection) {
			nameValueList = append(nameValueList, NameValue{
				Name:  selection.Find(".offer-attr-item-name").Text(),
				Value: selection.Find(".offer-attr-item-value").Text(),
			})
		})
		jb, err := json.Marshal(&nameValueList)
		if err != nil {
			log.Println(err)
			return
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", dir, "attr.json"), jb, os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	resss := scriptRegexp.FindAllStringSubmatch(string(contentHtml), -1)
	for _, i2 := range resss {
		if len(i2) >= 2 && strings.Contains(i2[1], "window.__INIT_DATA") && strings.Contains(i2[1], "window.__GLOBAL_DADA") {
			storeDataObj, err := vm.Run("var window = {};\n" + i2[1] + "\n;window.__INIT_DATA;")
			if err != nil {
				return err
			}
			obj := storeDataObj.Object()
			{
				globalData, err := obj.Get("globalData")
				if err != nil {
					return err
				}

				func() {
					//globalData.orderParamModel.orderParam.skuParam
					orderParamModel, err := globalData.Object().Get("orderParamModel")
					if err != nil {
						log.Println(err)
						return
					}
					orderParam, err := orderParamModel.Object().Get("orderParam")
					if err != nil {
						log.Println(err)
						return
					}
					skuParam, err := orderParam.Object().Get("skuParam")
					if err != nil {
						log.Println(err)
						return
					}
					if jb, err := skuParam.MarshalJSON(); err != nil {
						log.Println(err)
					} else {
						err := os.WriteFile(fmt.Sprintf("%s/%s", dir, "skuParam.json"), jb, os.ModePerm)
						if err != nil {
							log.Println(err)
							return
						}
					}
				}()

				skuModel, err := globalData.Object().Get("skuModel")
				if err != nil {
					return err
				}

				func() {
					skuInfoMap, err := skuModel.Object().Get("skuInfoMap")
					if err != nil {
						log.Println(err)
						return
					}
					if jb, err := skuInfoMap.MarshalJSON(); err != nil {
						log.Println(err)
					} else {
						err := os.WriteFile(fmt.Sprintf("%s/%s", dir, "skuInfoMap.json"), jb, os.ModePerm)
						if err != nil {
							log.Println(err)
							return
						}
					}
				}()

				func() {
					skuProps, err := skuModel.Object().Get("skuProps")
					if err != nil {
						log.Println(err)
						return
					}
					if jb, err := skuProps.MarshalJSON(); err != nil {
						log.Println(err)
					} else {
						err := os.WriteFile(fmt.Sprintf("%s/%s", dir, "skuProps.json"), jb, os.ModePerm)
						if err != nil {
							log.Println(err)
							return
						}
					}
				}()

				images, err := globalData.Object().Get("images")
				if err != nil {
					return err
				}

				imagesObject := images.Object()
				keys := imagesObject.Keys()
				for keyIndex, key := range keys {
					if imageItem, err := imagesObject.Get(key); err == nil {
						if imageSrc, err := imageItem.Object().Get("fullPathImageURI"); err == nil {
							_, fileName := filepath.Split(imageSrc.String())
							saveFile := fmt.Sprintf("%s/head/%d-%s", dir, 1000+keyIndex, strings.Split(fileName, "?")[0])
							_ = os.MkdirAll(fmt.Sprintf("%s/head", dir), os.ModePerm)
							if fi, err := os.Stat(saveFile); err == nil && fi.Size() > 0 {
								continue
							}
							log.Println(imageSrc.String())
							response, err := http.Get(imageSrc.String())
							if err != nil {
								log.Println(err)
								continue
							}
							defer response.Body.Close()
							imgBody, err := io.ReadAll(response.Body)
							if err != nil {
								log.Println(err)
								continue
							}
							err = os.WriteFile(saveFile, imgBody, os.ModePerm)
							if err != nil {
								return err
							}
						}
					}
				}

				/*if imagesJSON, err := images.MarshalJSON(); err == nil {
					log.Println(images)
					log.Println(string(imagesJSON))
				}*/
			}
		}
		/*if len(i2) >= 2 && strings.Contains(i2[1], "window.__STORE_DATA") && strings.Contains(i2[1], "window.isOnline") && false {
			//storeData := make(map[string]any)
			storeDataObj, err := vm.Run("var window = {};\n" + i2[1] + "\n;window.__STORE_DATA;")
			if err != nil {
				return err
			}
			obj := storeDataObj.Object()
			{
				globalData, err := obj.Get("globalData")
				if err != nil {
					return err
				}
				images, err := globalData.Object().Get("images")
				if err != nil {
					return err
				}
				if _, err := images.MarshalJSON(); err == nil {
					//log.Println(string(imagesJSON))
				}
			}
			//log.Println(storeData)
		}*/
	}

	return nil
}
func (m *Catch1688Service) Run(catchDir string, translate func(query string) (string, error)) error {
	dir, err := os.Open(catchDir)
	if err != nil {
		return err
	}
	typeDirs, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	for i := 0; i < len(typeDirs); i++ {
		subTypeDir, err := os.Open(catchDir + "/" + typeDirs[i].Name())
		if err != nil {
			return err
		}
		subTypeDirs, err := subTypeDir.Readdir(-1)
		if err != nil {
			return err
		}

		for ii := 0; ii < len(subTypeDirs); ii++ {
			goodsDir, err := os.Open(fmt.Sprintf("%s/%s", subTypeDir.Name(), subTypeDirs[ii].Name()))
			if err != nil {
				return err
			}
			goodsDirs, err := goodsDir.Readdir(-1)
			if err != nil {
				return err
			}

			for iii := 0; iii < len(goodsDirs); iii++ {
				if _, err := os.Stat(goodsDir.Name() + "/" + goodsDirs[iii].Name() + "/name.txt"); err != nil {
					title, err := translate(goodsDirs[iii].Name())
					if err != nil {
						return err
					}
					err = os.WriteFile(goodsDir.Name()+"/"+goodsDirs[iii].Name()+"/name.txt", []byte(title), os.ModePerm)
					if err != nil {
						return err
					}
				}

				err = m.readGoods(fmt.Sprintf("%s/%s/%s", subTypeDir.Name(), subTypeDirs[ii].Name(), goodsDirs[iii].Name()))
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func (m *Catch1688Service) getGoodsType(content string) (map[string]interface{}, map[string]interface{}) {

	goodsType := make(map[string]interface{})
	goodsInfo := make(map[string]interface{})

	//reg := regexp.MustCompile(`iDetailData.registerRenderData\(([\s\S]*)[\)^;]*`)

	resss := scriptRegexp.FindAllStringSubmatch(content, -1)
	log.Println(len(resss))

	for x := 0; x < len(resss); x++ {

		if strings.EqualFold(resss[x][1], "") == false {
			log.Println("----------------")

			reg := regexp.MustCompile(`\siDetailData.registerRenderData\(({[\s\S]+})\);\s+$`)
			if reg.MatchString(resss[x][1]) {
				//log.Println(resss[x][1])

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

				//log.Println(tesd)

				//return goodsType,goodsInfo
				json.Unmarshal([]byte(tesd), &goodsType)
				log.Println(goodsType)

			}
			///------------------------info--------------
			reg = regexp.MustCompile(`\s"sku"\s:\s({[\s\S]+})\s+};\s+iDetailData`)

			if reg.MatchString(resss[x][1]) {
				ress := reg.FindAllStringSubmatch(resss[x][1], -1)

				log.Println(ress[0][1])

				err := json.Unmarshal([]byte(ress[0][1]), &goodsInfo)
				log.Println(err)
				log.Println(goodsInfo)
			}

			log.Println("----------------")

		}

	}

	return goodsType, goodsInfo
}
func (m *Catch1688Service) Catch(CatchContent, Mark string, isGbk bool) {
	//b, err := ioutil.ReadAll(res.Body)

	addPriceRotia := 0.3
	brokerageRotia := 0.1
	expresstePrice := 500

	haveAdmin := m.Admin.FindAdminByAccount(db.Orm(), "admin")
	//邮件模板
	express := m.ExpressTemplate.GetExpressTemplateByOID(haveAdmin.OID)
	if express.ID == 0 {
		log.Println("没有创建快递模板无法添加产品")
		return
	}

	//content_item := service.URLS[i]

	_havg := m.Goods.FindGoodsLikeMark(Mark)
	if _havg.ID > 0 {
		return
	}

	body := bytes.NewBufferString(CatchContent)

	var doc *goquery.Document
	var err error
	if isGbk {
		reader := transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
		doc, err = goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		doc, err = goquery.NewDocumentFromReader(body)
		if err != nil {
			log.Fatal(err)
		}
	}

	goods := &model.Goods{}

	doc.Find("h1.d-title").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		title = strings.ReplaceAll(title, "批发", "")
		//fmt.Println(title)
		goods.Title = title
	})

	_goods := m.Goods.FindGoodsByTitle(haveAdmin.OID, goods.Title)
	if _goods.ID != 0 {
		return
	}
	//goods.Mark = Mark
	goods.OID = haveAdmin.OID
	goods.ExpressTemplateID = express.ID

	docHtml, err := doc.Html()
	log.Println(err)

	goodsType, goodsInfo := m.getGoodsType(docHtml)
	log.Println(goodsType, goodsInfo)

	totlStock := uint(0)
	minPrice := math.MaxFloat64

	specifications := make([]model.Specification, 0)

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

			specification := model.Specification{}
			specification.Label = html.UnescapeString(key)
			specification.MarketPrice = uint(_price)
			specification.Brokerage = uint(brokerage)
			specification.CostPrice = uint(costPrice)
			specification.Stock = uint(_canBookCount)
			specification.Num = 1
			specification.Weight = 200
			specification.OID = haveAdmin.OID
			specifications = append(specifications, specification)

			totlStock = totlStock + uint(specification.Stock)
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

		specification := model.Specification{}
		specification.Label = html.UnescapeString(goods.Title)
		specification.MarketPrice = uint(minPrice)
		specification.Brokerage = uint(brokerage)
		specification.CostPrice = uint(costPrice)
		specification.Stock = uint(totlStock)
		specification.Num = 1
		specification.Weight = 200
		specification.OID = haveAdmin.OID
		specifications = append(specifications, specification)

	}

	if len(specifications) == 0 {
		return
	}

	goods.Price = uint(minPrice)
	goods.Stock = uint(totlStock)

	categoryList := goodsType["categoryList"].([]interface{})

	categoryA := categoryList[0].(map[string]interface{})
	categoryB := categoryList[1].(map[string]interface{})

	gt, gtc := m.GoodsTypeService.AddGoodsTypeByNameByChild(categoryA["name"].(string), categoryB["name"].(string))
	goods.GoodsTypeID = gt.ID
	goods.GoodsTypeChildID = gtc.ID

	imageList, exist := doc.Find("div.mod-detail-version2018-gallery").Attr("data-gallery-image-list")
	if exist {
		imageLists := strings.Split(imageList, ",")
		log.Println(imageLists)

		//images := make([]string, 0)
		for i := 0; i < len(imageLists); i++ {
			if strings.EqualFold(imageLists[i], "") == false {
				imgPath := "" //todo tool.DownloadInternetImage(imageLists[i], "", "")
				log.Println(imgPath)
				if strings.EqualFold(imgPath, "") == false {
					//todo images = append(images, "//"+conf.Config.Domain+"/file/load?path="+imgPath)
					time.Sleep(200 * time.Millisecond)
				}

			}

		}

		//goods.Images = util.StructToJSON(images)
	} else {
		//images := make([]string, 0)
		doc.Find("div.mod-detail-gallery li.tab-trigger").Each(func(i int, selection *goquery.Selection) {

			imgPath, exist := selection.Attr("data-imgs")
			if exist {
				//log.Println(imgPath)
				imgPathMap := make(map[string]interface{})
				util.JSONToStruct(imgPath, &imgPathMap)
				if imgPathMap["original"] != nil {
					//todo imgPath := "" //todo tool.DownloadInternetImage(imgPathMap["original"].(string), "", "")
					//todo images = append(images, "//"+conf.Config.Domain+"/file/load?path="+imgPath)
					time.Sleep(200 * time.Millisecond)
				}
			}

		})
		//goods.Images = util.StructToJSON(images)

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

	//goods.Params = util.StructToJSON(attributes)

	/*doc.Find(".table-sku tbody tr").Each(func(i int, s *goquery.Selection) {

		spe := model.Specification{}
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
			log.Println(err)
			//b, err := ioutil.ReadAll(res.Body)
			reader := transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())
			b, err := ioutil.ReadAll(reader)
			log.Println(err)
			//fmt.Println(string(b))
			te := strings.TrimSpace(string(b))
			//fmt.Println(te)

			reg := regexp.MustCompile(`^var offer_details={"content":"(.+)"};`)
			ress := reg.FindStringSubmatch(te)
			//fmt.Println(len(ress))
			//fmt.Println(ress[1])

			red := strings.NewReader(ress[1])
			imgsDoc, err := goquery.NewDocumentFromReader(red)
			log.Println(err)

			//images := make([]string, 0)
			imgsDoc.Find("img").Each(func(i int, s *goquery.Selection) {
				pimgUrl, exits := s.Attr("src")
				if exits && strings.EqualFold(pimgUrl, "") == false {
					pimgUrl = strings.ReplaceAll(pimgUrl, `\`, "")
					pimgUrl = strings.ReplaceAll(pimgUrl, `"`, "")
					imgPath := "" //todo tool.DownloadInternetImage(pimgUrl, "", "")

					if strings.EqualFold(imgPath, "") == false {
						//todo images = append(images, "//"+conf.Config.Domain+"/file/load?path="+imgPath)
						time.Sleep(200 * time.Millisecond)
					}

				}
			})

			//goods.Pictures = util.StructToJSON(images)

		}

	})

	//goods.Videos = "[]"

	dao.Create(db.Orm(), goods)
	for s := 0; s < len(specifications); s++ {
		specifications[s].GoodsID = goods.ID
		dao.Create(db.Orm(), &(specifications[s]))
	}
}
func (m *Catch1688Service) URLCatch(URL string) {

	res, err := http.Get(URL)
	if err != nil {
		log.Println(err)
		time.Sleep(time.Hour * 3)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	log.Println(err)
	m.Catch(string(b), URL, true)

}
