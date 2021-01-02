package goods

import (
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/activity"
	"strconv"

	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
)

type GoodsService struct {
	dao.BaseDao
	TimeSell activity.TimeSellService
	Collage  activity.CollageService
}

func (service GoodsService) GetSpecification(ID uint64, target *dao.Specification) error {
	Orm := dao.Orm()
	err := service.Get(Orm, ID, &target)

	return err
}

/*func (service GoodsService) AddSpecification(context *gweb.Context) gweb.Result {
	item := &dao.Specification{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	}
	err = service.Add(Orm, item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}
}
func (service GoodsService) ListSpecification(context *gweb.Context) gweb.Result {
	GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	var dts []dao.Specification
	service.FindWhere(Orm, &dts, dao.Specification{GoodsID: GoodsID})
	return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionOK, Message: "OK", Data: dts}}
}*/
func (service GoodsService) DeleteSpecification(ID uint64) error {
	Orm := dao.Orm()
	err := service.Delete(Orm, &dao.Specification{}, ID)
	return err
}
func (service GoodsService) ChangeSpecification(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	item := &dao.Specification{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	}
	err = service.ChangeModel(Orm, GoodsID, item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}
}
func (service GoodsService) SaveGoods(goods dao.Goods, specifications []dao.Specification) error {
	Orm := dao.Orm()
	var err error
	tx := Orm.Begin()

	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if tx.NewRecord(&goods) {
		err = tx.Create(&goods).Error
	} else {
		//err = tx.Save(goods).Error
		err = tx.Model(&goods).Updates(goods).Error
	}

	if err != nil {
		return err
	}

	//添加或修改的时候不删除规格
	/*err=service.UnscopedDeleteWhere(tx,&dao.Specification{},"GoodsID=?",goods.ID)
	if err!=nil{
		return err
	}*/

	var total uint
	for _, value := range specifications {

		value.GoodsID = goods.ID

		if tx.NewRecord(&goods) {
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

	err = tx.Save(&goods).Error

	return err
}
func (service GoodsService) GetGoodsInfo(goods dao.Goods) dao.GoodsInfo {
	//Orm := dao.Orm()

	//user := dao.User{}
	//service.User.Get(Orm, UserID, &user)
	//brokerageProvisoConf := service.Configuration.GetConfiguration(play.ConfigurationKey_BrokerageProviso)
	//brokerageProvisoConfV, _ := strconv.ParseUint(brokerageProvisoConf.V, 10, 64)
	//vipdiscountConf := service.Configuration.GetConfiguration(play.ConfigurationKey_VIPDiscount)
	//VIPDiscount, _ := strconv.ParseUint(vipdiscountConf.V, 10, 64)
	timeSell := service.TimeSell.GetTimeSellByGoodsID(goods.ID, goods.OID)
	goodsInfo := dao.GoodsInfo{}
	goodsInfo.Goods = goods
	goodsInfo.Discounts = make([]dao.Discount, 0)

	if timeSell.IsEnable() {
		//Favoured:=uint64(util.Rounding45(float64(goods.Price)*(float64(timeSell.Discount)/float64(100)), 2))
		goodsInfo.Discounts = append(goodsInfo.Discounts, dao.Discount{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: "TimeSell", Discount: uint64(timeSell.Discount)})
	} else {
		collage := service.Collage.GetCollageByGoodsID(goods.ID, goods.OID)
		if collage.ID != 0 && collage.TotalNum > 0 {
			goodsInfo.Discounts = append(goodsInfo.Discounts, dao.Discount{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint64(collage.Discount)})
		}

	}
	return goodsInfo
}
func (service GoodsService) GetGoods(DB *gorm.DB, ID uint64) dao.GoodsInfo {
	Orm := dao.Orm()
	var goods dao.Goods
	err := service.Get(Orm, ID, &goods)
	glog.Trace(err)

	var specifications []dao.Specification
	err = service.FindWhere(Orm, &specifications, dao.Specification{GoodsID: ID})
	glog.Trace(err)

	goodsInfo := service.GetGoodsInfo(goods)
	goodsInfo.Specifications = specifications

	/*var mtimeSell dao.TimeSell
	err=TimeSellService{}.Get(Orm,goods.TimeSellID,&mtimeSell)
	glog.Trace(err)
	if mtimeSell.IsEnable(){
		timeSell = mtimeSell
	}else {
		timeSell = dao.TimeSell{}
	}*/

	return goodsInfo
	//return DB.Model(target).Related(&dao.Specification{}).Where("ID=?", ID).First(target).Error
	/*Orm := dao.Orm()
	err := service.Get(Orm, ID, &goods)
	glog.Trace(err)

	err = service.FindWhere(Orm, &specifications, dao.Specification{GoodsID: ID})
	glog.Trace(err)

	var mtimeSell dao.TimeSell
	err = TimeSellService{}.Get(Orm, goods.TimeSellID, &mtimeSell)
	glog.Trace(err)
	if mtimeSell.IsEnable() {
		timeSell = mtimeSell
	} else {
		timeSell = dao.TimeSell{}
	}

	return*/
	//return DB.Model(target).Related(&dao.Specification{}).Where("ID=?", ID).First(target).Error
}

func (service GoodsService) DeleteGoods(ID uint64) *result.ActionResult {
	Orm := dao.Orm()
	tx := Orm.Begin()
	err := service.Delete(tx, &dao.Goods{}, ID)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Where("GoodsID=?", ID).Delete(dao.Specification{}).Error
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
func (service GoodsService) DeleteGoodsType(ID uint64) *result.ActionResult {
	Orm := dao.Orm()
	tx := Orm.Begin()
	var gtcs []dao.GoodsTypeChild
	tx.Where(&dao.GoodsTypeChild{GoodsTypeID: ID}).Find(&gtcs) //Updates(map[string]interface{}{"GoodsTypeID": 0})

	var err error
	if len(gtcs) <= 0 {
		err = service.Delete(tx, &dao.GoodsType{}, ID)
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
func (service GoodsService) DeleteGoodsTypeChild(GoodsTypeChildID uint64) *result.ActionResult {
	Orm := dao.Orm()
	tx := Orm.Begin()
	tx.Model(&dao.Goods{GoodsTypeChildID: GoodsTypeChildID}).Updates(map[string]interface{}{"GoodsTypeChildID": 0})
	err := service.Delete(tx, &dao.GoodsTypeChild{}, GoodsTypeChildID)
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

func (service GoodsService) DeleteTimeSellGoods(DB *gorm.DB, GoodsID uint64, OID uint64) error {
	timesell := service.TimeSell.GetTimeSellByGoodsID(GoodsID, OID)
	err := service.DeleteWhere(DB, &dao.TimeSellGoods{}, "TimeSellHash=? and GoodsID=?", timesell.Hash, GoodsID) //Delete(DB, &dao.TimeSellGoods{}, timesell.ID)
	glog.Error(err)
	return err
}
func (service GoodsService) DeleteCollageGoods(DB *gorm.DB, GoodsID uint64, OID uint64) error {
	timesell := service.Collage.GetCollageByGoodsID(GoodsID, OID)

	err := service.DeleteWhere(DB, &dao.CollageGoods{}, "CollageHash=? and GoodsID=?", timesell.Hash, GoodsID) //Delete(DB, &dao.TimeSellGoods{}, timesell.ID)
	glog.Error(err)
	return err

	//err := service.Delete(DB, &dao.Collage{}, timesell.ID)
	//glog.Error(err)
	//return err
}
func (service GoodsService) FindGoodsByTimeSellID(TimeSellID uint64) []dao.Goods {
	Orm := dao.Orm()

	var timesell dao.TimeSell
	err := service.Get(Orm, TimeSellID, &timesell)
	glog.Error(err)

	var list []dao.Goods

	//rows, _ := Orm.Raw("")

	//todo:
	//err = service.FindWhere(Orm, &list, "ID=?", timesell.GoodsID)
	//glog.Error(err)
	return list
}
func (service GoodsService) FindGoodsByTimeSellHash(Hash string) []dao.Goods {
	Orm := dao.Orm()

	var GoodsIDs []uint64
	Orm.Model(&dao.TimeSell{}).Where("Hash=?", Hash).Pluck("GoodsID", &GoodsIDs)

	var list []dao.Goods
	err := service.FindWhere(Orm, &list, "ID in (?)", GoodsIDs)
	glog.Error(err)
	return list
}
func (service GoodsService) FindGoodsByCollageHash(Hash string) []dao.Goods {
	Orm := dao.Orm()

	var GoodsIDs []uint64
	Orm.Model(&dao.Collage{}).Where("Hash=?", Hash).Pluck("GoodsID", &GoodsIDs)

	var list []dao.Goods
	err := service.FindWhere(Orm, &list, "ID in (?)", GoodsIDs)
	glog.Error(err)
	return list
}
func (service GoodsService) FindGoodsByOrganizationIDAndGoodsID(OrganizationID uint64, GoodsID uint64) dao.Goods {
	var Goods dao.Goods
	dao.Orm().Model(&dao.Goods{}).Where("ID=? and OID=?", GoodsID, OrganizationID).First(&Goods)
	return Goods
}
func (service GoodsService) FindGoodsByTitle(Title string) dao.Goods {
	var Goods dao.Goods
	dao.Orm().Model(&dao.Goods{}).Where("Title=?", Title).First(&Goods)
	return Goods
}
func (service GoodsService) FindGoodsLikeMark(Mark string) dao.Goods {
	var Goods dao.Goods
	dao.Orm().Model(&dao.Goods{}).Where("Mark like ?", "%"+Mark+"%").First(&Goods)
	return Goods
}
func (service GoodsService) AllList() []dao.Goods {

	Orm := dao.Orm()

	var result []dao.Goods

	db := Orm.Model(&dao.Goods{}).Order("CreatedAt desc") //.Limit(10)

	db.Find(&result)

	return result

}
func (service GoodsService) GetDiscounts(GoodsID, OID uint64) []dao.Discount {
	discounts := make([]dao.Discount, 0)
	timeSell := service.TimeSell.GetTimeSellByGoodsID(GoodsID, OID)
	if timeSell.IsEnable() {
		//Favoured:=uint64(util.Rounding45(float64(value.Price)*(float64(timeSell.Discount)/float64(100)), 2))
		discounts = append(discounts, dao.Discount{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: "TimeSell", Discount: uint64(timeSell.Discount)})
	} else {
		collage := service.Collage.GetCollageByGoodsID(GoodsID, OID)
		if collage.ID != 0 && collage.TotalNum > 0 {
			discounts = append(discounts, dao.Discount{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint64(collage.Discount)})
		}
	}
	return discounts
}
func (service GoodsService) GetGoodsInfoList(UserID uint64, goodsList []dao.Goods) []dao.GoodsInfo {

	var results = make([]dao.GoodsInfo, 0)

	for _, value := range goodsList {
		goodsInfo := dao.GoodsInfo{}
		goodsInfo.Goods = value
		goodsInfo.Discounts = service.GetDiscounts(value.ID, value.OID)
		results = append(results, goodsInfo)
	}

	return results
}

type TopGoodsTypeChild struct {
	Name             string `gorm:"column:Name"`
	Image            string `gorm:"column:Image"`
	GoodsTypeChildID uint64 `gorm:"column:GoodsTypeChildID"`
	Price            uint64 `gorm:"column:Price"`
}

func (service GoodsService) GetTopGoodsTypeChild(DB *gorm.DB, Num uint64) []TopGoodsTypeChild {
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
func (service GoodsService) GoodsList(UserID uint64, SqlOrder string, Index int, where interface{}, args ...interface{}) []dao.GoodsInfo {
	Orm := dao.Orm()
	//var goodsList []dao.Goods
	//db := Orm.Model(&dao.Goods{}).Order("CountSale desc").Limit(10)
	//db.Find(&result)
	pager := service.FindWherePaging(Orm, SqlOrder, dao.Goods{}, Index, where, args)

	return service.GetGoodsInfoList(UserID, pager.Data.([]dao.Goods))
}
func (service GoodsService) HotListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID, GoodsTypeChildID, Num uint64) []dao.Goods {

	Orm := dao.Orm()

	var result []dao.Goods

	db := Orm.Model(&dao.Goods{}).Where("GoodsTypeID=? and GoodsTypeChildID=?", GoodsTypeID, GoodsTypeChildID).Order("CountSale desc").Limit(Num)

	db.Find(&result)

	return result

}
func (service GoodsService) NewListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID, GoodsTypeChildID, Num uint64) []dao.Goods {

	Orm := dao.Orm()

	var result []dao.Goods

	db := Orm.Model(&dao.Goods{}).Where("GoodsTypeID=? and GoodsTypeChildID=?", GoodsTypeID, GoodsTypeChildID).Order("CreatedAt desc").Limit(Num)

	db.Find(&result)

	return result

}
func (service GoodsService) HotList(count uint64) []dao.Goods {

	Orm := dao.Orm()

	var result []dao.Goods

	db := Orm.Model(&dao.Goods{}).Order("CountSale desc").Limit(count)

	db.Find(&result)

	return result

}
func (service GoodsService) ListAllGoodsType() []dao.GoodsType {
	/*Orm := dao.Orm()
	var gts []dao.GoodsType
	service.FindAll(Orm, &gts)
	return gts*/
	Orm := dao.Orm()
	var gts []dao.GoodsType
	var gtsIDs []uint64
	//service.FindWhere(Orm, &gts, dao.GoodsTypeChild{})
	Orm.Model(&dao.Goods{}).Group("GoodsTypeID").Pluck("GoodsTypeID", &gtsIDs)
	Orm.Model(&dao.GoodsType{}).Where("ID in (?)", gtsIDs).Find(&gts)
	return gts
}
func (service GoodsService) ListGoodsType(OID uint64) []dao.GoodsType {
	/*Orm := dao.Orm()
	var gts []dao.GoodsType
	service.FindAllByOID(Orm,&gts,OID)
	return gts*/
	Orm := dao.Orm()
	var gts []dao.GoodsType
	var gtsIDs []uint64
	//service.FindWhere(Orm, &gts, dao.GoodsTypeChild{})
	Orm.Model(&dao.Goods{}).Where("OID=?", OID).Group("GoodsTypeID").Pluck("GoodsTypeID", &gtsIDs)
	Orm.Model(&dao.GoodsType{}).Where("ID in (?)", gtsIDs).Find(&gts)
	return gts
}
func (service GoodsService) ListGoodsTypeByOIDForAdmin(OID uint64) []dao.GoodsType {
	Orm := dao.Orm()
	var gts []dao.GoodsType
	Orm.Model(&dao.GoodsType{}).Where("OID=?", OID).Find(&gts)
	return gts
}
func (service GoodsService) ListGoodsTypeForAdmin() []dao.GoodsType {
	Orm := dao.Orm()
	var gts []dao.GoodsType
	Orm.Model(&dao.GoodsType{}).Find(&gts)
	return gts
}

func (service GoodsService) ListGoodsTypeChildAll(OID uint64) []dao.GoodsTypeChild {
	Orm := dao.Orm()
	var gts []dao.GoodsTypeChild
	var gtsIDs []uint64
	//service.FindWhere(Orm, &gts, dao.GoodsTypeChild{})
	Orm.Model(&dao.Goods{}).Where("OID=?", OID).Group("GoodsTypeChildID").Pluck("GoodsTypeChildID", &gtsIDs)
	Orm.Model(&dao.GoodsTypeChild{}).Where("ID in (?)", gtsIDs).Find(&gts)
	return gts
}
func (service GoodsService) ListAllGoodsTypeChild(GoodsTypeID uint64) []dao.GoodsTypeChild {
	/*Orm := dao.Orm()
	var gts []dao.GoodsTypeChild
	service.FindWhere(Orm, &gts, dao.GoodsTypeChild{GoodsTypeID: GoodsTypeID})
	return gts*/
	Orm := dao.Orm()
	var gts []dao.GoodsTypeChild
	Orm.Model(&dao.GoodsTypeChild{}).Where("GoodsTypeID=?", GoodsTypeID).Find(&gts)
	return gts
}

func (service GoodsService) ListGoodsTypeChild(GoodsTypeID uint64) []dao.GoodsTypeChild {
	/*Orm := dao.Orm()
	var gts []dao.GoodsTypeChild
	service.FindWhere(Orm, &gts, dao.GoodsTypeChild{GoodsTypeID: GoodsTypeID})
	return gts*/
	Orm := dao.Orm()
	var gts []dao.GoodsTypeChild
	var gtsIDs []uint64
	//service.FindWhere(Orm, &gts, dao.GoodsTypeChild{})
	Orm.Model(&dao.Goods{}).Where("GoodsTypeID=?", GoodsTypeID).Group("GoodsTypeChildID").Pluck("GoodsTypeChildID", &gtsIDs)
	Orm.Model(&dao.GoodsTypeChild{}).Where("ID in (?)", gtsIDs).Find(&gts)
	return gts
}

func (service GoodsService) ListGoodsChildByGoodsTypeID(GoodsTypeID, GoodsTypeChildID uint64) []dao.Goods {
	Orm := dao.Orm()
	var gts []dao.Goods
	service.FindWhere(Orm, &gts, dao.Goods{GoodsTypeID: GoodsTypeID, GoodsTypeChildID: GoodsTypeChildID})
	return gts
}
func (service GoodsService) ListGoodsByGoodsTypeID(GoodsTypeID uint64) []dao.Goods {
	Orm := dao.Orm()
	var gts []dao.Goods
	service.FindWhere(Orm, &gts, dao.Goods{GoodsTypeID: GoodsTypeID})
	return gts
}
func (service GoodsService) AddGoodsTypeByNameByChild(name string, childName string) (goodsType dao.GoodsType, goodsTypeChild dao.GoodsTypeChild) {
	Orm := dao.Orm()

	var gt dao.GoodsType
	var gtc dao.GoodsTypeChild

	err := Orm.Model(&dao.GoodsType{}).Where("Name=?", name).First(&gt).Error
	if gorm.IsRecordNotFoundError(err) {
		gt.Name = name
		service.Save(Orm, &gt)
	}

	err = Orm.Model(&dao.GoodsTypeChild{}).Where("Name=? and GoodsTypeID=?", childName, gt.ID).First(&gtc).Error
	if gorm.IsRecordNotFoundError(err) {
		gtc.Name = childName
		gtc.GoodsTypeID = gt.ID
		service.Save(Orm, &gtc)
	}

	return gt, gtc
}

func (service GoodsService) RecommendGoods(OID, GoodsID, GoodsTypeID, GoodsTypeChildID uint64) []dao.Goods {
	var contentList []dao.Goods
	dao.Orm().
		Model(&dao.Goods{}).
		Where("OID=? and ID<>? and (GoodsTypeID=? or GoodsTypeChildID=?)", OID, GoodsID, GoodsTypeID, GoodsTypeChildID).
		Order("RAND()").
		Limit(6).
		Find(&contentList)
	return contentList
}
func (service GoodsService) ListGoodsByType(OID, GoodsTypeID, GoodsTypeChildID uint64) []dao.Goods {

	var contentList []dao.Goods

	if GoodsTypeID == 0 {
		dao.Orm().Model(&dao.Goods{}).Where("OID=?", OID).
			Order("CreatedAt desc").Order("ID desc").Find(&contentList)
		return contentList
	}

	if GoodsTypeChildID > 0 {

		dao.Orm().Model(&dao.Goods{}).Where("OID=? and GoodsTypeID=? and GoodsTypeChildID=?", OID, GoodsTypeID, GoodsTypeChildID).
			Order("CreatedAt desc").Order("ID desc").Find(&contentList)
		return contentList
	} else {
		dao.Orm().Model(&dao.Goods{}).Where("OID=? and GoodsTypeID=?", OID, GoodsTypeID).
			Order("CreatedAt desc").Order("ID desc").Find(&contentList)
		return contentList
	}

}
func (service GoodsService) GetGoodsTypeData(OID uint64) *dao.GoodsTypeData {

	goodsTypeData := &dao.GoodsTypeData{}

	rows, err := dao.Orm().Raw(`SELECT gt.*,gtc.* FROM GoodsTypeChild AS gtc LEFT JOIN GoodsType as gt ON (gt.ID=gtc.GoodsTypeID) WHERE OID=?`, OID).Rows()
	if glog.Error(err) {
		return goodsTypeData
	}

	defer rows.Close()

	for rows.Next() {
		var item dao.GoodsTypeGoodsTypeChild
		err := dao.Orm().ScanRows(rows, &item)
		glog.Error(err)

		goodsTypeItem := goodsTypeData.Get(item.GoodsType.ID)
		if goodsTypeItem.Item.ID == 0 {

			goodsTypeData.List = append(goodsTypeData.List, &dao.GoodsTypeItem{Item: &item.GoodsType, SubType: []*dao.GoodsTypeItemSub{{Item: &item.GoodsTypeChild, SubType: []*dao.GoodsTypeItemSub{}}}})

		} else {

			goodsTypeItem.SubType = append(goodsTypeItem.SubType, &dao.GoodsTypeItemSub{Item: &item.GoodsTypeChild, SubType: []*dao.GoodsTypeItemSub{}})

		}

	}

	return goodsTypeData

}
