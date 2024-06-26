package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/domain/cache"

	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type TopGoodsTypeChild struct {
	Name             string `gorm:"column:Name"`
	Image            string `gorm:"column:Image"`
	GoodsTypeChildID uint   `gorm:"column:GoodsTypeChildID"`
	Price            uint   `gorm:"column:Price"`
}

type GoodsTypeService struct {

}

func (m GoodsTypeService) GetTopGoodsTypeChild(DB *gorm.DB, Num uint) []TopGoodsTypeChild {
	list := make([]TopGoodsTypeChild, 0)
	//SELECT gtc.Name,gtc.Image,gtc.ID AS GoodsTypeChildID,gtc.GoodsTypeID AS GoodsTypeID,MIN(g.Price) FROM Goods AS g LEFT JOIN GoodsTypeChild AS gtc ON (gtc.GoodsTypeID=g.GoodsTypeID AND gtc.ID=g.GoodsTypeChildID) GROUP BY g.GoodsTypeID;
	rows, err := DB.Raw(`
SELECT 
gtc."Name" as Name,
gtc."Image" as Image,
gtc."ID" AS GoodsTypeChildID,
gtc."GoodsTypeID" AS GoodsTypeID,
MIN(g."Price") as Price
FROM "Goods" AS g LEFT JOIN "GoodsTypeChild" AS gtc ON (gtc."GoodsTypeID"=g."GoodsTypeID" AND gtc."ID"=g."GoodsTypeChildID") GROUP BY g."GoodsTypeID" limit ?
`, Num).Rows()
	if err != nil {
		return list
	}
	for rows.Next() {
		var result TopGoodsTypeChild
		DB.ScanRows(rows, &result)
		list = append(list, result)
	}
	return list
}

func (m GoodsTypeService) ListGoodsByOID(OID dao.PrimaryKey) []model.GoodsType {
	Orm := db.Orm()
	var menus []model.GoodsType
	Orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID}).Find(&menus)
	return menus
}

