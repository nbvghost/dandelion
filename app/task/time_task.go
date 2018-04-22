package service

import (
	"net/http"

	"fmt"
	"io/ioutil"

	"github.com/revel/modules/db/app"
	"server.local/gweb/tool"

	"encoding/json"
	"net/url"

	"strings"
)

func init() {

}
func ReadANetArticle() {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "http://ali-weixin-hot.showapi.com/articleTypeList", nil)
	tool.CheckError(err)
	req.Header.Add("Authorization", "APPCODE c69bc2ff3a994d22925d7a06709801fb")
	resp, err := client.Do(req)
	tool.CheckError(err)
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(b))

	d := make(map[string]interface{})
	err = json.Unmarshal(b, &d)
	if err != nil {
		tool.CheckError(err)
		fmt.Println(string(b))
		return
	}

	if d["showapi_res_body"] != nil {
		if d["showapi_res_body"].(map[string]interface{})["typeList"] != nil {

			arr := d["showapi_res_body"].(map[string]interface{})["typeList"].([]interface{})

			for _, value := range arr {
				item := value.(map[string]interface{})
				fmt.Println(item)

				v := url.Values{}
				v.Set("typeId", item["id"].(string))
				v.Set("page", "10")

				req, err := http.NewRequest("GET", "http://ali-weixin-hot.showapi.com/articleDetalList?"+v.Encode(), nil)
				tool.CheckError(err)
				req.Header.Add("Authorization", "APPCODE c69bc2ff3a994d22925d7a06709801fb")
				resp, err := client.Do(req)
				tool.CheckError(err)
				b, err := ioutil.ReadAll(resp.Body)
				fmt.Println(string(b))

				da := make(map[string]interface{})
				err = json.Unmarshal(b, &da)
				tool.CheckError(err)

				if da["showapi_res_body"] != nil {
					if da["showapi_res_body"].(map[string]interface{})["pagebean"] != nil {

						if da["showapi_res_body"].(map[string]interface{})["pagebean"].(map[string]interface{})["contentlist"] != nil {
							list := da["showapi_res_body"].(map[string]interface{})["pagebean"].(map[string]interface{})["contentlist"].([]interface{})
							for _, value := range list {

								_article := value.(map[string]interface{})
								fmt.Println(_article["url"])

								if _article["title"] == nil || _article["contentImg"] == nil {
									continue
								}

								Title := _article["title"].(string)
								Url := _article["url"].(string)
								if db.GetArticleByTitle(Title).ID != 0 {
									continue
								}

								html := ReadWeiXinArticle(Url)
								if strings.EqualFold(html, "") {
									continue
								}

								categ, _ := db.AddCategory(item["name"].(string))

								article := &db.Article{}
								article.Title = Title
								article.FromUrl = Url
								article.CategoryID = categ.ID
								article.Thumbnail = DownNetImage(_article["contentImg"].(string))
								article.Content = html
								db.AddArticle(article)

							}
						}

					}
				}

			}

		}
	}

}
func ReadBNetArticle() {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "http://jisuwxwzjx.market.alicloudapi.com/weixinarticle/channel", nil)
	tool.CheckError(err)
	req.Header.Add("Authorization", "APPCODE c69bc2ff3a994d22925d7a06709801fb")
	resp, err := client.Do(req)
	tool.CheckError(err)
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(b))

	d := make(map[string]interface{})
	err = json.Unmarshal(b, &d)
	if err != nil {
		tool.CheckError(err)
		fmt.Println(string(b))
		return
	}

	if d["result"] != nil {

		arr := d["result"].([]interface{})

		for _, value := range arr {
			item := value.(map[string]interface{})

			v := url.Values{}
			v.Set("channelid", item["channelid"].(string))

			req, err := http.NewRequest("GET", "http://jisuwxwzjx.market.alicloudapi.com/weixinarticle/get?"+v.Encode(), nil)
			tool.CheckError(err)
			req.Header.Add("Authorization", "APPCODE c69bc2ff3a994d22925d7a06709801fb")
			resp, err := client.Do(req)
			tool.CheckError(err)
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				tool.CheckError(err)
				continue
			}

			da := make(map[string]interface{})
			err = json.Unmarshal(b, &da)
			tool.CheckError(err)

			if da["result"] != nil {
				if da["result"].(map[string]interface{})["list"] != nil {

					list := da["result"].(map[string]interface{})["list"].([]interface{})
					for _, value := range list {

						_article := value.(map[string]interface{})

						if _article["title"] == nil || _article["pic"] == nil || _article["content"] == nil {
							continue
						}

						Title := _article["title"].(string)
						Url := _article["url"].(string)
						if db.GetArticleByTitle(Title).ID != 0 {
							continue
						}

						html := ReadWeiXinArticle(Url)
						if strings.EqualFold(html, "") {
							html = _article["content"].(string)
						}

						categ, _ := db.AddCategory(item["channel"].(string))

						article := &db.Article{}
						article.Title = Title
						article.FromUrl = Url
						article.CategoryID = categ.ID
						article.Thumbnail = DownNetImage(_article["pic"].(string))
						article.Content = html
						db.AddArticle(article)

					}
				}

			}

		}

	}

}
