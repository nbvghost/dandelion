package internal

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type ContentSubTypeDao struct{}

func (ContentSubTypeDao) FindContentSubTypesByNameAndContentItemID(Name string, ContentItemID dao.PrimaryKey) model.ContentSubType {
	Orm := db.Orm()
	var cst model.ContentSubType
	Orm.Where("ContentItemID=? and Name=?", ContentItemID, Name).First(&cst)
	return cst
}
func (ContentSubTypeDao) GetContentSubTypeByID(ID dao.PrimaryKey) model.ContentSubType {
	Orm := db.Orm()
	var menus model.ContentSubType
	Orm.Where(`"ID"=?`, ID).First(&menus)
	return menus
}

// uri 和 name 在 ContentItemID 下面唯一
func (ContentSubTypeDao) GetContentSubTypeByUri(OID, ContentItemID, ID dao.PrimaryKey, uri string) model.ContentSubType {
	Orm := db.Orm()
	var item model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"OID":           OID,
		"ContentItemID": ContentItemID,
		"Uri":           uri,
	}).Where(`"ID"<>?`, ID).First(&item)
	return item
}
func (ContentSubTypeDao) FindAllContentSubType(OID dao.PrimaryKey) []model.ContentSubType {
	Orm := db.Orm()
	var list []model.ContentSubType
	Orm.Model(&model.ContentSubType{}).Where(map[string]interface{}{"OID": OID}).Find(&list)
	return list
}
func (ContentSubTypeDao) FindContentSubTypesByContentItemIDs(ContentItemIDs []uint) []model.ContentSubType {
	Orm := db.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(`"ContentItemID" in (?)`, ContentItemIDs).Order(`"Sort" asc`).Find(&menus)
	return menus
}
func (ContentSubTypeDao) FindContentSubTypesByContentItemID(ContentItemID uint) []model.ContentSubType {
	Orm := db.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": 0,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}
func (ContentSubTypeDao) FindContentSubTypesByParentContentSubTypeID(ParentContentSubTypeID dao.PrimaryKey) []model.ContentSubType {
	Orm := db.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(`"ParentContentSubTypeID"=?`, ParentContentSubTypeID).Order(`"Sort" asc`).Find(&menus)
	return menus
}
func (ContentSubTypeDao) FindContentSubTypesByContentItemIDAndParentContentSubTypeID(ContentItemID, ParentContentSubTypeID dao.PrimaryKey) []model.ContentSubType {
	Orm := db.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": ParentContentSubTypeID,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}

// 获取ID，返回子类ID,包括本身
func (ContentSubTypeDao) GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID dao.PrimaryKey) []dao.PrimaryKey {
	var IDList []dao.PrimaryKey
	db.Orm().Model(&model.ContentSubType{}).Where(`"ContentItemID"=? and ("ID"=? or "ParentContentSubTypeID"=?)`, ContentItemID, ContentSubTypeID, ContentSubTypeID).Pluck(`"ID"`, &IDList)
	return IDList
}
func (ContentSubTypeDao) GetContentSubTypeByName(OID, ContentItemID, ID dao.PrimaryKey, Name string) model.ContentSubType {
	Orm := db.Orm()
	var menus model.ContentSubType
	Orm.Where(map[string]interface{}{
		"OID":           OID,
		"ContentItemID": ContentItemID,
		"Name":          Name,
	}).Where(`"ID"<>?`, ID).First(&menus)
	return menus

}
func (ContentSubTypeDao) GetContentSubTypeByNameContentItemIDParentContentSubTypeID(Name string, ContentItemID, ParentContentSubTypeID uint) model.ContentSubType {
	Orm := db.Orm()
	var menus model.ContentSubType

	Orm.Where(map[string]any{"Name": Name, "ContentItemID": ContentItemID, "ParentContentSubTypeID": ParentContentSubTypeID}).First(&menus)

	return menus
}
