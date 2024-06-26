package catch

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/library/result"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var client = &http.Client{}

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}
	client.Jar =jar
}

type Api struct {
}

func (Api) AddGoods(Title string, GoodsTypeID int, GoodsTypeChildID int) error {
	//http://admin.dev.com/api/goods/add-goods

	params := url.Values{}
	params.Set("goods", fmt.Sprintf(`{"Title":"%s","GoodsTypeID":%d,"GoodsTypeChildID":%d}`, Title, GoodsTypeID, GoodsTypeChildID))

	form, err := client.PostForm("https://admin.sites.ink/api/goods/add-goods", params)
	if err != nil {
		return err
	}
	defer form.Body.Close()

	body, err := io.ReadAll(form.Body)
	if err != nil {
		return err
	}

	ar := result.ActionResult{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return err
	}
	return nil
}
func (Api) Login(account, password string) error {
	params := url.Values{}
	params.Set("account", account)
	params.Set("password", password)

	form, err := client.PostForm("https://admin.sites.ink/api/account/login", params)
	if err != nil {
		return err
	}
	defer form.Body.Close()

	body, err := io.ReadAll(form.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	return nil
}
