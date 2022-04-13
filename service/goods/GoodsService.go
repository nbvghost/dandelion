package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/tag"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/internal/repository"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/dandelion/service/pinyin"
	"github.com/nbvghost/gpa/params"
	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/tool/object"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
)

type GoodsService struct {
	model.BaseDao
	TimeSell      activity.TimeSellService
	Collage       activity.CollageService
	PinyinService pinyin.Service
}

func (service GoodsService) GetGoodsTypeByUri(OID types.PrimaryKey, GoodsTypeUri, GoodsTypeChildUri string) (model.GoodsType, model.GoodsTypeChild) {
	Orm := singleton.Orm()
	var item model.GoodsType
	var itemSub model.GoodsTypeChild

	itemMap := map[string]interface{}{"OID": OID, "Uri": GoodsTypeUri}
	Orm.Model(model.GoodsType{}).Where(itemMap).First(&item)

	itemSubMap := map[string]interface{}{
		"OID":         OID,
		"GoodsTypeID": item.ID,
		"Uri":         GoodsTypeChildUri,
	}
	Orm.Model(model.GoodsTypeChild{}).Where(itemSubMap).First(&itemSub)
	if itemSub.IsZero() {
		itemSub.Uri = "all"
	}
	return item, itemSub
}
func (service GoodsService) ListGoodsByOID(OID types.PrimaryKey) []model.GoodsType {
	Orm := singleton.Orm()
	var menus []model.GoodsType
	Orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID}).Find(&menus)
	return menus
}
func (service GoodsService) FindGoodsTags(OID types.PrimaryKey) ([]extends.Tag, error) {
	//SELECT unnest("Tags") as Tag,count("Tags") as Count FROM "Content" where  group by unnest("Tags");
	var tags []extends.Tag
	err := singleton.Orm().Model(model.Goods{}).Select(`unnest("Tags") as "Name",count("Tags") as "Count"`).Where(map[string]interface{}{
		"OID": OID,
	}).Where(`array_length("Tags",1)>0`).Group(`unnest("Tags")`).Find(&tags).Error
	tags = tag.CreateUri(tags)
	return tags, err
}
func (service GoodsService) FindGoodsByTag(OID types.PrimaryKey, tag extends.Tag, _pageIndex int, orders ...extends.Order) (pageIndex, pageSize int, total int64, list []model.Goods, err error) {
	//select * from "Content" where array_length("Tags",1) is null;
	db := singleton.Orm().Model(model.Goods{}).Where(`"OID"=?`, OID).
		Where(`array_length("Tags",1) is not null`).
		Where(`"Tags" @> array[?]`, tag.Name)

	db.Count(&total)

	for _, v := range orders {
		db.Order(fmt.Sprintf(`"%s" %s`, v.ColumnName, v.Method))
	}

	pageSize = 20

	err = db.Limit(pageSize).Offset(_pageIndex * pageSize).Find(&list).Error
	pageIndex = _pageIndex

	return
}
func (service GoodsService) PaginationGoods(OID, GoodsTypeID, GoodsTypeChildID types.PrimaryKey, pageIndex int) (int, int, int, []*model.Goods, error) {
	if GoodsTypeID == 0 {
		return repository.Goods.FindByOIDLimit(OID, params.NewLimit(pageIndex, 18))
	}
	if GoodsTypeChildID == 0 {
		return repository.Goods.FindByOIDAndGoodsTypeIDLimit(OID, GoodsTypeID, params.NewLimit(pageIndex, 18))
	}
	return repository.Goods.FindByOIDAndGoodsTypeIDAndGoodsTypeChildIDLimit(OID, GoodsTypeID, GoodsTypeChildID, params.NewLimit(pageIndex, 20))
}
func (service GoodsService) GetSpecification(ID types.PrimaryKey, target *model.Specification) error {
	Orm := singleton.Orm()
	err := service.Get(Orm, ID, &target)

	return err
}

