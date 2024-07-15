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
	"github.com/nbvghost/dandelion/service/serviceargument"
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

var baseUrl = "https://admin.sites.ink"

func SetBaseURL(url string) {
	baseUrl = url
}

type Api struct {
}

func (Api) GetSpecification(GoodsID dao.PrimaryKey) ([]model.Specification, error) {
	params := url.Values{}
	params.Set("goods-id", fmt.Sprintf("%d", GoodsID))

	form, err := client.Get(baseUrl + "/api/goods/specification?" + params.Encode())
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
			Specifications []model.Specification
		}
		Now int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}
	return ar.Data.Specifications, nil
}
func (Api) PostSpecification(GoodsID dao.PrimaryKey, si []serviceargument.SpecificationInfo) ([]model.Specification, error) {
	goodsType := map[string]any{
		"GoodsID": GoodsID,
		"List":    si,
	}

	goodsBytes, err := json.Marshal(goodsType)
	if err != nil {
		return nil, err
	}

	form, err := client.Post(baseUrl+"/api/goods/specification", "application/json", bytes.NewReader(goodsBytes))
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
			Specifications []model.Specification
		}
		Now int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}
	return ar.Data.Specifications, nil
}
func (Api) UpdateGoodsSkuLabelData(ID dao.PrimaryKey, label, image string) (map[dao.PrimaryKey][]*model.GoodsSkuLabelData, error) {
	goodsType := &model.GoodsSkuLabelData{
		Label: label,
		Image: image,
	}
	goodsType.ID = ID

	goodsBytes, err := json.Marshal(goodsType)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", baseUrl+"/api/goods/sku-label-data", bytes.NewReader(goodsBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	form, err := client.Do(request)
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
			SkuLabelDataList map[dao.PrimaryKey][]*model.GoodsSkuLabelData
		}
		Now int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}
	return ar.Data.SkuLabelDataList, nil
}
func (Api) AddGoodsSkuLabelData(GoodsID, GoodsSkuLabelID dao.PrimaryKey, label, name, image string) (map[dao.PrimaryKey][]*model.GoodsSkuLabelData, error) {
	goodsType := &model.GoodsSkuLabelData{
		GoodsSkuLabelID: GoodsSkuLabelID,
		GoodsID:         GoodsID,
		Label:           label,
		Name:            name,
		Image:           image,
	}

	goodsBytes, err := json.Marshal(goodsType)
	if err != nil {
		return nil, err
	}

	form, err := client.Post(baseUrl+"/api/goods/sku-label-data", "application/json", bytes.NewReader(goodsBytes))
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
			SkuLabelDataList map[dao.PrimaryKey][]*model.GoodsSkuLabelData
		}
		Now int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}
	return ar.Data.SkuLabelDataList, nil
}
func (Api) GetSKULabelData(GoodsID dao.PrimaryKey, goodsSkuLabelData *model.GoodsSkuLabelData) (map[dao.PrimaryKey][]*model.GoodsSkuLabelData, error) {
	//https://admin.sites.ink/api/goods/sku-label
	params := url.Values{}
	params.Set("goods-id", fmt.Sprintf("%d", GoodsID))
	params.Set("goods-sku-label-id", fmt.Sprintf("%d", goodsSkuLabelData.ID))
	params.Set("name", goodsSkuLabelData.Name)

	request, err := http.NewRequest("GET", baseUrl+"/api/goods/sku-label-data?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	form, err := client.Do(request)
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
			SkuLabelDataList map[dao.PrimaryKey][]*model.GoodsSkuLabelData
		}
		Now int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}
	return ar.Data.SkuLabelDataList, nil
}
func (Api) PutGoodsSKULabel(GoodsID dao.PrimaryKey, GoodsSkuLabelList []model.GoodsSkuLabel) ([]model.GoodsSkuLabel, error) {
	//https://admin.sites.ink/api/goods/sku-label
	goodsType := map[string]any{
		"GoodsID":   GoodsID,
		"LabelList": GoodsSkuLabelList,
	}

	goodsBytes, err := json.Marshal(goodsType)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", baseUrl+"/api/goods/sku-label", bytes.NewReader(goodsBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	form, err := client.Do(request)
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
			SkuLabelList []model.GoodsSkuLabel
		}
		Now int64
	}{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}
	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}
	return ar.Data.SkuLabelList, nil
}
func (Api) AddGoodsTypeChild(GoodsTypeID dao.PrimaryKey, Name string) (*model.GoodsTypeChild, error) {
	goodsType := &model.GoodsTypeChild{Name: Name, GoodsTypeID: GoodsTypeID}

	goodsBytes, err := json.Marshal(goodsType)
	if err != nil {
		return nil, err
	}

	form, err := client.Post(baseUrl+"/api/goods/add-goods-type-child", "application/json", bytes.NewReader(goodsBytes))
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
func (Api) GetGoodsTypeChild(ID dao.PrimaryKey, Name string, GoodsTypeID dao.PrimaryKey) (*model.GoodsTypeChild, error) {
	params := url.Values{}
	params.Set("ID", fmt.Sprintf("%d", ID))
	params.Set("GoodsTypeID", fmt.Sprintf("%d", GoodsTypeID))
	params.Set("Name", Name)

	form, err := client.Get(baseUrl + "/api/goods/get-goods-type-child?" + params.Encode())
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
func (Api) AddGoodsType(Name string) (*model.GoodsType, error) {
	goodsType := &model.GoodsType{Name: Name}

	goodsBytes, err := json.Marshal(goodsType)
	if err != nil {
		return nil, err
	}

	form, err := client.Post(baseUrl+"/api/goods/add-goods-type", "application/json", bytes.NewReader(goodsBytes))
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
func (Api) GetGoodsType(ID dao.PrimaryKey, Name string) (*model.GoodsType, error) {
	params := url.Values{}
	params.Set("ID", fmt.Sprintf("%d", ID))
	params.Set("Name", Name)

	form, err := client.Get(baseUrl + "/api/goods/get-goods-type?" + params.Encode())
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
func (Api) UpdateGoods(ID dao.PrimaryKey, goods *model.Goods, Title string, Images sqltype.Array[string], GoodsTypeID dao.PrimaryKey, GoodsTypeChildID dao.PrimaryKey) (*model.Goods, error) {
	//https://admin.sites.ink/api/goods/change-goods

	//goods := &model.Goods{Entity: dao.Entity{ID: ID}}
	goods.ID = ID
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

	form, err := client.PostForm(baseUrl+"/api/goods/change-goods", params)
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
	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}
	return ar.Data.Goods, nil
}
func (Api) AddGoods(Title string, GoodsTypeID dao.PrimaryKey, GoodsTypeChildID dao.PrimaryKey) (*model.Goods, error) {
	//http://admin.dev.com/api/goods/add-goods

	params := url.Values{}
	params.Set("goods", fmt.Sprintf(`{"Title":"%s","GoodsTypeID":%d,"GoodsTypeChildID":%d}`, Title, GoodsTypeID, GoodsTypeChildID))

	form, err := client.PostForm(baseUrl+"/api/goods/add-goods", params)
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
	if ar.Code != 0 {
		return ar.Data.Goods, errors.New(ar.Message)
	}
	return ar.Data.Goods, nil
}
func (Api) Login(account, password string) error {
	params := url.Values{}
	params.Set("account", account)
	params.Set("password", password)

	form, err := client.PostForm(baseUrl+"/api/account/login", params)
	if err != nil {
		return err
	}
	defer form.Body.Close()

	body, err := io.ReadAll(form.Body)
	if err != nil {
		return err
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
		return err
	}
	if ar.Code != 0 {
		return errors.New(ar.Message)
	}
	log.Println(string(body))
	return nil
}
