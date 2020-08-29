package service

import (
	"github.com/nbvghost/dandelion/app/service/dao"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
	"strconv"
)

type ContentService struct {
	dao.BaseDao
}

func (self ContentService) GetContentIDs(OID uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if OID <= 0 {
		return levea
	}
	Orm.Model(&dao.Content{}).Where("OID=?", OID).Pluck("ID", &levea)
	return levea
}
func (self ContentService) ListContentSubTypeAction(context *gweb.Context) gweb.Result {
	//Organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)
	PID, _ := strconv.ParseUint(context.Request.URL.Query().Get("PID"), 10, 64)

	content := self.GetContentByID(ContentID)

	csts := self.FindContentSubTypesByContentIDAndParentContentSubTypeID(content.ID, PID)

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: csts}}
}
func (self ContentService) GetClassifyByName(Name string, ContentID, ParentContentSubTypeID uint64) dao.ContentSubType {
	Orm := dao.Orm()
	var menus dao.ContentSubType

	Orm.Where("Name=?", Name).Where("ContentID=? and ParentContentSubTypeID=?", ContentID, ParentContentSubTypeID).First(&menus)

	return menus

}
func (self ContentService) FindContentSubTypesByContentID(ContentID uint64) []dao.ContentSubType {
	Orm := dao.Orm()
	var menus []dao.ContentSubType
	Orm.Model(dao.ContentSubType{}).Where("ContentID=? and ParentContentSubTypeID=0", ContentID).Order("Sort asc").Find(&menus)
	return menus

}
func (self ContentService) FindContentSubTypesByParentContentSubTypeID(ParentContentSubTypeID uint64) []dao.ContentSubType {
	Orm := dao.Orm()
	var menus []dao.ContentSubType
	Orm.Model(dao.ContentSubType{}).Where("ParentContentSubTypeID=?", ParentContentSubTypeID).Order("Sort asc").Find(&menus)
	return menus
}
func (self ContentService) FindContentSubTypesByContentIDAndParentContentSubTypeID(ContentID, ParentContentSubTypeID uint64) []dao.ContentSubType {
	Orm := dao.Orm()
	var menus []dao.ContentSubType
	Orm.Model(dao.ContentSubType{}).Where("ContentID=? and ParentContentSubTypeID=?", ContentID, ParentContentSubTypeID).Order("Sort asc").Find(&menus)
	return menus
}

//-----------------------------------------Content----------------------------------------------------------

func (service ContentService) GetContentByIDAndOID(ID, OID uint64) dao.Content {
	Orm := dao.Orm()
	var menus dao.Content

	Orm.Where("ID=? and OID=?", ID, OID).First(&menus)

	return menus
}
func (service ContentService) GetContentByID(ID uint64) dao.Content {
	Orm := dao.Orm()
	var menus dao.Content
	Orm.Where("ID=?", ID).First(&menus)
	return menus
}
func (service ContentService) GetContentByNameAndOID(Name string, OID uint64) dao.Content {
	Orm := dao.Orm()
	var menus dao.Content

	Orm.Where("Name=? and OID=?", Name, OID).First(&menus)

	return menus
}

func (service ContentService) ListContentType() []dao.ContentType {
	Orm := dao.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var list []dao.ContentType
	err := service.FindAll(Orm, &list)
	glog.Trace(err)
	return list
}
func (service ContentService) ListContentTypeByType(Type string) dao.ContentType {
	Orm := dao.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var list dao.ContentType
	err := service.FindWhere(Orm, &list, "Type=?", Type)
	glog.Trace(err)
	return list
}
func (service ContentService) FindContentSubTypesByNameAndContentID(Name string, ContentID uint64) dao.ContentSubType {
	Orm := dao.Orm()
	var cst dao.ContentSubType
	Orm.Where("ContentID=? and Name=?", ContentID, Name).First(&cst)
	return cst
}
