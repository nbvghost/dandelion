package content

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

func (service ContentService) FindAllContentSubType(OID dao.PrimaryKey) []model.ContentSubType {
	Orm := db.Orm()
	var list []model.ContentSubType
	Orm.Model(&model.ContentSubType{}).Where(map[string]interface{}{"OID": OID}).Find(&list)
	return list
}
func (service ContentService) FindContentSubTypesByContentItemIDs(ContentItemIDs []uint) []model.ContentSubType {
	Orm := db.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where("ContentItemID in (?)", ContentItemIDs).Order("Sort asc").Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByContentItemID(ContentItemID uint) []model.ContentSubType {
	Orm := db.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": 0,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByParentContentSubTypeID(ParentContentSubTypeID dao.PrimaryKey) []model.ContentSubType {
	Orm := db.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where("ParentContentSubTypeID=?", ParentContentSubTypeID).Order("Sort asc").Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByContentItemIDAndParentContentSubTypeID(ContentItemID, ParentContentSubTypeID dao.PrimaryKey) []model.ContentSubType {
	Orm := db.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": ParentContentSubTypeID,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}

// 获取ID，返回子类ID,包括本身
func (service ContentService) GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID dao.PrimaryKey) []dao.PrimaryKey {
	var IDList []dao.PrimaryKey
	db.Orm().Model(&model.ContentSubType{}).Where(`"ContentItemID"=? and ("ID"=? or "ParentContentSubTypeID"=?)`, ContentItemID, ContentSubTypeID, ContentSubTypeID).Pluck(`"ID"`, &IDList)
	return IDList
}
func (service ContentService) GetContentSubTypeByName(OID, ContentItemID, ID dao.PrimaryKey, Name string) model.ContentSubType {
	Orm := db.Orm()
	var menus model.ContentSubType
	Orm.Where(map[string]interface{}{
		"OID":           OID,
		"ContentItemID": ContentItemID,
		"Name":          Name,
	}).Where(`"ID"<>?`, ID).First(&menus)
	return menus

}
func (service ContentService) GetContentSubTypeByNameContentItemIDParentContentSubTypeID(Name string, ContentItemID, ParentContentSubTypeID uint) model.ContentSubType {
	Orm := db.Orm()
	var menus model.ContentSubType

	Orm.Where("Name=?", Name).Where("ContentItemID=? and ParentContentSubTypeID=?", ContentItemID, ParentContentSubTypeID).First(&menus)

	return menus

}
