package goods

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/domain/cache"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/activity"
	"github.com/nbvghost/dandelion/service/internal/express"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"github.com/nbvghost/tool/object"
	"gorm.io/gorm/clause"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/pkg/errors"
)

type GoodsService struct {
	model.BaseDao
	TimeSell               activity.TimeSellService
	Collage                activity.CollageService
	GoodsTypeService       GoodsTypeService
	AttributesService      AttributesService
	SpecificationService   SpecificationService
	SKUService             SKUService
	ExpressTemplateService express.ExpressTemplateService
}

func (m GoodsService) PaginationGoodsDetail(context constrain.IContext, OID, GoodsTypeID, GoodsTypeChildID dao.PrimaryKey, filterOption []serviceargument.Option, sortMethod *serviceargument.SortMethod, pageIndex, pageSize int) (serviceargument.Pagination[*extends.GoodsDetail], *serviceargument.Options) {
	orm := db.Orm().Model(&model.Goods{}).Where(`"Goods"."OID"=?`, OID)

	if GoodsTypeID > 0 && GoodsTypeChildID == 0 {
		orm.Where(`"Goods"."GoodsTypeID"=?`, GoodsTypeID)
	}
	if GoodsTypeID > 0 && GoodsTypeChildID > 0 {
		orm.Where(`"Goods"."GoodsTypeID"=? and "Goods"."GoodsTypeChildID"=?`, GoodsTypeID, GoodsTypeChildID)
	}

	orm.Joins(`left join "Specification" on ("Goods"."ID" = "Specification"."GoodsID")`)
	orm.Joins(`left join "GoodsAttributes" on ("Goods"."ID" = "GoodsAttributes"."GoodsID")`)
	orm.Joins(`left join "GoodsSkuLabel" on ("Goods"."ID" = "GoodsSkuLabel"."GoodsID")`)
	orm.Joins(`left join "GoodsSkuLabelData" on ("Goods"."ID" = "GoodsSkuLabelData"."GoodsID")`)
	orm.Group(`"Goods"."ID"`)
	if sortMethod != nil {
		if len(sortMethod.Field) > 0 {
			method := "asc"
			if !strings.EqualFold(sortMethod.Method, "asc") {
				method = "desc"
			}
			switch sortMethod.Field {
			case "position":
				orm.Order(`"Goods"."CountSale" desc,"Goods"."CountView" desc,"Goods"."UpdatedAt" desc,"Goods"."ID" desc`)
			case "name":
				orm.Order(fmt.Sprintf(`"Goods"."Title" %s,"Goods"."CountSale" desc,"Goods"."CountView" desc,"Goods"."UpdatedAt" desc,"Goods"."ID" desc`, method))
			case "price":
				orm.Order(fmt.Sprintf(`"Goods"."Price" %s,"Goods"."CountSale" desc,"Goods"."CountView" desc,"Goods"."UpdatedAt" desc,"Goods"."ID" desc`, method))
			}

		} else {
			orm.Order(`"Goods"."CountSale" desc,"Goods"."CountView" desc,"Goods"."UpdatedAt" desc,"Goods"."ID" desc`)
		}
	} else {
		orm.Order(`"Goods"."CountSale" desc,"Goods"."CountView" desc,"Goods"."UpdatedAt" desc,"Goods"."ID" desc`)
	}

	orm.Select(`
"Goods".* as "Goods",
json_agg("Specification") as "Specification",
json_agg("GoodsAttributes") as "GoodsAttributes",
json_agg("GoodsSkuLabel") as "GoodsSkuLabel",
json_agg("GoodsSkuLabelData") as "GoodsSkuLabelData",
json_agg("Specification"."MarketPrice")::jsonb as "MarketPrice",
json_agg("Specification"."Weight")::jsonb as "Weight",
json_agg("Specification"."Label")::jsonb as "Label",
json_agg("Specification"."Num")::jsonb as "Num",
json_agg(concat("GoodsAttributes"."Name", '-', "GoodsAttributes"."Value"))::jsonb as "AttrKey",
json_agg("GoodsSkuLabel"."Name")::jsonb as "SKULabel",
json_agg("GoodsSkuLabelData"."Label")::jsonb as "SKUDataLabel"
`)

	mainOrm := db.Orm().Table(`(?) as g`, orm)

	for _, option := range filterOption {
		switch option.Type {
		case serviceargument.OptionsTypeAttribute:
			mainOrm.Where(fmt.Sprintf(`jsonb_path_exists(g."AttrKey", '$[*] ? (@ == "%s")')`, option.Key+"-"+option.Value[0].Value))
		case serviceargument.OptionsTypeSpecification:
			mainOrm.Where(fmt.Sprintf(`jsonb_path_exists(g."SKULabel", '$[*] ? (@ == "%s")')`, option.Key))
			mainOrm.Where(fmt.Sprintf(`jsonb_path_exists(g."SKUDataLabel", '$[*] ? (@ == "%s")')`, option.Value[0].Value))
		case serviceargument.OptionsTypePackageNum:
			pas := strings.Split(option.Value[0].Value, "-")
			if len(pas) >= 2 {
				mainOrm.Where(fmt.Sprintf(`jsonb_path_exists(g."Num", '$[*] ? (@ >= %d && @ <= %d)')`, object.ParseInt(pas[0]), object.ParseInt(pas[1])))
			}
		case serviceargument.OptionsTypeWeight:
			pas := strings.Split(option.Value[0].Value, "-")
			if len(pas) >= 2 {
				mainOrm.Where(fmt.Sprintf(`jsonb_path_exists(g."Weight", '$[*] ? (@ >= %d && @ <= %d)')`, object.ParseInt(pas[0]), object.ParseInt(pas[1])))
			}
		case serviceargument.OptionsTypePrice:
			pas := strings.Split(option.Value[0].Value, "-")
			if len(pas) >= 2 {
				mainOrm.Where(fmt.Sprintf(`jsonb_path_exists(g."MarketPrice", '$[*] ? (@ >= %d && @ <= %d)')`, object.ParseInt(pas[0]), object.ParseInt(pas[1])))
			}
		}
	}

	var total int64
	mainOrm.Count(&total)

	goodsList := make([]*extends.GoodsDetail, 0)

	mainOrm.Find(&goodsList)

	var options = m.goodsOptions(goodsList)

	if pageIndex < 0 {
		pageIndex = 0
	}

	goodsList = make([]*extends.GoodsDetail, 0)
	mainOrm.Limit(pageSize).Offset(pageSize * pageIndex).Find(&goodsList)

	s := serviceargument.NewPagination[*extends.GoodsDetail](pageIndex, pageSize, int(total), goodsList)

	return s, options
}

