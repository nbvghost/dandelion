package internal

import (
	"context"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type ContentSubTypeDao struct{}

func (ContentSubTypeDao) FindContentSubTypesByNameAndContentItemID(ctx context.Context, Name string, ContentItemID dao.PrimaryKey) model.ContentSubType {
	Orm := db.GetDB(ctx)
	var cst model.ContentSubType
	Orm.Where("ContentItemID=? and Name=?", ContentItemID, Name).First(&cst)
	return cst
}
func (ContentSubTypeDao) GetContentSubTypeByID(ctx context.Context, ID dao.PrimaryKey) model.ContentSubType {
	Orm := db.GetDB(ctx)
	var menus model.ContentSubType
	Orm.Where(`"ID"=?`, ID).First(&menus)
	return menus
}

// uri 和 name 在 ContentItemID 下面唯一
func (ContentSubTypeDao) GetContentSubTypeByUri(ctx context.Context, OID, ContentItemID, ID dao.PrimaryKey, uri string) model.ContentSubType {
	Orm := db.GetDB(ctx)
	var item model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"OID":           OID,
		"ContentItemID": ContentItemID,
		"Uri":           uri,
	}).Where(`"ID"<>?`, ID).First(&item)
	return item
}
func (ContentSubTypeDao) FindAllContentSubType(ctx context.Context, OID dao.PrimaryKey) []model.ContentSubType {
	Orm := db.GetDB(ctx)
	var list []model.ContentSubType
	Orm.Model(&model.ContentSubType{}).Where(map[string]interface{}{"OID": OID}).Find(&list)
	return list
}
func (ContentSubTypeDao) FindContentSubTypesByContentItemIDs(ctx context.Context, ContentItemIDs []uint) []model.ContentSubType {
	Orm := db.GetDB(ctx)
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(`"ContentItemID" in (?)`, ContentItemIDs).Order(`"Sort" asc`).Find(&menus)
	return menus
}
func (ContentSubTypeDao) FindContentSubTypesByContentItemID(ctx context.Context, ContentItemID uint) []model.ContentSubType {
	Orm := db.GetDB(ctx)
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": 0,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}
func (ContentSubTypeDao) FindContentSubTypesByParentContentSubTypeID(ctx context.Context, ParentContentSubTypeID dao.PrimaryKey) []model.ContentSubType {
	Orm := db.GetDB(ctx)
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(`"ParentContentSubTypeID"=?`, ParentContentSubTypeID).Order(`"Sort" asc`).Find(&menus)
	return menus
}
func (ContentSubTypeDao) FindContentSubTypesByContentItemIDAndParentContentSubTypeID(ctx context.Context, ContentItemID, ParentContentSubTypeID dao.PrimaryKey) []model.ContentSubType {
	Orm := db.GetDB(ctx)
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": ParentContentSubTypeID,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}

// 获取ID，返回子类ID,包括本身
func (ContentSubTypeDao) GetContentSubTypeAllIDByID(ctx context.Context, ContentItemID, ContentSubTypeID dao.PrimaryKey) []dao.PrimaryKey {
	var IDList []dao.PrimaryKey
	db.GetDB(ctx).Model(&model.ContentSubType{}).Where(`"ContentItemID"=? and ("ID"=? or "ParentContentSubTypeID"=?)`, ContentItemID, ContentSubTypeID, ContentSubTypeID).Pluck(`"ID"`, &IDList)
	return IDList
}
func (ContentSubTypeDao) GetContentSubTypeByName(ctx context.Context, OID, ContentItemID, ID dao.PrimaryKey, Name string) model.ContentSubType {
	Orm := db.GetDB(ctx)
	var menus model.ContentSubType
	Orm.Where(map[string]interface{}{
		"OID":           OID,
		"ContentItemID": ContentItemID,
		"Name":          Name,
	}).Where(`"ID"<>?`, ID).First(&menus)
	return menus

}
func (ContentSubTypeDao) GetContentSubTypeByNameContentItemIDParentContentSubTypeID(ctx context.Context, Name string, ContentItemID, ParentContentSubTypeID uint) model.ContentSubType {
	Orm := db.GetDB(ctx)
	var menus model.ContentSubType

	Orm.Where(map[string]any{"Name": Name, "ContentItemID": ContentItemID, "ParentContentSubTypeID": ParentContentSubTypeID}).First(&menus)

	return menus
}
