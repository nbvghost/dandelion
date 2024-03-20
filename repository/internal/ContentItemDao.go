package internal

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"gorm.io/gorm"
)

type ContentItemDao struct{}

func (ContentItemDao) GetContentItemByUri(OID, ID dao.PrimaryKey, uri string) model.ContentItem {
	Orm := db.Orm()
	var item model.ContentItem
	item.OID = OID
	item.Uri = uri
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{"OID": item.OID, "Uri": item.Uri}).Where(`"ID"<>?`, ID).First(&item)
	return item
}
func (ContentItemDao) FindContentItemByType(Type model.ContentTypeType, OID dao.PrimaryKey) []model.ContentItem {
	Orm := db.Orm()
	menus := make([]model.ContentItem, 0)
	Orm.Where(map[string]interface{}{
		"Type": Type,
		"OID":  OID,
	}).Find(&menus)
	return menus
}
func (ContentItemDao) ExistContentItemByNameAndOID(OID, ID dao.PrimaryKey, Name string) model.ContentItem {
	Orm := db.Orm()
	var menus model.ContentItem
	Orm.Where(`"OID"=?`, OID).Where(map[string]interface{}{"Name": Name}).Where(`"ID"<>?`, ID).First(&menus)
	return menus
}
func (ContentItemDao) GetContentItemByIDAndOID(ID, OID uint) model.ContentItem {
	Orm := db.Orm()
	var menus model.ContentItem

	Orm.Where(`"ID"=? and "OID"=?`, ID, OID).First(&menus)

	return menus
}
func (ContentItemDao) GetContentItemByID(ID dao.PrimaryKey) model.ContentItem {
	Orm := db.Orm()
	var menus model.ContentItem
	Orm.Where(`"ID"=?`, ID).First(&menus)
	return menus
}
func (ContentItemDao) ListContentItemByOID(OID dao.PrimaryKey) []model.ContentItem {
	Orm := db.Orm()
	var menus []model.ContentItem
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Order(`"Sort"`).Order(`"UpdatedAt" desc`).Find(&menus)
	return menus
}
func (ContentItemDao) GetContentItemByTypeTemplateName(db *gorm.DB, OID dao.PrimaryKey, typ model.ContentTypeType, templateName string) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(`"OID"=? And "Type"=? And "TemplateName"=?`, OID, typ, templateName).First(&contentItem)
	return &contentItem
}
func (ContentItemDao) GetContentItemOfIndex(db *gorm.DB, OID dao.PrimaryKey) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(`"OID"=? And "Type"=?`, OID, model.ContentTypeIndex).First(&contentItem)
	return &contentItem
}
func (ContentItemDao) GetContentItemOfProducts(db *gorm.DB, OID dao.PrimaryKey) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(`"OID"=? And "Type"=?`, OID, model.ContentTypeProducts).First(&contentItem)
	if contentItem.ID == 0 {
		return nil
	}
	return &contentItem
}
func (ContentItemDao) ListContentItemByOIDMap(OID dao.PrimaryKey) map[dao.PrimaryKey]model.ContentItem {
	Orm := db.Orm()
	var menus []model.ContentItem
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Order(`"Sort"`).Order(`"UpdatedAt" desc`).Find(&menus)
	m := make(map[dao.PrimaryKey]model.ContentItem)
	for i, v := range menus {
		m[v.ID] = menus[i]
	}
	return m
}

func (ContentItemDao) GetContentItemIDs(OID dao.PrimaryKey) []uint {
	Orm := db.Orm()
	var levea []uint
	if OID <= 0 {
		return levea
	}
	Orm.Model(&model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Pluck(`"ID"`, &levea)
	return levea
}
func (ContentItemDao) FindContentItemByShowAtHome(OID dao.PrimaryKey) []*model.ContentItem {
	Orm := db.Orm()
	var levea []*model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(map[string]interface{}{"OID": OID, "ShowAtHome": true}).Order(`"Sort"`).Find(&levea)
	return levea
}
func (ContentItemDao) FindContentItemByTypeTemplate(oid dao.PrimaryKey, contentType string, templateName string, pageIndex int) (int64, []*model.ContentItem) {
	var list []*model.ContentItem
	var total int64

	d := db.Orm().Model(model.ContentItem{}).Order(`"Sort"`).
		Where(`"OID"=? and "Type"=? and "TemplateName"=?`, oid, contentType, templateName)

	d.Count(&total)
	d.Offset(pageIndex * 20).Limit(20).Find(&list)
	return total, list
}
