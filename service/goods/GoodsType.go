package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/pinyin"
	"github.com/nbvghost/gpa/types"
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
	PinyinService pinyin.Service
}

func (service GoodsTypeService) GetTopGoodsTypeChild(DB *gorm.DB, Num uint) []TopGoodsTypeChild {
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

func (service GoodsTypeService) ListGoodsByOID(OID types.PrimaryKey) []model.GoodsType {
	Orm := singleton.Orm()
	var menus []model.GoodsType
	Orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID}).Find(&menus)
	return menus
}

func (service GoodsTypeService) DeleteGoodsType(ID types.PrimaryKey) *result.ActionResult {
	Orm := singleton.Orm()
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
func (service GoodsTypeService) DeleteGoodsTypeChild(GoodsTypeChildID types.PrimaryKey) *result.ActionResult {
	Orm := singleton.Orm()
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

func (service GoodsTypeService) GetGoodsType(ID types.PrimaryKey) model.GoodsType {
	Orm := singleton.Orm()
	var goodsType model.GoodsType
	Orm.Model(&model.GoodsType{}).Where(map[string]interface{}{"ID": ID}).First(&goodsType)
	return goodsType
}
func (service GoodsTypeService) GetGoodsTypeChild(ID types.PrimaryKey) model.GoodsTypeChild {
	Orm := singleton.Orm()
	var typeChild model.GoodsTypeChild
	Orm.Model(&model.GoodsTypeChild{}).Where(map[string]interface{}{"ID": ID}).First(&typeChild)
	return typeChild
}
func (service GoodsTypeService) AllGoodsType() []extends.AllGoodsType {
	gts := service.ListAllGoodsType()

	list := make([]extends.AllGoodsType, 0)
	for i := range gts {
		subs := service.ListAllGoodsTypeChild(gts[i].ID)

		goodsList := service.ListGoodsByGoodsTypeID(gts[i].ID)
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
func (service GoodsTypeService) ListAllGoodsType() []model.GoodsType {
	/*Orm := singleton.Orm()
	var gts []model.GoodsType
	service.FindAll(Orm, &gts)
	return gts*/
	Orm := singleton.Orm()
	var gts []model.GoodsType
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Group("GoodsTypeID").Pluck("GoodsTypeID", &gtsIDs)
	Orm.Model(&model.GoodsType{}).Where(`"ID" in (?)`, gtsIDs).Find(&gts)
	return gts
}
func (service GoodsTypeService) ListGoodsType(OID types.PrimaryKey) []model.GoodsType {
	/*Orm := singleton.Orm()
	var gts []model.GoodsType
	service.FindAllByOID(Orm,&gts,OID)
	return gts*/
	Orm := singleton.Orm()
	var gts []model.GoodsType
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Where(`"OID"=?`, OID).Group("GoodsTypeID").Pluck("GoodsTypeID", &gtsIDs)
	Orm.Model(&model.GoodsType{}).Where(`"ID" in (?)`, gtsIDs).Find(&gts)
	return gts
}
func (service GoodsTypeService) ListGoodsTypeByOIDForAdmin(OID types.PrimaryKey) []model.GoodsType {
	Orm := singleton.Orm()
	var gts []model.GoodsType
	Orm.Model(&model.GoodsType{}).Where(`"OID"=?`, OID).Find(&gts)
	return gts
}
func (service GoodsTypeService) ListGoodsTypeForAdmin() []model.GoodsType {
	Orm := singleton.Orm()
	var gts []model.GoodsType
	Orm.Model(&model.GoodsType{}).Find(&gts)
	return gts
}

func (service GoodsTypeService) ListGoodsTypeChildAll(OID types.PrimaryKey) []model.GoodsTypeChild {
	Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Where(`"OID"=?`, OID).Group("GoodsTypeChildID").Pluck("GoodsTypeChildID", &gtsIDs)
	Orm.Model(&model.GoodsTypeChild{}).Where(`"ID" in (?)`, gtsIDs).Find(&gts)
	return gts
}
func (service GoodsTypeService) ListAllGoodsTypeChild(GoodsTypeID types.PrimaryKey) []model.GoodsTypeChild {
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

func (service GoodsTypeService) ListGoodsTypeChild(GoodsTypeID types.PrimaryKey) []model.GoodsTypeChild {
	/*Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	service.FindWhere(Orm, &gts, model.GoodsTypeChild{GoodsTypeID: GoodsTypeID})
	return gts*/
	Orm := singleton.Orm()
	var gts []model.GoodsTypeChild
	var gtsIDs []uint
	//service.FindWhere(Orm, &gts, model.GoodsTypeChild{})
	Orm.Model(&model.Goods{}).Where(`"GoodsTypeID"=?`, GoodsTypeID).Group("GoodsTypeChildID").Pluck("GoodsTypeChildID", &gtsIDs)
	Orm.Model(&model.GoodsTypeChild{}).Where(`"ID" in (?)`, gtsIDs).Find(&gts)
	return gts
}

func (service GoodsTypeService) ListGoodsChildByGoodsTypeID(GoodsTypeID, GoodsTypeChildID types.PrimaryKey) []types.IEntity {
	//Orm := singleton.Orm()
	//var gts []model.Goods
	//service.FindWhere(Orm, &gts, model.Goods{GoodsTypeID: GoodsTypeID, GoodsTypeChildID: GoodsTypeChildID})
	gts := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"GoodsTypeID"=? and "GoodsTypeChildID"=?`, GoodsTypeID, GoodsTypeChildID).List()
	return gts
}
func (service GoodsTypeService) ListGoodsByGoodsTypeID(GoodsTypeID types.PrimaryKey) []types.IEntity {
	//Orm := singleton.Orm()
	//var gts []model.Goods
	//service.FindWhere(Orm, &gts, model.Goods{GoodsTypeID: GoodsTypeID})
	gts := dao.Find(singleton.Orm(), &model.Goods{}).Where(`"GoodsTypeID"=?`, GoodsTypeID).List()
	return gts
}
func (service GoodsTypeService) AddGoodsTypeByNameByChild(name string, childName string) (goodsType model.GoodsType, goodsTypeChild model.GoodsTypeChild) {
	Orm := singleton.Orm()

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
func (service GoodsTypeService) HotListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID, GoodsTypeChildID, Num uint) []model.Goods {

	Orm := singleton.Orm()

	var result []model.Goods

	db := Orm.Model(&model.Goods{}).Where(`"GoodsTypeID"=? and "GoodsTypeChildID"=?`, GoodsTypeID, GoodsTypeChildID).Order(`"CountSale" desc`).Limit(int(Num))

	db.Find(&result)

	return result

}
func (service GoodsTypeService) NewListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID, GoodsTypeChildID, Num uint) []model.Goods {

	Orm := singleton.Orm()

	var result []model.Goods

	db := Orm.Model(&model.Goods{}).Where(`"GoodsTypeID"=? and "GoodsTypeChildID"=?`, GoodsTypeID, GoodsTypeChildID).Order(`"CreatedAt" desc`).Limit(int(Num))

	db.Find(&result)

	return result

}

func (service GoodsTypeService) ListGoodsByType(OID, GoodsTypeID, GoodsTypeChildID types.PrimaryKey) []model.Goods {

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
func (service GoodsTypeService) GetGoodsTypeData(OID types.PrimaryKey) *extends.GoodsTypeData {

	goodsTypeData := &extends.GoodsTypeData{}

	rows, err := singleton.Orm().Raw(`SELECT gt.*,gtc.* FROM "GoodsTypeChild" AS gtc LEFT JOIN "GoodsType" as gt ON (gt."ID"=gtc."GoodsTypeID") WHERE gtc."OID"=?`, OID).Rows()
	if err != nil {
		return goodsTypeData
	}

	defer rows.Close()

	for rows.Next() {
		var item extends.GoodsTypeGoodsTypeChild
		err := singleton.Orm().ScanRows(rows, &item)
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
func (service GoodsTypeService) getGoodsTypeByName(orm *gorm.DB, OID types.PrimaryKey, name string) (model.GoodsType, error) {
	var gt model.GoodsType
	err := orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Name": name}).First(&gt).Error
	return gt, err
}
func (service GoodsTypeService) getGoodsTypeByUri(orm *gorm.DB, OID types.PrimaryKey, uri string) (model.GoodsType, error) {
	var gt model.GoodsType
	err := orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Uri": uri}).First(&gt).Error
	return gt, err
}
func (service GoodsTypeService) AddGoodsType(OID types.PrimaryKey, name string) error {
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
func (service GoodsTypeService) ChangeGoodsType(OID, ID types.PrimaryKey, name string) error {
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
func (service GoodsTypeService) getGoodsTypeChildByName(orm *gorm.DB, OID, GoodsTypeID types.PrimaryKey, name string) (model.GoodsTypeChild, error) {
	var gt model.GoodsTypeChild
	err := orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "GoodsTypeID": GoodsTypeID, "Name": name}).First(&gt).Error
	return gt, err
}
func (service GoodsTypeService) getGoodsTypeChildByUri(orm *gorm.DB, OID, GoodsTypeID types.PrimaryKey, uri string) (model.GoodsTypeChild, error) {
	var gt model.GoodsTypeChild
	err := orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "GoodsTypeID": GoodsTypeID, "Uri": uri}).First(&gt).Error
	return gt, err
}
func (service GoodsTypeService) AddGoodsTypeChild(OID, GoodsTypeID types.PrimaryKey, name, image string) error {
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
func (service GoodsTypeService) ChangeGoodsTypeChild(OID, ID types.PrimaryKey, name, image string) error {
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
