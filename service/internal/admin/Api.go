package admin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
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
	client.Jar = jar
}

type Api struct {
}

func (Api) AddGoodsTypeChild(GoodsTypeID dao.PrimaryKey, Name string) (*model.GoodsTypeChild, error) {
	goodsType := &model.GoodsTypeChild{Name: Name, GoodsTypeID: GoodsTypeID}

	goodsBytes, err := json.Marshal(goodsType)
	if err != nil {
		return nil, err
	}

	form, err := client.Post("https://admin.sites.ink/api/goods/add-goods-type-child", "application/json", bytes.NewReader(goodsBytes))
	if err != nil {
		return nil, err
	}
	defer form.Body.Close()

	body, err := io.ReadAll(form.Body)
	if err != nil {
		return nil, err
	}

	ar := struct {
		Code    result.ActionResultCode
		Message string
		Data    *model.GoodsTypeChild
		Now     int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}
	return ar.Data, nil
}
func (Api) GetGoodsTypeChild(ID dao.PrimaryKey, Name string) (*model.GoodsTypeChild, error) {
	params := url.Values{}
	params.Set("ID", fmt.Sprintf("%d", ID))
	params.Set("Name", Name)

	form, err := client.Get("https://admin.sites.ink/api/goods/get-goods-type-child?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer form.Body.Close()

	body, err := io.ReadAll(form.Body)
	if err != nil {
		return nil, err
	}

	ar := struct {
		Code    result.ActionResultCode
		Message string
		Data    *model.GoodsTypeChild
		Now     int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	return ar.Data, nil
}

func (Api) AddGoodsType(Name string) (*model.GoodsType, error) {
	goodsType := &model.GoodsType{Name: Name}

	goodsBytes, err := json.Marshal(goodsType)
	if err != nil {
		return nil, err
	}

	form, err := client.Post("https://admin.sites.ink/api/goods/add-goods-type", "application/json", bytes.NewReader(goodsBytes))
	if err != nil {
		return nil, err
	}
	defer form.Body.Close()

	body, err := io.ReadAll(form.Body)
	if err != nil {
		return nil, err
	}

	ar := struct {
		Code    result.ActionResultCode
		Message string
		Data    *model.GoodsType
		Now     int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}
	return ar.Data, nil
}
func (Api) GetGoodsType(ID dao.PrimaryKey, Name string) (*model.Goods, error) {
	params := url.Values{}
	params.Set("ID", fmt.Sprintf("%d", ID))
	params.Set("Name", Name)

	form, err := client.Get("https://admin.sites.ink/api/goods/get-goods-type?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer form.Body.Close()

	body, err := io.ReadAll(form.Body)
	if err != nil {
		return nil, err
	}

	ar := struct {
		Code    result.ActionResultCode
		Message string
		Data    *model.Goods
		Now     int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	return ar.Data, nil
}

func (Api) UpdateGoods(ID dao.PrimaryKey, Title string, Images sqltype.Array[string], GoodsTypeID dao.PrimaryKey, GoodsTypeChildID dao.PrimaryKey) (*model.Goods, error) {
	//https://admin.sites.ink/api/goods/change-goods

	goods := &model.Goods{Entity: dao.Entity{ID: ID}}
	goods.Images = Images
	goods.Title = Title
	goods.GoodsTypeID = GoodsTypeID
	goods.GoodsTypeChildID = GoodsTypeChildID

	goodsBytes, err := json.Marshal(goods)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("goods", string(goodsBytes))

	form, err := client.PostForm("https://admin.sites.ink/api/goods/change-goods", params)
	if err != nil {
		return nil, err
	}
	defer form.Body.Close()

	body, err := io.ReadAll(form.Body)
	if err != nil {
		return nil, err
	}

	ar := struct {
		Code    result.ActionResultCode
		Message string
		Data    struct {
			Goods *model.Goods
		}
		Now int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	return ar.Data.Goods, nil
}
func (Api) AddGoods(Title string, GoodsTypeID dao.PrimaryKey, GoodsTypeChildID dao.PrimaryKey) (*model.Goods, error) {
	//http://admin.dev.com/api/goods/add-goods

	params := url.Values{}
	params.Set("goods", fmt.Sprintf(`{"Title":"%s","GoodsTypeID":%d,"GoodsTypeChildID":%d}`, Title, GoodsTypeID, GoodsTypeChildID))

	form, err := client.PostForm("https://admin.sites.ink/api/goods/add-goods", params)
	if err != nil {
		return nil, err
	}
	defer form.Body.Close()

	body, err := io.ReadAll(form.Body)
	if err != nil {
		return nil, err
	}

	ar := struct {
		Code    result.ActionResultCode
		Message string
		Data    struct {
			Goods *model.Goods
		}
		Now int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	return ar.Data.Goods, nil
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
