package content

import (
	"github.com/nbvghost/dandelion/library/db"
	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

func (service ContentService) ListContentItemByOIDMap(OID dao.PrimaryKey) map[dao.PrimaryKey]model.ContentItem {
	Orm := db.Orm()
	var menus []model.ContentItem
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Order(`"Sort"`).Order(`"UpdatedAt" desc`).Find(&menus)
	m := make(map[dao.PrimaryKey]model.ContentItem)
	for i, v := range menus {
		m[v.ID] = menus[i]
	}
	return m
}
func (service ContentService) ListContentItemByOID(OID dao.PrimaryKey) []model.ContentItem {
	Orm := db.Orm()
	var menus []model.ContentItem
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Order(`"Sort"`).Order(`"UpdatedAt" desc`).Find(&menus)
	return menus
}
func (service ContentService) GetContentItemByTypeTemplateName(db *gorm.DB, OID dao.PrimaryKey, typ model.ContentTypeType, templateName string) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(`"OID"=? And "Type"=? And "TemplateName"=?`, OID, typ, templateName).First(&contentItem)
	return &contentItem
}
func (service ContentService) GetContentItemOfIndex(db *gorm.DB, OID dao.PrimaryKey) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(`"OID"=? And "Type"=?`, OID, model.ContentTypeIndex).First(&contentItem)
	return &contentItem
}
func (service ContentService) GetContentItemOfProducts(db *gorm.DB, OID dao.PrimaryKey) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(`"OID"=? And "Type"=?`, OID, model.ContentTypeProducts).First(&contentItem)
	if contentItem.ID == 0 {
		return nil
	}
	return &contentItem
}

func (service ContentService) GetContentItemIDs(OID dao.PrimaryKey) []uint {
	Orm := db.Orm()
	var levea []uint
	if OID <= 0 {
		return levea
	}
	Orm.Model(&model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Pluck(`"ID"`, &levea)
	return levea
}
func (service ContentService) FindContentItemByShowAtHome(OID dao.PrimaryKey) []*model.ContentItem {
	Orm := db.Orm()
	var levea []*model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(map[string]interface{}{"OID": OID, "ShowAtHome": true}).Order(`"Sort"`).Find(&levea)
	return levea
}
func (service ContentService) FindContentItemByTypeTemplate(oid dao.PrimaryKey, contentType string, templateName string, pageIndex int) (int64, []*model.ContentItem) {
	var list []*model.ContentItem
	var total int64

	d := db.Orm().Model(model.ContentItem{}).Order(`"Sort"`).
		Where(`"OID"=? and "Type"=? and "TemplateName"=?`, oid, contentType, templateName)

	d.Count(&total)
	d.Offset(pageIndex * 20).Limit(20).Find(&list)
	return total, list
}
