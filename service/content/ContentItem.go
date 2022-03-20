package content

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/gpa/types"
	"gorm.io/gorm"
)

func (service ContentService) ListContentItemByOID(OID types.PrimaryKey) []model.ContentItem {
	Orm := singleton.Orm()
	var menus []model.ContentItem
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Order(`"Sort"`).Order(`"UpdatedAt" desc`).Find(&menus)
	return menus
}
func (service ContentService) GetContentItemDefault(db *gorm.DB, OID types.PrimaryKey) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where("OID=? And Type=?", OID, model.ContentTypeProducts).First(&contentItem)
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
