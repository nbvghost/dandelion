package goods

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/internal/repository"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/dandelion/service/pinyin"
	"github.com/nbvghost/gpa/types"
	"github.com/pkg/errors"
)

type GoodsService struct {
	model.BaseDao
	TimeSell             activity.TimeSellService
	Collage              activity.CollageService
	PinyinService        pinyin.Service
	GoodsTypeService     GoodsTypeService
	AttributesService    AttributesService
	SpecificationService SpecificationService
}

func (service GoodsService) PaginationGoods(OID, GoodsTypeID, GoodsTypeChildID types.PrimaryKey, pageIndex int) (int, int, int, []*model.Goods) {
	//pageIndex, pageSize, total, list, err :
	if GoodsTypeID == 0 {
		db := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"OID"=?`, OID)
		total := db.Limit(pageIndex, 18)
		goodsList := db.List()
		list := make([]*model.Goods, 0)
		for i := range goodsList {
			list = append(list, goodsList[i].(*model.Goods))
		}
		return pageIndex, 18, int(total), list //repository.Goods.FindByOIDLimit(OID, params.NewLimit(pageIndex, 18))
	}
	if GoodsTypeChildID == 0 {
		db := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"OID"=? and "GoodsTypeID"=?`, OID, GoodsTypeID)
		total := db.Limit(pageIndex, 20)
		goodsList := db.List()
		list := make([]*model.Goods, 0)
		for i := range goodsList {
			list = append(list, goodsList[i].(*model.Goods))
		}
		return pageIndex, 20, int(total), list
		//return repository.Goods.FindByOIDAndGoodsTypeIDLimit(OID, GoodsTypeID, params.NewLimit(pageIndex, 18))
	}
	db := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"OID"=? and "GoodsTypeID"=? and "GoodsTypeChildID"=?`, OID, GoodsTypeID, GoodsTypeChildID)
	total := db.Limit(pageIndex, 20)
	goodsList := db.List()
	list := make([]*model.Goods, 0)
	for i := range goodsList {
		list = append(list, goodsList[i].(*model.Goods))
	}
	return pageIndex, 20, int(total), list
	//return repository.Goods.FindByOIDAndGoodsTypeIDAndGoodsTypeChildIDLimit(OID, GoodsTypeID, GoodsTypeChildID, params.NewLimit(pageIndex, 20))
}

/*
	func (service GoodsService) AddSpecification(context *gweb.Context) (r gweb.Result,err error) {
		content_item := &model.Specification{}
		err := util.RequestBodyToJSON(context.Request.Body, content_item)
		if err != nil {
			return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
		}
		err = service.Add(Orm, content_item)
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}
	}

	func (service GoodsService) ListSpecification(context *gweb.Context) (r gweb.Result,err error) {
		GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
		var dts []model.Specification
		service.FindWhere(Orm, &dts, model.Specification{GoodsID: GoodsID})
		return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: dts}}
	}
*/

func (service GoodsService) getGoodsByUri(OID types.PrimaryKey, uri string) model.Goods {
	Orm := singleton.Orm()
	var goods model.Goods
	goods.OID = OID
	goods.Uri = uri
	Orm.Model(model.Goods{}).Where(map[string]interface{}{"OID": goods.OID, "Uri": goods.Uri}).First(&goods)
	return goods
}
func (service GoodsService) SaveGoods(OID types.PrimaryKey, goods *model.Goods, specifications []model.Specification) error {
	Orm := singleton.Orm()
	var err error
	tx := Orm.Begin()

	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	if len(goods.Title) == 0 {
		return errors.Errorf("请指定产品标题")
	}

	g := service.FindGoodsByTitle(goods.Title)
	if !g.IsZero() && g.ID != goods.ID {
		return errors.Errorf("重复的产品标题")
	}

	uri := service.PinyinService.AutoDetectUri(goods.Title)
	g = service.getGoodsByUri(OID, uri)
	if !g.IsZero() {
		uri = fmt.Sprintf("%s-%d", uri, time.Now().Unix())
	}
	goods.Uri = uri
	if goods.ID.IsZero() {
		err = tx.Create(goods).Error
	} else {
		//err = tx.Save(goods).Error
		err = tx.Model(goods).Updates(goods).Error
	}

	if err != nil {
		return err
	}

	//添加或修改的时候不删除规格
	/*err=service.UnscopedDeleteWhere(tx,&model.Specification{},"GoodsID=?",goods.ID)
	if err!=nil{
		return err
	}*/

	if len(specifications) > 0 {
		var total uint
		for _, value := range specifications {
			value.GoodsID = goods.ID
			if value.ID.IsZero() {
				err = tx.Create(&value).Error
				total = total + value.Stock
			} else {
				err = tx.Save(&value).Error
				//err = tx.Model(&goods).Updates(goods).Error
				total = total + value.Stock
			}
			if err != nil {
				return err
			}
		}
		goods.Stock = total
	}
	err = tx.Save(goods).Error
	return err
}

func (service GoodsService) GetGoodsInfo(goods *model.Goods) (*extends.GoodsInfo, error) {
	Orm := singleton.Orm()

	//Orm := singleton.Orm()

	//user := model.User{}
	//service.User.Get(Orm, UserID, &user)
	//brokerageProvisoConf := service.Configuration.GetConfiguration(play.ConfigurationKey_BrokerageProviso)
	//brokerageProvisoConfV, _ := strconv.ParseUint(brokerageProvisoConf.V, 10, 64)
	//vipdiscountConf := service.Configuration.GetConfiguration(play.ConfigurationKey_VIPDiscount)
	//VIPDiscount, _ := strconv.ParseUint(vipdiscountConf.V, 10, 64)
	timeSell := service.TimeSell.GetTimeSellByGoodsID(goods.ID, goods.OID)
	goodsInfo := extends.GoodsInfo{}
	goodsInfo.Goods = *goods
	goodsInfo.GoodsType = service.GoodsTypeService.GetGoodsType(goods.GoodsTypeID)
	goodsInfo.GoodsTypeChild = service.GoodsTypeService.GetGoodsTypeChild(goods.GoodsTypeChildID)
	goodsInfo.Discounts = make([]extends.Discount, 0)

	var goodsRating extends.GoodsRating
	singleton.Orm().Model(&model.GoodsReview{}).Where(`"GoodsID"=?`, goods.ID).Select(`SUM("Rating") as "Rating",COUNT("ID") as "RatingCount"`).Scan(&goodsRating)
	goodsInfo.Rating = goodsRating

	if timeSell.IsEnable() {
		//Favoured:=uint(util.Rounding45(float64(goods.Price)*(float64(timeSell.Discount)/float64(100)), 2))
		goodsInfo.Discounts = append(goodsInfo.Discounts, extends.Discount{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: "TimeSell", Discount: uint(timeSell.Discount)})
	} else {
		collage := service.Collage.GetCollageByGoodsID(goods.ID, goods.OID)
		if collage.ID != 0 && collage.TotalNum > 0 {
			goodsInfo.Discounts = append(goodsInfo.Discounts, extends.Discount{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint(collage.Discount)})
		}

	}

	goodsSkuLabelDataMap := make(map[types.PrimaryKey][]*model.GoodsSkuLabelData, 0)
	{
		goodsSkuLabelData := dao.Find(Orm, &model.GoodsSkuLabelData{}).Where(`"GoodsID"=?`, goods.ID).List()
		for i := range goodsSkuLabelData {
			item := goodsSkuLabelData[i].(*model.GoodsSkuLabelData)
			if _, ok := goodsSkuLabelDataMap[item.GoodsSkuLabelID]; !ok {
				goodsSkuLabelDataMap[item.GoodsSkuLabelID] = make([]*model.GoodsSkuLabelData, 0)
			}
			goodsSkuLabelDataMap[item.GoodsSkuLabelID] = append(goodsSkuLabelDataMap[item.GoodsSkuLabelID], item)
		}
	}

	goodsSkuLabel := dao.Find(Orm, &model.GoodsSkuLabel{}).Where(`"GoodsID"=? and "Abel"=?`, goods.ID, true).Order(`"Image" desc`).List()
	skuLabelList := make([]extends.SkuLabel, len(goodsSkuLabel))
	for i := range goodsSkuLabel {
		item := goodsSkuLabel[i].(*model.GoodsSkuLabel)
		skuLabelList[i].Label = item
		skuLabelList[i].Data = goodsSkuLabelDataMap[item.ID]
	}

	specifications := dao.Find(Orm, &model.Specification{}).Where(`"GoodsID"=?`, goods.ID).Order(`"CreatedAt" desc`).List() //service.FindWhere(Orm, &specifications, model.Specification{GoodsID: goods.ID})
	for i := range specifications {
		specification := specifications[i].(*model.Specification)
		goodsInfo.Specifications = append(goodsInfo.Specifications, specification)
	}

	goodsInfo.SkuLabels = skuLabelList

	attributesGroup, err := repository.GoodsAttributesGroup.FindByGoodsID(goods.ID)
	if err != nil {
		return nil, err
	}
	attributes, err := repository.GoodsAttributes.FindByGoodsID(goods.ID)
	if err != nil {
		return nil, err
	}
	for _, v := range attributesGroup {
		goodsAttributes := extends.GoodsAttributes{}
		goodsAttributes.GroupID = v.ID
		goodsAttributes.GroupName = v.Name
		for _, vv := range attributes {
			if vv.GroupID == v.ID {
				goodsAttribute := extends.GoodsAttribute{}
				goodsAttribute.ID = vv.ID
				goodsAttribute.Name = vv.Name
				goodsAttribute.Value = vv.Value
				goodsAttributes.Attrs = append(goodsAttributes.Attrs, goodsAttribute)
			}
		}
		goodsInfo.Attributes = append(goodsInfo.Attributes, goodsAttributes)
	}

	return &goodsInfo, nil
}
func (service GoodsService) GetGoodsInfoList(goodsList []model.Goods) []extends.GoodsInfo {

	var results = make([]extends.GoodsInfo, 0)

	for _, value := range goodsList {
		goodsInfo := extends.GoodsInfo{}
		goodsInfo.Goods = value
		goodsInfo.Discounts = service.GetDiscounts(value.ID, value.OID)
		goodsInfo.GoodsType = service.GoodsTypeService.GetGoodsType(value.GoodsTypeID)
		goodsInfo.GoodsTypeChild = service.GoodsTypeService.GetGoodsTypeChild(value.GoodsTypeChildID)
		results = append(results, goodsInfo)
	}

	return results
}
func (service GoodsService) GetGoods(DB *gorm.DB, context constrain.IContext, ID types.PrimaryKey) (*extends.GoodsInfo, error) {
	Orm := singleton.Orm()
	//var goods model.Goods
	goods := dao.GetByPrimaryKey(Orm, &model.Goods{}, ID).(*model.Goods)

	return service.GetGoodsInfo(goods)
}

func (service GoodsService) DeleteGoods(ID types.PrimaryKey) *result.ActionResult {
	Orm := singleton.Orm()
	tx := Orm.Begin()
	err := dao.DeleteByPrimaryKey(tx, &model.Goods{}, ID)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Where(`"GoodsID"=?`, ID).Delete(model.Specification{}).Error
	if err != nil {
		tx.Rollback()
	}

	defer func() {
		if err == nil {
			tx.Commit()
		}
	}()

	return (&result.ActionResult{}).SmartError(err, "删除成功", nil)
}

func (service GoodsService) DeleteTimeSellGoods(DB *gorm.DB, GoodsID types.PrimaryKey, OID types.PrimaryKey) error {
	timesell := service.TimeSell.GetTimeSellByGoodsID(GoodsID, OID)
	err := dao.DeleteBy(DB, &model.TimeSellGoods{}, map[string]interface{}{
		"TimeSellHash": timesell.Hash,
		"GoodsID":      GoodsID,
	}) //Delete(DB, &model.TimeSellGoods{}, timesell.ID)
	log.Println(err)
	return err
}
func (service GoodsService) DeleteCollageGoods(DB *gorm.DB, GoodsID types.PrimaryKey, OID types.PrimaryKey) error {
	timesell := service.Collage.GetCollageByGoodsID(GoodsID, OID)

	err := dao.DeleteBy(DB, &model.CollageGoods{}, map[string]interface{}{
		"CollageHash": timesell.Hash,
		"GoodsID":     GoodsID,
	}) //Delete(DB, &model.TimeSellGoods{}, timesell.ID)
	log.Println(err)
	return err

	//err := service.Delete(DB, &model.Collage{}, timesell.ID)
	//log.Println(err)
	//return err
}
func (service GoodsService) FindGoodsByTimeSellID(TimeSellID types.PrimaryKey) []model.Goods {
	Orm := singleton.Orm()

	//var timesell model.TimeSell
	timesell := dao.GetByPrimaryKey(Orm, &model.TimeSell{}, TimeSellID).(*model.TimeSell)

	var list []model.Goods

	//rows, _ := Orm.Raw("")

	//todo:没有写完整
	//err = service.FindWhere(Orm, &list, "ID=?", timesell.GoodsID)
	//log.Println(err)
	//dao.FindBy(Orm, &model.Goods{}, map[string]any{"ID": timesell.})
	log.Println(timesell)
	return list
}
func (service GoodsService) FindGoodsByTimeSellHash(Hash string) []types.IEntity {
	Orm := singleton.Orm()

	var GoodsIDs []types.PrimaryKey
	Orm.Model(&model.TimeSell{}).Where(`"Hash"=?`, Hash).Pluck("GoodsID", &GoodsIDs)

	//var list []model.Goods
	//err := service.FindWhere(Orm, &list, `"ID" in (?)`, GoodsIDs)
	//list := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"ID" in (?)`, GoodsIDs).List()
	return service.ListGoodsByIDs(GoodsIDs)
}
func (service GoodsService) FindGoodsByCollageHash(Hash string) []types.IEntity {
	Orm := singleton.Orm()

	var GoodsIDs []types.PrimaryKey
	Orm.Model(&model.Collage{}).Where(`"Hash"=?`, Hash).Pluck("GoodsID", &GoodsIDs)

	//var list []model.Goods
	//err := service.FindWhere(Orm, &list, `"ID" in (?)`, GoodsIDs)
	//log.Println(err)
	//return list
	//list := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"ID" in (?)`, GoodsIDs).List()
	return service.ListGoodsByIDs(GoodsIDs)
}
func (service GoodsService) ListGoodsByIDs(goodsIDs []types.PrimaryKey) []types.IEntity {
	list := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"ID" in (?)`, goodsIDs).List()
	return list
}
func (service GoodsService) FindGoodsByOrganizationIDAndGoodsID(OrganizationID types.PrimaryKey, GoodsID types.PrimaryKey) model.Goods {
	var Goods model.Goods
	singleton.Orm().Model(&model.Goods{}).Where(`"ID"=? and "OID"=?`, GoodsID, OrganizationID).First(&Goods)
	return Goods
}
func (service GoodsService) FindGoodsByTitle(Title string) model.Goods {
	var Goods model.Goods
	singleton.Orm().Model(&model.Goods{}).Where(`"Title"=?`, Title).First(&Goods)
	return Goods
}
func (service GoodsService) FindGoodsLikeMark(Mark string) model.Goods {
	var Goods model.Goods
	singleton.Orm().Model(&model.Goods{}).Where(`"Mark" like ?`, "%"+Mark+"%").First(&Goods)
	return Goods
}
func (service GoodsService) AllList() []model.Goods {

	Orm := singleton.Orm()

	var result []model.Goods

	db := Orm.Model(&model.Goods{}).Order(`"CreatedAt" desc`) //.Limit(10)

	db.Find(&result)

	return result

}
func (service GoodsService) GetDiscounts(GoodsID, OID types.PrimaryKey) []extends.Discount {
	discounts := make([]extends.Discount, 0)
	timeSell := service.TimeSell.GetTimeSellByGoodsID(GoodsID, OID)
	if timeSell.IsEnable() {
		//Favoured:=uint(util.Rounding45(float64(value.Price)*(float64(timeSell.Discount)/float64(100)), 2))
		discounts = append(discounts, extends.Discount{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: "TimeSell", Discount: uint(timeSell.Discount)})
	} else {
		collage := service.Collage.GetCollageByGoodsID(GoodsID, OID)
		if collage.ID != 0 && collage.TotalNum > 0 {
			discounts = append(discounts, extends.Discount{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint(collage.Discount)})
		}
	}
	return discounts
}

func (service GoodsService) GoodsList(SqlOrder string, Index, Size int, where interface{}, args ...interface{}) extends.GoodsInfoPagination {
	Orm := singleton.Orm()
	//var goodsList []model.Goods
	//db := Orm.Model(&model.Goods{}).Order("CountSale desc").Limit(10)
	//db.Find(&result)

	pager := service.FindWherePaging(Orm, SqlOrder, model.Goods{}, Index, Size, where, args...)

	return extends.GoodsInfoPagination{
		List:  service.GetGoodsInfoList(pager.Data.([]model.Goods)),
		Total: pager.Total,
		Index: pager.Offset,
		Size:  pager.Limit,
	}
}