/*func (service GoodsService) AddSpecification(context *gweb.Context) (r gweb.Result,err error) {
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
}*/
func (service GoodsService) DeleteSpecification(ID types.PrimaryKey) error {
	Orm := singleton.Orm()
	err := service.Delete(Orm, &model.Specification{}, ID)
	return err
}
func (service GoodsService) DeleteGoodsAttributes(ID types.PrimaryKey) error {

	return repository.GoodsAttributes.DeleteByID(ID).Err
}
func (service GoodsService) AddGoodsAttributes(goodsID, groupID types.PrimaryKey, name, value string) error {
	if goodsID == 0 || groupID == 0 {
		return errors.New(fmt.Sprintf("产品ID不能为空或组ID不能为空"))
	}
	if strings.EqualFold(name, "") || strings.EqualFold(value, "") {
		return nil
	}
	hasAttr, err := repository.GoodsAttributes.GetByGoodsIDAndName(goodsID, name)
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", name))
	}
	err = repository.GoodsAttributes.Save(&model.GoodsAttributes{
		GoodsID: goodsID,
		GroupID: groupID,
		Name:    name,
		Value:   value,
	})
	if glog.Error(err) {
		return err
	}
	return nil
}
func (service GoodsService) ListGoodsAttributesGroupByGoodsID(goodsID types.PrimaryKey) ([]*model.GoodsAttributesGroup, error) {

	return repository.GoodsAttributesGroup.FindByGoodsID(goodsID)
}
func (service GoodsService) GetGoodsAttributesGroup(ID types.PrimaryKey) types.IEntity {
	return repository.GoodsAttributesGroup.GetByID(ID)
}
func (service GoodsService) DeleteGoodsAttributesGroup(ID types.PrimaryKey) error {
	attrs, err := service.ListGoodsAttributesByGroupID(ID)
	if err != nil {
		return err
	}
	if len(attrs) > 0 {
		return errors.New(fmt.Sprintf("属性组包含子属性，无法删除"))
	}
	del := repository.GoodsAttributesGroup.DeleteByID(ID)
	return del.Err
}
func (service GoodsService) ListGoodsAttributesByGroupID(attributesGroupID types.PrimaryKey) ([]*model.GoodsAttributes, error) {
	return repository.GoodsAttributes.FindByGroupID(attributesGroupID)
}
func (service GoodsService) ChangeGoodsAttributesGroup(id types.PrimaryKey, groupName string) error {
	if id == 0 {
		return errors.New(fmt.Sprintf("ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil
	}
	hasAttr, err := repository.GoodsAttributesGroup.GetByName(groupName)
	if err != nil {
		return err
	}
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}

	update := repository.GoodsAttributesGroup.UpdateByID(id, map[string]interface{}{"Name": groupName})
	if glog.Error(update.Err) {
		return err
	}
	return nil
}
func (service GoodsService) AddGoodsAttributesGroup(goodsID types.PrimaryKey, groupName string) error {
	if goodsID == 0 {
		return errors.New(fmt.Sprintf("产品ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil
	}
	hasAttr, err := repository.GoodsAttributesGroup.GetByGoodsIDAndName(goodsID, groupName)
	if err != nil {
		return err
	}
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}
	err = repository.GoodsAttributesGroup.Save(&model.GoodsAttributesGroup{
		GoodsID: goodsID,
		Name:    groupName,
	})
	if glog.Error(err) {
		return err
	}
	return nil
}
func (service GoodsService) ChangeSpecification(context *gweb.Context) (r gweb.Result, err error) {
	Orm := singleton.Orm()
	//GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	GoodsID := object.ParseUint(context.PathParams["GoodsID"])
	item := &model.Specification{}
	err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}
	err = service.ChangeModel(Orm, types.PrimaryKey(GoodsID), item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
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

	var total uint
	for _, value := range specifications {

		value.GoodsID = goods.ID

		if goods.ID.IsZero() {
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

	err = tx.Save(goods).Error

	return err
}
func (service GoodsService) GetGoodsInfo(goods model.Goods) (*extends.GoodsInfo, error) {
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
	goodsInfo.Goods = goods
	goodsInfo.GoodsType = service.GetGoodsType(goods.GoodsTypeID)
	goodsInfo.GoodsTypeChild = service.GetGoodsTypeChild(goods.GoodsTypeChildID)
	goodsInfo.Discounts = make([]extends.Discount, 0)

	if timeSell.IsEnable() {
		//Favoured:=uint(util.Rounding45(float64(goods.Price)*(float64(timeSell.Discount)/float64(100)), 2))
		goodsInfo.Discounts = append(goodsInfo.Discounts, extends.Discount{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: "TimeSell", Discount: uint(timeSell.Discount)})
	} else {
		collage := service.Collage.GetCollageByGoodsID(goods.ID, goods.OID)
		if collage.ID != 0 && collage.TotalNum > 0 {
			goodsInfo.Discounts = append(goodsInfo.Discounts, extends.Discount{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint(collage.Discount)})
		}

	}

	var specifications []model.Specification
	err := service.FindWhere(Orm, &specifications, model.Specification{GoodsID: goods.ID})
	if err != nil {
		return nil, err
	}

	goodsInfo.Specifications = specifications

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
		goodsInfo.GoodsType = service.GetGoodsType(value.GoodsTypeID)
		goodsInfo.GoodsTypeChild = service.GetGoodsTypeChild(value.GoodsTypeChildID)
		results = append(results, goodsInfo)
	}

	return results
}
func (service GoodsService) GetGoods(DB *gorm.DB, context constrain.IContext, ID types.PrimaryKey) (*extends.GoodsInfo, error) {
	Orm := singleton.Orm()
	var goods model.Goods
	err := service.Get(Orm, ID, &goods)
	glog.Error(err)
	return service.GetGoodsInfo(goods)
}

func (service GoodsService) DeleteGoods(ID types.PrimaryKey) *result.ActionResult {
	Orm := singleton.Orm()
	tx := Orm.Begin()
	err := service.Delete(tx, &model.Goods{}, ID)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Where("GoodsID=?", ID).Delete(model.Specification{}).Error
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
func (service GoodsService) DeleteGoodsType(ID types.PrimaryKey) *result.ActionResult {
	Orm := singleton.Orm()
	tx := Orm.Begin()
	var gtcs []model.GoodsTypeChild
	tx.Where(&model.GoodsTypeChild{GoodsTypeID: ID}).Find(&gtcs) //Updates(map[string]interface{}{"GoodsTypeID": 0})

	var err error
	if len(gtcs) <= 0 {
		err = service.Delete(tx, &model.GoodsType{}, ID)
		if err != nil {
			tx.Rollback()
		}
	} else {
		return (&result.ActionResult{}).SmartError(err, "包含子类数据，不能删除", nil)
	}

	defer func() {
		if err == nil {
			tx.Commit()
		}
	}()
	return (&result.ActionResult{}).SmartError(err, "删除成功", nil)
}
func (service GoodsService) DeleteGoodsTypeChild(GoodsTypeChildID types.PrimaryKey) *result.ActionResult {
	Orm := singleton.Orm()
	tx := Orm.Begin()
	tx.Model(&model.Goods{GoodsTypeChildID: GoodsTypeChildID}).Updates(map[string]interface{}{"GoodsTypeChildID": 0})
	err := service.Delete(tx, &model.GoodsTypeChild{}, GoodsTypeChildID)
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
	err := service.DeleteWhere(DB, &model.TimeSellGoods{}, map[string]interface{}{
		"TimeSellHash": timesell.Hash,
		"GoodsID":      GoodsID,
	}) //Delete(DB, &model.TimeSellGoods{}, timesell.ID)
	glog.Error(err)
	return err
}
func (service GoodsService) DeleteCollageGoods(DB *gorm.DB, GoodsID types.PrimaryKey, OID types.PrimaryKey) error {
	timesell := service.Collage.GetCollageByGoodsID(GoodsID, OID)

	err := service.DeleteWhere(DB, &model.CollageGoods{}, map[string]interface{}{
		"CollageHash": timesell.Hash,
		"GoodsID":     GoodsID,
	}) //Delete(DB, &model.TimeSellGoods{}, timesell.ID)
	glog.Error(err)
	return err

	//err := service.Delete(DB, &model.Collage{}, timesell.ID)
	//glog.Error(err)
	//return err
}
func (service GoodsService) FindGoodsByTimeSellID(TimeSellID types.PrimaryKey) []model.Goods {
	Orm := singleton.Orm()

	var timesell model.TimeSell
	err := service.Get(Orm, TimeSellID, &timesell)
	glog.Error(err)

	var list []model.Goods

	//rows, _ := Orm.Raw("")

	//todo:没有写完整
	//err = service.FindWhere(Orm, &list, "ID=?", timesell.GoodsID)
	//glog.Error(err)
	return list
}
func (service GoodsService) FindGoodsByTimeSellHash(Hash string) []model.Goods {
	Orm := singleton.Orm()

	var GoodsIDs []uint
	Orm.Model(&model.TimeSell{}).Where("Hash=?", Hash).Pluck("GoodsID", &GoodsIDs)

	var list []model.Goods
	err := service.FindWhere(Orm, &list, "ID in (?)", GoodsIDs)
	glog.Error(err)
	return list
}
func (service GoodsService) FindGoodsByCollageHash(Hash string) []model.Goods {
	Orm := singleton.Orm()

	var GoodsIDs []uint
	Orm.Model(&model.Collage{}).Where("Hash=?", Hash).Pluck("GoodsID", &GoodsIDs)

	var list []model.Goods
	err := service.FindWhere(Orm, &list, "ID in (?)", GoodsIDs)
	glog.Error(err)
	return list
}
func (service GoodsService) FindGoodsByOrganizationIDAndGoodsID(OrganizationID types.PrimaryKey, GoodsID types.PrimaryKey) model.Goods {
	var Goods model.Goods
	singleton.Orm().Model(&model.Goods{}).Where("ID=? and OID=?", GoodsID, OrganizationID).First(&Goods)
	return Goods
}
func (service GoodsService) FindGoodsByTitle(Title string) model.Goods {
	var Goods model.Goods
	singleton.Orm().Model(&model.Goods{}).Where("Title=?", Title).First(&Goods)
	return Goods
}
func (service GoodsService) FindGoodsLikeMark(Mark string) model.Goods {
	var Goods model.Goods
	singleton.Orm().Model(&model.Goods{}).Where("Mark like ?", "%"+Mark+"%").First(&Goods)
	return Goods
}
func (service GoodsService) AllList() []model.Goods {

	Orm := singleton.Orm()

	var result []model.Goods

	db := Orm.Model(&model.Goods{}).Order("CreatedAt desc") //.Limit(10)

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

type TopGoodsTypeChild struct {
	Name             string `gorm:"column:Name"`
	Image            string `gorm:"column:Image"`
	GoodsTypeChildID uint   `gorm:"column:GoodsTypeChildID"`
	Price            uint   `gorm:"column:Price"`
}

func (service GoodsService) GetTopGoodsTypeChild(DB *gorm.DB, Num uint) []TopGoodsTypeChild {
	list := make([]TopGoodsTypeChild, 0)
	//SELECT gtc.Name,gtc.Image,gtc.ID AS GoodsTypeChildID,gtc.GoodsTypeID AS GoodsTypeID,MIN(g.Price) FROM Goods AS g LEFT JOIN GoodsTypeChild AS gtc ON (gtc.GoodsTypeID=g.GoodsTypeID AND gtc.ID=g.GoodsTypeChildID) GROUP BY g.GoodsTypeID;
	rows, err := DB.Raw("SELECT gtc.Name as Name,gtc.Image as Image,gtc.ID AS GoodsTypeChildID,gtc.GoodsTypeID AS GoodsTypeID,MIN(g.Price) as Price FROM Goods AS g LEFT JOIN GoodsTypeChild AS gtc ON (gtc.GoodsTypeID=g.GoodsTypeID AND gtc.ID=g.GoodsTypeChildID) GROUP BY g.GoodsTypeID limit ?", Num).Rows()
	if glog.Error(err) {
		return list
	}
	for rows.Next() {
		var result TopGoodsTypeChild
		DB.ScanRows(rows, &result)
		list = append(list, result)
	}
	return list
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
func (service GoodsService) HotListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID, GoodsTypeChildID, Num uint) []model.Goods {

	Orm := singleton.Orm()

	var result []model.Goods

	db := Orm.Model(&model.Goods{}).Where("GoodsTypeID=? and GoodsTypeChildID=?", GoodsTypeID, GoodsTypeChildID).Order("CountSale desc").Limit(int(Num))

	db.Find(&result)

	return result

}
func (service GoodsService) NewListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID, GoodsTypeChildID, Num uint) []model.Goods {

	Orm := singleton.Orm()

	var result []model.Goods

	db := Orm.Model(&model.Goods{}).Where("GoodsTypeID=? and GoodsTypeChildID=?", GoodsTypeID, GoodsTypeChildID).Order("CreatedAt desc").Limit(int(Num))

	db.Find(&result)

	return result

}
func (service GoodsService) HotSaleList(OID types.PrimaryKey, count uint) []model.Goods {
	Orm := singleton.Orm()
	var result []model.Goods
	db := Orm.Model(&model.Goods{}).Where(map[string]interface{}{"OID": OID}).Order(`"CountSale" desc`).Limit(int(count))
	db.Find(&result)
	return result
}
func (service GoodsService) HotViewList(OID types.PrimaryKey, count uint) []model.Goods {
	Orm := singleton.Orm()
	var result []model.Goods
	db := Orm.Model(&model.Goods{}).Where(map[string]interface{}{"OID": OID}).Order(`"CountView" desc`).Limit(int(count))
	db.Find(&result)
	return result
}
func (service GoodsService) GetGoodsType(ID types.PrimaryKey) model.GoodsType {
	Orm := singleton.Orm()
	var result model.GoodsType
	Orm.Model(&model.GoodsType{}).Where(map[string]interface{}{"ID": ID}).First(&result)
	return result
}
func (service GoodsService) GetGoodsTypeChild(ID types.PrimaryKey) model.GoodsTypeChild {
	Orm := singleton.Orm()
	var result model.GoodsTypeChild
	Orm.Model(&model.GoodsTypeChild{}).Where(map[string]interface{}{"ID": ID}).First(&result)
	return result
}
func (service GoodsService) AllGoodsType() []extends.AllGoodsType {
	gts := service.ListAllGoodsType()

	list := make([]extends.AllGoodsType, 0)
	for i := range gts {
		subs := service.ListAllGoodsTypeChild(gts[i].ID)

		goodsList := service.ListGoodsByGoodsTypeID(gts[i].ID)
		var maxPrice uint = 0
		var minPrice uint = 0
		for ii := range goodsList {
			if ii == 0 {
				minPrice = goodsList[ii].Price
			}
			if goodsList[ii].Price > maxPrice {
				maxPrice = goodsList[ii].Price
			}
			if goodsList[ii].Price < minPrice {
				minPrice = goodsList[ii].Price
			}
		}
		list = append(list, extends.AllGoodsType{
			GoodsType: gts[i],
			Sub:       subs,
			MaxPrice:  uint(maxPrice),
			MinPrice:  uint(minPrice),
		})
	}
	return list
}
func (service GoodsService) ListAllGoodsType() []model.GoodsType {
	/*Orm := singleton.Orm()
	var gts []model.GoodsType
	service.FindAll(Orm, &gts)
	return gts*/
	Orm := singleton.Orm()
	var gts []model.GoodsType
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Group("GoodsTypeID").Pluck("GoodsTypeID", &gtsIDs)
	Orm.Model(&model.GoodsType{}).Where("ID in (?)", gtsIDs).Find(&gts)
	return gts
}
func (service GoodsService) ListGoodsType(OID uint) []model.GoodsType {
	/*Orm := singleton.Orm()
	var gts []model.GoodsType
	service.FindAllByOID(Orm,&gts,OID)
	return gts*/
	Orm := singleton.Orm()
	var gts []model.GoodsType
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Where("OID=?", OID).Group("GoodsTypeID").Pluck("GoodsTypeID", &gtsIDs)
	Orm.Model(&model.GoodsType{}).Where("ID in (?)", gtsIDs).Find(&gts)
	return gts
}
func (service GoodsService) ListGoodsTypeByOIDForAdmin(OID uint) []model.GoodsType {
	Orm := singleton.Orm()
	var gts []model.GoodsType
	Orm.Model(&model.GoodsType{}).Where("OID=?", OID).Find(&gts)
	return gts
}
func (service GoodsService) ListGoodsTypeForAdmin() []model.GoodsType {
	Orm := singleton.Orm()
	var gts []model.GoodsType
	Orm.Model(&model.GoodsType{}).Find(&gts)
	return gts
}

func (service GoodsService) ListGoodsTypeChildAll(OID types.PrimaryKey) []model.GoodsTypeChild {
	Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Where("OID=?", OID).Group("GoodsTypeChildID").Pluck("GoodsTypeChildID", &gtsIDs)
	Orm.Model(&model.GoodsTypeChild{}).Where("ID in (?)", gtsIDs).Find(&gts)
	return gts
}
func (service GoodsService) ListAllGoodsTypeChild(GoodsTypeID types.PrimaryKey) []model.GoodsTypeChild {
	/*Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	service.FindWhere(Orm, &gts, model.GoodsTypeChild{GoodsTypeID: GoodsTypeID})
	return gts*/
	Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	Orm.Model(&model.GoodsTypeChild{}).Where(map[string]interface{}{
		"GoodsTypeID": GoodsTypeID,
	}).Find(&gts)
	return gts
}

func (service GoodsService) ListGoodsTypeChild(GoodsTypeID types.PrimaryKey) []model.GoodsTypeChild {
	/*Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	service.FindWhere(Orm, &gts, model.GoodsTypeChild{GoodsTypeID: GoodsTypeID})
	return gts*/
	Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Where("GoodsTypeID=?", GoodsTypeID).Group("GoodsTypeChildID").Pluck("GoodsTypeChildID", &gtsIDs)
	Orm.Model(&model.GoodsTypeChild{}).Where("ID in (?)", gtsIDs).Find(&gts)
	return gts
}

func (service GoodsService) ListGoodsChildByGoodsTypeID(GoodsTypeID, GoodsTypeChildID types.PrimaryKey) []model.Goods {
	Orm := singleton.Orm()
	var gts []model.Goods
	service.FindWhere(Orm, &gts, model.Goods{GoodsTypeID: GoodsTypeID, GoodsTypeChildID: GoodsTypeChildID})
	return gts
}
func (service GoodsService) ListGoodsByGoodsTypeID(GoodsTypeID types.PrimaryKey) []model.Goods {
	Orm := singleton.Orm()
	var gts []model.Goods
	service.FindWhere(Orm, &gts, model.Goods{GoodsTypeID: GoodsTypeID})
	return gts
}
func (service GoodsService) AddGoodsTypeByNameByChild(name string, childName string) (goodsType model.GoodsType, goodsTypeChild model.GoodsTypeChild) {
	Orm := singleton.Orm()

	var gt model.GoodsType
	var gtc model.GoodsTypeChild

	err := Orm.Model(&model.GoodsType{}).Where("Name=?", name).First(&gt).Error
	if gorm.ErrRecordNotFound == err {
		gt.Name = name
		service.Save(Orm, &gt)
	}

	err = Orm.Model(&model.GoodsTypeChild{}).Where("Name=? and GoodsTypeID=?", childName, gt.ID).First(&gtc).Error
	if gorm.ErrRecordNotFound == err {
		gtc.Name = childName
		gtc.GoodsTypeID = gt.ID
		service.Save(Orm, &gtc)
	}

	return gt, gtc
}

func (service GoodsService) RecommendGoods(OID, GoodsID, GoodsTypeID, GoodsTypeChildID types.PrimaryKey) []model.Goods {
	var contentList []model.Goods
	singleton.Orm().
		Model(&model.Goods{}).
		Where("OID=? and ID<>? and (GoodsTypeID=? or GoodsTypeChildID=?)", OID, GoodsID, GoodsTypeID, GoodsTypeChildID).
		Order("RAND()").
		Limit(6).
		Find(&contentList)
	return contentList
}
func (service GoodsService) ListGoodsByType(OID, GoodsTypeID, GoodsTypeChildID types.PrimaryKey) []model.Goods {

	var contentList []model.Goods

	if GoodsTypeID == 0 {
		singleton.Orm().Model(&model.Goods{}).Where(`"OID"=?`, OID).
			Order(`"CreatedAt" desc`).Order(`"ID" desc`).Find(&contentList)
		return contentList
	}

	if GoodsTypeChildID > 0 {
		singleton.Orm().Model(&model.Goods{}).Where(`"OID"=? and "GoodsTypeID"=? and "GoodsTypeChildID"=?`, OID, GoodsTypeID, GoodsTypeChildID).
			Order(`"CreatedAt" desc`).Order(`"ID" desc`).Find(&contentList)
		return contentList
	} else {
		singleton.Orm().Model(&model.Goods{}).Where(`"OID"=? and "GoodsTypeID"=?`, OID, GoodsTypeID).
			Order(`"CreatedAt" desc`).Order(`"ID" desc`).Find(&contentList)
		return contentList
	}

}
func (service GoodsService) GetGoodsTypeData(OID types.PrimaryKey) *model.GoodsTypeData {

	goodsTypeData := &model.GoodsTypeData{}

	rows, err := singleton.Orm().Raw(`SELECT gt.*,gtc.* FROM GoodsTypeChild AS gtc LEFT JOIN GoodsType as gt ON (gt.ID=gtc.GoodsTypeID) WHERE OID=?`, OID).Rows()
	if glog.Error(err) {
		return goodsTypeData
	}

	defer rows.Close()

	for rows.Next() {
		var item model.GoodsTypeGoodsTypeChild
		err := singleton.Orm().ScanRows(rows, &item)
		glog.Error(err)

		goodsTypeItem := goodsTypeData.Get(item.GoodsType.ID)
		if goodsTypeItem.Item.ID == 0 {

			goodsTypeData.List = append(goodsTypeData.List, &model.GoodsTypeItem{Item: &item.GoodsType, SubType: []*model.GoodsTypeItemSub{{Item: &item.GoodsTypeChild, SubType: []*model.GoodsTypeItemSub{}}}})

		} else {

			goodsTypeItem.SubType = append(goodsTypeItem.SubType, &model.GoodsTypeItemSub{Item: &item.GoodsTypeChild, SubType: []*model.GoodsTypeItemSub{}})

		}

	}

	return goodsTypeData

}
func (service GoodsService) getGoodsTypeByName(orm *gorm.DB, OID types.PrimaryKey, name string) (model.GoodsType, error) {
	var gt model.GoodsType
	err := orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Name": name}).First(&gt).Error
	return gt, err
}
func (service GoodsService) getGoodsTypeByUri(orm *gorm.DB, OID types.PrimaryKey, uri string) (model.GoodsType, error) {
	var gt model.GoodsType
	err := orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Uri": uri}).First(&gt).Error
	return gt, err
}
func (service GoodsService) AddGoodsType(OID types.PrimaryKey, name string) error {
	orm := singleton.Orm()
	gt, _ := service.getGoodsTypeByName(orm, OID, name)
	if !gt.IsZero() {
		return errors.Errorf("重复的名字:%s", name)
	}

	uri := service.PinyinService.AutoDetectUri(name)
	gt, _ = service.getGoodsTypeByUri(orm, OID, uri)
	if !gt.IsZero() {
		gt.Uri = fmt.Sprintf("%s-%d", gt.Uri, time.Now().Unix())
	}
	gt.OID = OID
	gt.Name = name
	gt.Uri = uri
	return orm.Model(model.GoodsType{}).Create(&gt).Error
}
func (service GoodsService) ChangeGoodsType(OID, ID types.PrimaryKey, name string) error {
	orm := singleton.Orm()
	gt, _ := service.getGoodsTypeByName(orm, OID, name)
	if gt.ID == ID {
		return nil
	}
	if !gt.IsZero() {
		return errors.Errorf("重复的名字:%s", name)
	}

	uri := service.PinyinService.AutoDetectUri(name)
	gt, _ = service.getGoodsTypeByUri(orm, OID, uri)
	if !gt.IsZero() {
		gt.Uri = fmt.Sprintf("%s-%d", gt.Uri, time.Now().Unix())
	}
	gt.Name = name
	gt.Uri = uri
	return orm.Model(model.GoodsType{}).Where(`"ID"=?`, ID).Updates(map[string]interface{}{
		"Name": gt.Name,
		"Uri":  gt.Uri,
	}).Error
}
func (service GoodsService) getGoodsTypeChildByName(orm *gorm.DB, OID, GoodsTypeID types.PrimaryKey, name string) (model.GoodsTypeChild, error) {
	var gt model.GoodsTypeChild
	err := orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "GoodsTypeID": GoodsTypeID, "Name": name}).First(&gt).Error
	return gt, err
}
func (service GoodsService) getGoodsTypeChildByUri(orm *gorm.DB, OID, GoodsTypeID types.PrimaryKey, uri string) (model.GoodsTypeChild, error) {
	var gt model.GoodsTypeChild
	err := orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "GoodsTypeID": GoodsTypeID, "Uri": uri}).First(&gt).Error
	return gt, err
}
func (service GoodsService) AddGoodsTypeChild(OID, GoodsTypeID types.PrimaryKey, name, image string) error {
	if GoodsTypeID == 0 {
		return errors.Errorf("没有指定父类ID")
	}
	orm := singleton.Orm()
	gtc := service.GetGoodsType(GoodsTypeID)
	if gtc.IsZero() {
		return errors.Errorf("不存在父类:%d", GoodsTypeID)
	}

	gt, _ := service.getGoodsTypeChildByName(orm, OID, GoodsTypeID, name)
	if !gt.IsZero() {
		return errors.Errorf("重复的名字:%s", name)
	}

	uri := service.PinyinService.AutoDetectUri(name)
	gt, _ = service.getGoodsTypeChildByUri(orm, OID, GoodsTypeID, uri)
	if !gt.IsZero() {
		gt.Uri = fmt.Sprintf("%s-%d", gt.Uri, time.Now().Unix())
	}
	gt.OID = OID
	gt.Name = name
	gt.Uri = uri
	gt.Image = image
	gt.GoodsTypeID = GoodsTypeID
	return orm.Model(model.GoodsTypeChild{}).Create(&gt).Error
}
func (service GoodsService) ChangeGoodsTypeChild(OID, ID types.PrimaryKey, name, image string) error {
	orm := singleton.Orm()
	gtc := service.GetGoodsTypeChild(ID)
	if gtc.IsZero() {
		return errors.Errorf("记录不存在")
	}

	gt, _ := service.getGoodsTypeChildByName(orm, OID, gtc.GoodsTypeID, name)
	if gt.ID == ID {
		return nil
	}
	if !gt.IsZero() {
		return errors.Errorf("重复的名字:%s", name)
	}

	uri := service.PinyinService.AutoDetectUri(name)
	gt, _ = service.getGoodsTypeChildByUri(orm, OID, gtc.GoodsTypeID, uri)
	if !gt.IsZero() {
		gt.Uri = fmt.Sprintf("%s-%d", gt.Uri, time.Now().Unix())
	}
	gt.Name = name
	gt.Uri = uri
	gt.Image = image
	return orm.Model(model.GoodsTypeChild{}).Where(`"ID"=?`, ID).Updates(map[string]interface{}{
		"Name":  gt.Name,
		"Uri":   gt.Uri,
		"Image": gt.Image,
	}).Error
}
