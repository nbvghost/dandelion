package content

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/gpa/types"
)

func (service ContentService) FindAllContentSubType(OID types.PrimaryKey) []model.ContentSubType {
	Orm := singleton.Orm()
	var list []model.ContentSubType
	Orm.Model(&model.ContentSubType{}).Where(map[string]interface{}{"OID": OID}).Find(&list)
	return list
}
func (service ContentService) FindContentSubTypesByContentItemIDs(ContentItemIDs []uint) []model.ContentSubType {
	Orm := singleton.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where("ContentItemID in (?)", ContentItemIDs).Order("Sort asc").Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByContentItemID(ContentItemID uint) []model.ContentSubType {
	Orm := singleton.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": 0,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByParentContentSubTypeID(ParentContentSubTypeID types.PrimaryKey) []model.ContentSubType {
	Orm := singleton.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where("ParentContentSubTypeID=?", ParentContentSubTypeID).Order("Sort asc").Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByContentItemIDAndParentContentSubTypeID(ContentItemID, ParentContentSubTypeID types.PrimaryKey) []model.ContentSubType {
	Orm := singleton.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": ParentContentSubTypeID,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}

//获取ID，返回子类ID,包括本身
func (service ContentService) GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID types.PrimaryKey) []types.PrimaryKey {
	var IDList []types.PrimaryKey
	singleton.Orm().Model(&model.ContentSubType{}).Where(`"ContentItemID"=? and ("ID"=? or "ParentContentSubTypeID"=?)`, ContentItemID, ContentSubTypeID, ContentSubTypeID).Pluck(`"ID"`, &IDList)
	return IDList
}
func (service ContentService) GetClassifyByName(Name string, ContentItemID, ParentContentSubTypeID types.PrimaryKey) model.ContentSubType {
	Orm := singleton.Orm()
	var menus model.ContentSubType

	Orm.Where(map[string]interface{}{
		"Name":                   Name,
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": ParentContentSubTypeID,
	}).First(&menus)

	return menus

}
func (service ContentService) GetContentSubTypeByNameContentItemIDParentContentSubTypeID(Name string, ContentItemID, ParentContentSubTypeID uint) model.ContentSubType {
	Orm := singleton.Orm()
	var menus model.ContentSubType

	Orm.Where("Name=?", Name).Where("ContentItemID=? and ParentContentSubTypeID=?", ContentItemID, ParentContentSubTypeID).First(&menus)

	return menus

}