func (m GoodsTypeService) DeleteGoodsType(ID dao.PrimaryKey) *result.ActionResult {
	Orm := db.Orm()
	tx := Orm.Begin()
	var gtcs []model.GoodsTypeChild
	tx.Where(&model.GoodsTypeChild{GoodsTypeID: ID}).Find(&gtcs) //Updates(map[string]interface{}{"GoodsTypeID": 0})

	var err error
	if len(gtcs) <= 0 {
		err = dao.DeleteByPrimaryKey(tx, &model.GoodsType{}, ID)
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
func (m GoodsTypeService) DeleteGoodsTypeChild(GoodsTypeChildID dao.PrimaryKey) *result.ActionResult {
	Orm := db.Orm()
	tx := Orm.Begin()
	tx.Model(&model.Goods{GoodsTypeChildID: GoodsTypeChildID}).Updates(map[string]interface{}{"GoodsTypeChildID": 0})
	err := dao.DeleteByPrimaryKey(tx, &model.GoodsTypeChild{}, GoodsTypeChildID)
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

func (m GoodsTypeService) GetGoodsType(ID dao.PrimaryKey) model.GoodsType {
	Orm := db.Orm()
	var goodsType model.GoodsType
	Orm.Model(&model.GoodsType{}).Where(map[string]interface{}{"ID": ID}).First(&goodsType)
	return goodsType
}
func (m GoodsTypeService) GetGoodsTypeChild(ID dao.PrimaryKey) model.GoodsTypeChild {
	Orm := db.Orm()
	var typeChild model.GoodsTypeChild
	Orm.Model(&model.GoodsTypeChild{}).Where(map[string]interface{}{"ID": ID}).First(&typeChild)
	return typeChild
}
func (m GoodsTypeService) AllGoodsType() []extends.AllGoodsType {
	gts := m.ListAllGoodsType()

	list := make([]extends.AllGoodsType, 0)
	for i := range gts {
		subs := m.ListAllGoodsTypeChild(gts[i].ID)

		goodsList := m.ListGoodsByGoodsTypeID(gts[i].ID)
		var maxPrice uint = 0
		var minPrice uint = 0
		for ii := range goodsList {
			item := goodsList[ii].(*model.Goods)
			if ii == 0 {
				minPrice = item.Price
			}
			if item.Price > maxPrice {
				maxPrice = item.Price
			}
			if item.Price < minPrice {
				minPrice = item.Price
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
func (m GoodsTypeService) ListAllGoodsType() []model.GoodsType {
	/*Orm := singleton.Orm()
	var gts []model.GoodsType
	service.FindAll(Orm, &gts)
	return gts*/
	Orm := db.Orm()
	var gts []model.GoodsType
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Group("GoodsTypeID").Pluck("GoodsTypeID", &gtsIDs)
	Orm.Model(&model.GoodsType{}).Where(`"ID" in (?)`, gtsIDs).Find(&gts)
	return gts
}
func (m GoodsTypeService) ListGoodsType(OID dao.PrimaryKey) []model.GoodsType {
	/*Orm := singleton.Orm()
	var gts []model.GoodsType
	service.FindAllByOID(Orm,&gts,OID)
	return gts*/
	Orm := db.Orm()

	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Where(`"OID"=?`, OID).Group("GoodsTypeID").Pluck("GoodsTypeID", &gtsIDs)

	var gts []model.GoodsType
	Orm.Model(&model.GoodsType{}).Where(`"ID" in (?)`, gtsIDs).Find(&gts)
	return gts
}
func (m GoodsTypeService) ListGoodsTypeByOIDForAdmin(OID dao.PrimaryKey) []model.GoodsType {
	Orm := db.Orm()
	var gts []model.GoodsType
	Orm.Model(&model.GoodsType{}).Where(`"OID"=?`, OID).Find(&gts)
	return gts
}
func (m GoodsTypeService) ListGoodsTypeForAdmin() []model.GoodsType {
	Orm := db.Orm()
	var gts []model.GoodsType
	Orm.Model(&model.GoodsType{}).Find(&gts)
	return gts
}

func (m GoodsTypeService) ListGoodsTypeChildAll(OID dao.PrimaryKey) []model.GoodsTypeChild {
	Orm := db.Orm()
	var gts []model.GoodsTypeChild
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Where(`"OID"=?`, OID).Group("GoodsTypeChildID").Pluck("GoodsTypeChildID", &gtsIDs)
	Orm.Model(&model.GoodsTypeChild{}).Where(`"ID" in (?)`, gtsIDs).Find(&gts)
	return gts
}
func (m GoodsTypeService) ListAllGoodsTypeChild(GoodsTypeID dao.PrimaryKey) []model.GoodsTypeChild {
	/*Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	service.FindWhere(Orm, &gts, model.GoodsTypeChild{GoodsTypeID: GoodsTypeID})
	return gts*/
	Orm := db.Orm()
	var gts []model.GoodsTypeChild
	Orm.Model(&model.GoodsTypeChild{}).Where(map[string]interface{}{
		"GoodsTypeID": GoodsTypeID,
	}).Find(&gts)
	return gts
}

func (m GoodsTypeService) ListGoodsTypeChild(GoodsTypeID dao.PrimaryKey) []model.GoodsTypeChild {
	/*Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	service.FindWhere(Orm, &gts, model.GoodsTypeChild{GoodsTypeID: GoodsTypeID})
	return gts*/
	Orm := db.Orm()
	var gts []model.GoodsTypeChild
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Where(`"GoodsTypeID"=?`, GoodsTypeID).Group("GoodsTypeChildID").Pluck("GoodsTypeChildID", &gtsIDs)
	Orm.Model(&model.GoodsTypeChild{}).Where(`"ID" in (?)`, gtsIDs).Find(&gts)
	return gts
}

func (m GoodsTypeService) ListGoodsChildByGoodsTypeID(GoodsTypeID, GoodsTypeChildID dao.PrimaryKey) []dao.IEntity {
	//Orm := singleton.Orm()
	//var gts []model.Goods
	//service.FindWhere(Orm, &gts, model.Goods{GoodsTypeID: GoodsTypeID, GoodsTypeChildID: GoodsTypeChildID})
	gts := dao.Find(db.Orm(), &model.Goods{}).Where(`"GoodsTypeID"=? and "GoodsTypeChildID"=?`, GoodsTypeID, GoodsTypeChildID).List()
	return gts
}
func (m GoodsTypeService) ListGoodsByGoodsTypeID(GoodsTypeID dao.PrimaryKey) []dao.IEntity {
	//Orm := singleton.Orm()
	//var gts []model.Goods
	//service.FindWhere(Orm, &gts, model.Goods{GoodsTypeID: GoodsTypeID})
	gts := dao.Find(db.Orm(), &model.Goods{}).Where(`"GoodsTypeID"=?`, GoodsTypeID).List()
	return gts
}
func (m GoodsTypeService) AddGoodsTypeByNameByChild(name string, childName string) (goodsType model.GoodsType, goodsTypeChild model.GoodsTypeChild) {
	Orm := db.Orm()

	var gt model.GoodsType
	var gtc model.GoodsTypeChild

	err := Orm.Model(&model.GoodsType{}).Where(`"Name"=?`, name).First(&gt).Error
	if gorm.ErrRecordNotFound == err {
		gt.Name = name
		dao.Save(Orm, &gt)
	}

	err = Orm.Model(&model.GoodsTypeChild{}).Where(`"Name"=? and "GoodsTypeID"=?`, childName, gt.ID).First(&gtc).Error
	if gorm.ErrRecordNotFound == err {
		gtc.Name = childName
		gtc.GoodsTypeID = gt.ID
		dao.Save(Orm, &gtc)
	}

	return gt, gtc
}
func (m GoodsTypeService) HotListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID, GoodsTypeChildID, Num uint) []model.Goods {

	Orm := db.Orm()

	var result []model.Goods

	db := Orm.Model(&model.Goods{}).Where(`"GoodsTypeID"=? and "GoodsTypeChildID"=?`, GoodsTypeID, GoodsTypeChildID).Order(`"CountSale" desc`).Limit(int(Num))

	db.Find(&result)

	return result

}
func (m GoodsTypeService) NewListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID, GoodsTypeChildID, Num uint) []model.Goods {

	Orm := db.Orm()

	var result []model.Goods

	db := Orm.Model(&model.Goods{}).Where(`"GoodsTypeID"=? and "GoodsTypeChildID"=?`, GoodsTypeID, GoodsTypeChildID).Order(`"CreatedAt" desc`).Limit(int(Num))

	db.Find(&result)

	return result

}

func (m GoodsTypeService) ListGoodsByType(OID, GoodsTypeID, GoodsTypeChildID dao.PrimaryKey) []model.Goods {

	var contentList []model.Goods

	if GoodsTypeID == 0 {
		db.Orm().Model(&model.Goods{}).Where(`"OID"=?`, OID).
			Order(`"CreatedAt" desc`).Order(`"ID" desc`).Find(&contentList)
		return contentList
	}

	if GoodsTypeChildID > 0 {
		db.Orm().Model(&model.Goods{}).Where(`"OID"=? and "GoodsTypeID"=? and "GoodsTypeChildID"=?`, OID, GoodsTypeID, GoodsTypeChildID).
			Order(`"CreatedAt" desc`).Order(`"ID" desc`).Find(&contentList)
		return contentList
	} else {
		db.Orm().Model(&model.Goods{}).Where(`"OID"=? and "GoodsTypeID"=?`, OID, GoodsTypeID).
			Order(`"CreatedAt" desc`).Order(`"ID" desc`).Find(&contentList)
		return contentList
	}

}
func (m GoodsTypeService) GetGoodsTypeData(OID dao.PrimaryKey) *extends.GoodsTypeData {

	goodsTypeData := &extends.GoodsTypeData{}

	rows, err := db.Orm().Raw(`SELECT gt.*,gtc.* FROM "GoodsTypeChild" AS gtc LEFT JOIN "GoodsType" as gt ON (gt."ID"=gtc."GoodsTypeID") WHERE gtc."OID"=?`, OID).Rows()
	if err != nil {
		return goodsTypeData
	}

	defer rows.Close()

	for rows.Next() {
		var item extends.GoodsTypeGoodsTypeChild
		err := db.Orm().ScanRows(rows, &item)
		if err != nil {
			log.Println(err)
		}

		goodsTypeItem := goodsTypeData.Get(item.GoodsType.ID)
		if goodsTypeItem.Item.ID == 0 {

			goodsTypeData.List = append(goodsTypeData.List, &extends.GoodsTypeItem{Item: &item.GoodsType, SubType: []*extends.GoodsTypeItemSub{{Item: &item.GoodsTypeChild, SubType: []*extends.GoodsTypeItemSub{}}}})

		} else {

			goodsTypeItem.SubType = append(goodsTypeItem.SubType, &extends.GoodsTypeItemSub{Item: &item.GoodsTypeChild, SubType: []*extends.GoodsTypeItemSub{}})

		}

	}

	return goodsTypeData

}
func (m GoodsTypeService) getGoodsTypeByName(orm *gorm.DB, OID dao.PrimaryKey, name string) (model.GoodsType, error) {
	var gt model.GoodsType
	err := orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Name": name}).First(&gt).Error
	return gt, err
}
func (m GoodsTypeService) getGoodsTypeByUri(orm *gorm.DB, OID dao.PrimaryKey, uri string) (model.GoodsType, error) {
	var gt model.GoodsType
	err := orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Uri": uri}).First(&gt).Error
	return gt, err
}
func (m GoodsTypeService) AddGoodsType(OID dao.PrimaryKey, goodsType *model.GoodsType) error {
	orm := db.Orm()
	gt, _ := m.getGoodsTypeByName(orm, OID, goodsType.Name)
	if !gt.IsZero() {
		return errors.Errorf("重复的名字:%s", goodsType.Name)
	}

	uri := cache.Cache.ChinesePinyinCache.AutoDetectUri(goodsType.Name)
	gt, _ = m.getGoodsTypeByUri(orm, OID, uri)
	if !gt.IsZero() {
		gt.Uri = fmt.Sprintf("%s-%d", gt.Uri, time.Now().Unix())
	}
	gt.OID = OID
	gt.Name = goodsType.Name
	gt.Uri = uri
	gt.Introduction = goodsType.Introduction
	return orm.Model(model.GoodsType{}).Create(&gt).Error
}
func (m GoodsTypeService) ChangeGoodsType(OID dao.PrimaryKey, goodsType *model.GoodsType) error {
	orm := db.Orm()
	gt, _ := m.getGoodsTypeByName(orm, OID, goodsType.Name)
	if !gt.IsZero() && gt.ID != goodsType.ID {
		return errors.Errorf("重复的名字:%s", goodsType.Name)
	}

	uri := cache.Cache.ChinesePinyinCache.AutoDetectUri(goodsType.Name)
	gt, _ = m.getGoodsTypeByUri(orm, OID, uri)
	if !gt.IsZero() {
		gt.Uri = fmt.Sprintf("%s-%d", gt.Uri, time.Now().Unix())
	}
	gt.Name = goodsType.Name
	gt.Uri = uri
	return orm.Model(model.GoodsType{}).Where(`"ID"=?`, goodsType.ID).Updates(map[string]interface{}{
		"Name":         gt.Name,
		"Uri":          gt.Uri,
		"Introduction": goodsType.Introduction,
		"Badge":        goodsType.Badge,
		"Image":        goodsType.Image,
		"IsStickyTop":  goodsType.IsStickyTop,
	}).Error
}
func (m GoodsTypeService) getGoodsTypeChildByName(orm *gorm.DB, OID, GoodsTypeID dao.PrimaryKey, name string) (model.GoodsTypeChild, error) {
	var gt model.GoodsTypeChild
	err := orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "GoodsTypeID": GoodsTypeID, "Name": name}).First(&gt).Error
	return gt, err
}
func (m GoodsTypeService) getGoodsTypeChildByUri(orm *gorm.DB, OID, GoodsTypeID dao.PrimaryKey, uri string) (model.GoodsTypeChild, error) {
	var gt model.GoodsTypeChild
	err := orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "GoodsTypeID": GoodsTypeID, "Uri": uri}).First(&gt).Error
	return gt, err
}
func (m GoodsTypeService) AddGoodsTypeChild(OID, GoodsTypeID dao.PrimaryKey, name, image string) error {
	if GoodsTypeID == 0 {
		return errors.Errorf("没有指定父类ID")
	}
	orm := db.Orm()
	gtc := m.GetGoodsType(GoodsTypeID)
	if gtc.IsZero() {
		return errors.Errorf("不存在父类:%d", GoodsTypeID)
	}

	gt, _ := m.getGoodsTypeChildByName(orm, OID, GoodsTypeID, name)
	if !gt.IsZero() {
		return errors.Errorf("重复的名字:%s", name)
	}

	uri := cache.Cache.ChinesePinyinCache.AutoDetectUri(name)
	gt, _ = m.getGoodsTypeChildByUri(orm, OID, GoodsTypeID, uri)
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
func (m GoodsTypeService) ChangeGoodsTypeChild(OID, ID dao.PrimaryKey, name, image string) error {
	orm := db.Orm()
	gtc := m.GetGoodsTypeChild(ID)
	if gtc.IsZero() {
		return errors.Errorf("记录不存在")
	}

	gt, _ := m.getGoodsTypeChildByName(orm, OID, gtc.GoodsTypeID, name)
	if gt.ID != ID && !gt.IsZero() {
		//说明，修改的是同一条记录，如果不相等的话，说明有存在相同name的记录。
		return errors.Errorf("重复的名字:%s", name)
	}

	uri := cache.Cache.ChinesePinyinCache.AutoDetectUri(name)
	gt, _ = m.getGoodsTypeChildByUri(orm, OID, gtc.GoodsTypeID, uri)
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

func (m GoodsTypeService) StickyTopGoodsTypeList(orm *gorm.DB, OID dao.PrimaryKey) []model.GoodsType {
	var gt []model.GoodsType
	orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "IsStickyTop": true}).Order(`"UpdatedAt" desc`).Find(&gt)
	return gt
}
