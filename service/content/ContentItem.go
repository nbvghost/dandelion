package content

import (
	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/gpa/types"
)

func (service ContentService) ListContentItemByOID(OID types.PrimaryKey) []model.ContentItem {
	Orm := singleton.Orm()
	var menus []model.ContentItem
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Order(`"Sort"`).Order(`"UpdatedAt" desc`).Find(&menus)
	return menus
}
func (service ContentService) GetContentItemByTypeTemplateName(db *gorm.DB, OID types.PrimaryKey, typ model.ContentTypeType, templateName string) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(`"OID"=? And "Type"=? And "TemplateName"=?`, OID, typ, templateName).First(&contentItem)
	return &contentItem
}
func (service ContentService) GetContentItemOfIndex(db *gorm.DB, OID types.PrimaryKey) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(`"OID"=? And "Type"=?`, OID, model.ContentTypeIndex).First(&contentItem)
	return &contentItem
}
func (service ContentService) GetContentItemOfProducts(db *gorm.DB, OID types.PrimaryKey) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where(`"OID"=? And "Type"=?`, OID, model.ContentTypeProducts).First(&contentItem)
	if contentItem.ID == 0 {
		return nil
	}
	return &contentItem
}

func (service ContentService) GetContentItemIDs(OID types.PrimaryKey) []uint {
	Orm := singleton.Orm()
	var levea []uint
	if OID <= 0 {
		return levea
	}
	Orm.Model(&model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Pluck(`"ID"`, &levea)
	return levea
}