func (m GoodsService) goodsOptions(list []*extends.GoodsDetail) *serviceargument.Options {
	var options = &serviceargument.Options{}
	for p := range list {
		detail := list[p]
		{
			for i := range detail.GoodsAttributes {
				item := detail.GoodsAttributes[i]
				options.AddAttributes(serviceargument.OptionsTypeAttribute, item.Name, item.Name, item.Value)
			}
		}
		{
			/*skuList := m.SKUService.SkuLabel(detail.GoodsSkuLabel, detail.GoodsSkuLabelData)
			for i := range skuList {
				item := skuList[i]
				for i2 := range item.Data {
					itemData := item.Data[i2]
					options.AddAttributes(serviceargument.OptionsTypeSpecification, item.Label.Name, item.Label.Label, itemData.Label)
				}
			}*/
		}

		numList := make([]int, 0)
		weightList := make([]int, 0)
		marketPriceList := make([]int, 0)
		for _, spec := range detail.Specification {
			numList = append(numList, int(spec.Num))
			weightList = append(weightList, int(spec.Weight))
			weightList = append(weightList, int(spec.MarketPrice))
		}

		sort.Ints(numList)
		sort.Ints(weightList)
		sort.Ints(marketPriceList)

		{
			for i, u := range numList {
				if i == 0 {
					options.AddAttributes(serviceargument.OptionsTypePackageNum, "PackingQuantity", "Packing quantity", "0-"+object.ParseString(u))
				} else {
					options.AddAttributes(serviceargument.OptionsTypePackageNum, "PackingQuantity", "Packing quantity", object.ParseString(numList[i-1])+"-"+object.ParseString(u))
				}
			}
		}

		{
			for i, u := range weightList {
				if i == 0 {
					options.AddAttributes(serviceargument.OptionsTypeWeight, "Weight", "Weight", "0-"+object.ParseString(u))
				} else {
					options.AddAttributes(serviceargument.OptionsTypeWeight, "Weight", "Weight", object.ParseString(weightList[i-1])+"-"+object.ParseString(u))
				}
			}
		}
		{
			for i, u := range marketPriceList {
				if i == 0 {
					options.AddAttributes(serviceargument.OptionsTypePrice, "Price", "Price", "0-"+object.ParseString(u))
				} else {
					options.AddAttributes(serviceargument.OptionsTypePrice, "Price", "Price", object.ParseString(marketPriceList[i-1])+"-"+object.ParseString(u))
				}
			}
		}
	}
	return options
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

func (m GoodsService) getGoodsByUri(OID dao.PrimaryKey, uri string) model.Goods {
	Orm := db.Orm()
	var goods model.Goods
	goods.OID = OID
	goods.Uri = uri
	Orm.Model(model.Goods{}).Where(map[string]interface{}{"OID": goods.OID, "Uri": goods.Uri}).First(&goods)
	return goods
}

func (m GoodsService) SaveGoods(tx *gorm.DB, OID dao.PrimaryKey, goods *model.Goods, specifications []model.Specification) (*model.Goods, error) {
	if goods.Tags == nil {
		goods.Tags = make(pq.StringArray, 0)
	}

	goods.OID = OID
	if len(goods.Title) == 0 {
		return goods, errors.Errorf("请指定产品标题")
	}

	g := m.FindGoodsByTitle(OID, goods.Title)
	if g.IsZero() == false && g.ID != goods.ID {
		return &g, errors.Errorf("重复的产品标题")
	}

	uri := cache.Cache.ChinesePinyinCache.AutoDetectUri(goods.Title)
	g = m.getGoodsByUri(OID, uri)
	if !g.IsZero() {
		uri = fmt.Sprintf("%s-%d", uri, time.Now().Unix())
	}
	goods.Uri = uri
	if goods.ID.IsZero() {
		err := tx.Create(goods).Error
		if err != nil {
			return goods, err
		}
	} else {
		//err = tx.Save(goods).Error
		err := tx.Model(goods).Updates(goods).Error
		if err != nil {
			return goods, err
		}
	}

	//添加或修改的时候不删除规格
	/*err=service.UnscopedDeleteWhere(tx,&model.Specification{},"GoodsID=?",goods.ID)
	if err!=nil{
		return err
	}*/

	if len(specifications) > 0 {
		var total uint
		for i := range specifications {
			value := specifications[i]
			value.OID = OID
			value.GoodsID = goods.ID
			if value.ID.IsZero() {
				err := tx.Create(&value).Error
				if err != nil {
					return goods, err
				}
				total = total + value.Stock
			} else {
				err := tx.Save(&value).Error
				if err != nil {
					return goods, err
				}
				//err = tx.Model(&goods).Updates(goods).Error
				total = total + value.Stock
			}
		}
		goods.Stock = total
	}

	if goods.ExpressTemplateID == 0 {
		expressTemplate := m.ExpressTemplateService.GetExpressTemplateByOID(OID)
		if expressTemplate.IsZero() {
			return goods, errors.New("请设置快递信息")
		}
		goods.ExpressTemplateID = expressTemplate.ID
	}
	err := tx.Save(goods).Error
	return goods, err
}

func (m GoodsService) Rating(goodsID dao.PrimaryKey) *extends.GoodsRating {
	var goodsRating extends.GoodsRating
	db.Orm().Model(&model.GoodsReview{}).Where(`"GoodsID"=?`, goodsID).Select(`SUM("Rating") as "Rating",COUNT("ID") as "RatingCount"`).Scan(&goodsRating)
	return &goodsRating
}
func (m GoodsService) GetGoodsInfo(goods *model.Goods) (*extends.GoodsMix, error) {
	Orm := db.Orm()

	//Orm := singleton.Orm()

	//user := model.User{}
	//service.User.Get(Orm, UserID, &user)
	//brokerageProvisoConf := service.Configuration.GetConfiguration(play.ConfigurationKey_BrokerageProviso)
	//brokerageProvisoConfV, _ := strconv.ParseUint(brokerageProvisoConf.V, 10, 64)
	//vipdiscountConf := service.Configuration.GetConfiguration(play.ConfigurationKey_VIPDiscount)
	//VIPDiscount, _ := strconv.ParseUint(vipdiscountConf.V, 10, 64)
	timeSell := m.TimeSell.GetTimeSellByGoodsID(goods.ID, goods.OID)
	goodsInfo := extends.GoodsMix{}
	goodsInfo.Goods = *goods
	goodsInfo.GoodsType = m.GoodsTypeService.GetGoodsType(goods.GoodsTypeID)
	goodsInfo.GoodsTypeChild = m.GoodsTypeService.GetGoodsTypeChild(goods.GoodsTypeChildID)
	goodsInfo.Discounts = make([]sqltype.Discount, 0)

	goodsInfo.Rating = *m.Rating(goods.ID)

	if timeSell.IsEnable() {
		//Favoured:=uint(util.Rounding45(float64(goods.Price)*(float64(timeSell.Discount)/float64(100)), 2))
		goodsInfo.Discounts = append(goodsInfo.Discounts, sqltype.Discount{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: "TimeSell", Discount: uint(timeSell.Discount)})
	} else {
		collage := m.Collage.GetCollageByGoodsID(goods.ID, goods.OID)
		if collage.ID != 0 && collage.TotalNum > 0 {
			goodsInfo.Discounts = append(goodsInfo.Discounts, sqltype.Discount{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint(collage.Discount)})
		}

	}

	specifications := dao.Find(Orm, &model.Specification{}).Where(`"GoodsID"=?`, goods.ID).Order(`"LabelIndex"::text asc`).List() //service.FindWhere(Orm, &specifications, model.Specification{GoodsID: goods.ID})
	for i := range specifications {
		specification := specifications[i].(*model.Specification)
		goodsInfo.Specifications = append(goodsInfo.Specifications, specification)
	}

	goodsInfo.SkuLabels = m.SKUService.SkuLabelByGoodsID(Orm, goods.ID) //skuLabelList

	attributesGroup := m.AttributesService.FindGroupByGoodsID(Orm, goods.ID) //repository.GoodsAttributesGroup.FindByGoodsID(goods.ID)
	attributes := m.AttributesService.FindByGoodsID(Orm, goods.ID)           //repository.GoodsAttributes.FindByGoodsID(goods.ID)

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

func (m GoodsService) GetGoods(DB *gorm.DB, context constrain.IContext, ID dao.PrimaryKey) (*extends.GoodsMix, error) {
	Orm := db.Orm()
	//var goods model.Goods
	goods := dao.GetByPrimaryKey(Orm, &model.Goods{}, ID).(*model.Goods)

	return m.GetGoodsInfo(goods)
}

func (m GoodsService) DeleteGoods(ID dao.PrimaryKey) *result.ActionResult {
	Orm := db.Orm()
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

func (m GoodsService) DeleteTimeSellGoods(DB *gorm.DB, GoodsID dao.PrimaryKey, OID dao.PrimaryKey) error {
	timesell := m.TimeSell.GetTimeSellByGoodsID(GoodsID, OID)
	err := dao.DeleteBy(DB, &model.TimeSellGoods{}, map[string]interface{}{
		"TimeSellHash": timesell.Hash,
		"GoodsID":      GoodsID,
	}) //Delete(DB, &model.TimeSellGoods{}, timesell.ID)
	if err != nil {
		log.Println(err)
	}
	return err
}
func (m GoodsService) DeleteCollageGoods(DB *gorm.DB, GoodsID dao.PrimaryKey, OID dao.PrimaryKey) error {
	timesell := m.Collage.GetCollageByGoodsID(GoodsID, OID)

	err := dao.DeleteBy(DB, &model.CollageGoods{}, map[string]interface{}{
		"CollageHash": timesell.Hash,
		"GoodsID":     GoodsID,
	}) //Delete(DB, &model.TimeSellGoods{}, timesell.ID)
	if err != nil {
		log.Println(err)
	}
	return err

	//err := service.Delete(DB, &model.Collage{}, timesell.ID)
	//log.Println(err)
	//return err
}
func (m GoodsService) FindGoodsByTimeSellID(TimeSellID dao.PrimaryKey) []model.Goods {
	Orm := db.Orm()

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
func (m GoodsService) FindGoodsByTimeSellHash(Hash string) []dao.IEntity {
	Orm := db.Orm()

	var GoodsIDs []dao.PrimaryKey
	Orm.Model(&model.TimeSell{}).Where(`"Hash"=?`, Hash).Pluck("GoodsID", &GoodsIDs)

	//var list []model.Goods
	//err := service.FindWhere(Orm, &list, `"ID" in (?)`, GoodsIDs)
	//list := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"ID" in (?)`, GoodsIDs).List()
	return m.ListGoodsByIDs(GoodsIDs)
}
func (m GoodsService) FindGoodsByCollageHash(Hash string) []dao.IEntity {
	Orm := db.Orm()

	var GoodsIDs []dao.PrimaryKey
	Orm.Model(&model.Collage{}).Where(`"Hash"=?`, Hash).Pluck("GoodsID", &GoodsIDs)

	//var list []model.Goods
	//err := service.FindWhere(Orm, &list, `"ID" in (?)`, GoodsIDs)
	//log.Println(err)
	//return list
	//list := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"ID" in (?)`, GoodsIDs).List()
	return m.ListGoodsByIDs(GoodsIDs)
}
func (m GoodsService) ListGoodsByIDs(goodsIDs []dao.PrimaryKey) []dao.IEntity {
	list := dao.Find(db.Orm(), &model.Goods{}).Where(`"ID" in (?)`, goodsIDs).List()
	return list
}
func (m GoodsService) FindGoodsByOrganizationIDAndGoodsID(OrganizationID dao.PrimaryKey, GoodsID dao.PrimaryKey) model.Goods {
	var Goods model.Goods
	db.Orm().Model(&model.Goods{}).Where(`"ID"=? and "OID"=?`, GoodsID, OrganizationID).First(&Goods)
	return Goods
}
func (m GoodsService) FindGoodsByTitle(OID dao.PrimaryKey, Title string) model.Goods {
	var Goods model.Goods
	db.Orm().Model(&model.Goods{}).Where(`"OID"=? and "Title"=?`, OID, Title).First(&Goods)
	return Goods
}
func (m GoodsService) FindGoodsLikeMark(Mark string) model.Goods {
	var Goods model.Goods
	db.Orm().Model(&model.Goods{}).Where(`"Mark" like ?`, "%"+Mark+"%").First(&Goods)
	return Goods
}
func (m GoodsService) AllList() []model.Goods {

	Orm := db.Orm()

	var result []model.Goods

	db := Orm.Model(&model.Goods{}).Order(`"CreatedAt" desc`) //.Limit(10)

	db.Find(&result)

	return result

}
func (m GoodsService) GetDiscounts(GoodsID, OID dao.PrimaryKey) []sqltype.Discount {
	discounts := make([]sqltype.Discount, 0)
	timeSell := m.TimeSell.GetTimeSellByGoodsID(GoodsID, OID)
	if timeSell.IsEnable() {
		//Favoured:=uint(util.Rounding45(float64(value.Price)*(float64(timeSell.Discount)/float64(100)), 2))
		discounts = append(discounts, sqltype.Discount{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: sqltype.DiscountTypeNameTimeSell, Discount: uint(timeSell.Discount)})
	} else {
		collage := m.Collage.GetCollageByGoodsID(GoodsID, OID)
		if collage.ID != 0 && collage.TotalNum > 0 {
			discounts = append(discounts, sqltype.Discount{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: sqltype.DiscountTypeNameCollage, Discount: uint(collage.Discount)})
		}
	}
	return discounts
}

func (m GoodsService) GoodsList(queryParam *serviceargument.ListQueryParam, oid dao.PrimaryKey, orderBy clause.OrderByColumn, pageNo, pageSize int) *result.Pagination {
	Orm := db.Orm()
	//var goodsList []model.Goods
	//db := Orm.Model(&model.Goods{}).Order("CountSale desc").Limit(10)
	//db.Find(&result)

	if pageSize <= 0 {
		pageSize = 10
	}

	pageIndex := pageNo - 1
	if pageIndex < 0 {
		pageIndex = 0
	}

	var goodsList []model.Goods

	orm := Orm.Model(model.Goods{})
	if oid > 0 {
		orm = orm.Where(`"OID"=?`, oid)
	}

	orm = orm.Where(`"Hide"=?`, 0)

	if queryParam != nil {
		if queryParam.GoodsTypeID > 0 {
			orm = orm.Where(`"GoodsTypeID"=?`, queryParam.GoodsTypeID)
		}
		if queryParam.GoodsTypeChildID > 0 {
			orm = orm.Where(`"GoodsTypeChildID"=?`, queryParam.GoodsTypeChildID)
		}
	}

	var recordsTotal int64

	orm.Count(&recordsTotal).Limit(pageSize).Offset(pageSize * pageIndex).Order(orderBy).Order(`"ID" DESC`).Find(&goodsList)
	//pager := service.FindWherePaging(Orm, SqlOrder, model.Goods{}, Index, Size, where, args...)
	return result.NewPagination(pageNo, pageSize, recordsTotal, m.getGoodsInfoList(goodsList))

	/*return extends.GoodsInfoPagination{
		List:  service.GetGoodsInfoList(goodsList),
		Total: recordsTotal,
		Index: pager.Offset,
		Size:  pager.Limit,
	}*/
}
func (m GoodsService) getGoodsInfoList(goodsList []model.Goods) []extends.GoodsMix {
	var results = make([]extends.GoodsMix, 0)
	for _, value := range goodsList {
		goodsInfo := extends.GoodsMix{}
		goodsInfo.Goods = value
		goodsInfo.Discounts = m.GetDiscounts(value.ID, value.OID)
		goodsInfo.GoodsType = m.GoodsTypeService.GetGoodsType(value.GoodsTypeID)
		goodsInfo.GoodsTypeChild = m.GoodsTypeService.GetGoodsTypeChild(value.GoodsTypeChildID)
		results = append(results, goodsInfo)
	}
	return results
}

/*func (m GoodsService) getGoodsInfoList(goodsList []model.Goods) []extends.GoodsMix {
	var results = make([]extends.GoodsMix, 0)
	for _, value := range goodsList {
		goodsInfo := extends.GoodsMix{}
		goodsInfo.Goods = value
		goodsInfo.Discounts = m.GetDiscounts(value.ID, value.OID)
		goodsInfo.GoodsType = m.GoodsTypeService.GetGoodsType(value.GoodsTypeID)
		goodsInfo.GoodsTypeChild = m.GoodsTypeService.GetGoodsTypeChild(value.GoodsTypeChildID)
		results = append(results, goodsInfo)
	}
	return results
}*/
