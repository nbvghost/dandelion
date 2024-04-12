package internal

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"gorm.io/gorm"
)

type ContentConfigDao struct{}

func (m ContentConfigDao) AddContentConfig(db *gorm.DB, company *model.Organization) error {
	Orm := db
	item := m.GetContentConfig(db, company.ID)
	if (&item).IsZero() {
		err := Orm.Create(&model.ContentConfig{OID: company.ID, Name: company.Name}).Error
		return err
	}
	return nil
}

func (m ContentConfigDao) GetContentConfig(orm *gorm.DB, OID dao.PrimaryKey) model.ContentConfig {
	var contentConfig model.ContentConfig
	orm.Model(&model.ContentConfig{}).Where(map[string]interface{}{"OID": OID}).First(&contentConfig)
	return contentConfig
}
