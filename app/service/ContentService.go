package service

import (
	"dandelion/app/service/dao"

	"dandelion/app/play"
	"dandelion/app/util"
	"errors"
	"strconv"
	"strings"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
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
func (service ContentService) ChangeClassify(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.ContentSubType{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	m := service.GetClassifyByName(item.Name, item.ContentID, item.ParentContentSubTypeID)
	if m.ID != 0 && m.ID != item.ID {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("名字重复，修改失败"), "", nil)}
	}
	err = service.ChangeModel(Orm, ID, &dao.ContentSubType{Name: item.Name})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (service ContentService) DeleteClassify(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	css := service.FindContentSubTypesByParentContentSubTypeID(ID)
	if len(css) > 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("包含子项内容，无法删除"), "删除成功", nil)}
	}
	articles := ArticleService{}.FindArticleByContentSubTypeID(ID)
	if len(articles) > 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("包含文章，无法删除"), "删除成功", nil)}
	}

	item := &dao.ContentSubType{}
	err := service.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (service ContentService) ListClassify(context *gweb.Context) gweb.Result {
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	list := service.FindContentSubTypesByContentID(ContentID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}

}
func (service ContentService) ListChildClassify(context *gweb.Context) gweb.Result {
	ParentContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ParentContentSubTypeID"], 10, 64)
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	list := service.FindContentSubTypesByContentIDAndParentContentSubTypeID(ContentID, ParentContentSubTypeID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}

}
func (service ContentService) AddClassify(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	item := &dao.ContentSubType{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	have := service.GetClassifyByName(item.Name, item.ContentID, item.ParentContentSubTypeID)
	if have.ID != 0 && have.ID != item.ID {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("这个名字已经被使用了"), "", nil)}
	}

	//item.OID = company.ID
	err = service.Add(Orm, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (service ContentService) GetContentSubTypeAction(context *gweb.Context) gweb.Result {
	ContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)

	Orm := dao.Orm()
	var menus dao.ContentSubType
	var pmenus dao.ContentSubType

	Orm.Where("ID=?", ContentSubTypeID).First(&menus)

	if menus.ID > 0 {
		Orm.Where("ID=?", menus.ParentContentSubTypeID).First(&pmenus)
	}
	results := make(map[string]interface{})
	results["ContentSubType"] = menus
	results["ParentContentSubType"] = pmenus

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", results)}
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
func (service ContentService) ListContentTypeAction(context *gweb.Context) gweb.Result {
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: service.ListContentType()}}
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
	err := Orm.Where("ContentID=? and Name=?", ContentID, Name).First(&cst).Error
	glog.Error(err)
	return cst
}

func (service ContentService) AddContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	item := &dao.Content{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	have := service.GetContentByNameAndOID(item.Name, company.ID)
	if have.ID != 0 || strings.EqualFold(item.Name, "") {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("这个名字已经被使用了"), "", nil)}
	}

	var mt dao.ContentType
	Orm.Where("ID=?", item.ContentTypeID).First(&mt)
	if mt.ID == 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("没有找到数据"), "", nil)}
	}

	item.OID = company.ID
	item.Type = mt.Type
	err = service.Add(Orm, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (service ContentService) GetContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Content{}
	err := service.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (service ContentService) ListContentsAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var dts []dao.Content
	Orm.Model(dao.Content{}).Where("OID=?", company.ID).Order("Sort").Find(&dts)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", dts)}
}

func (service ContentService) DeleteContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)

	css := service.FindContentSubTypesByContentID(ID)
	if len(css) > 0 {

		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("包含子项内容无法删除"), "删除成功", nil)}
	}
	item := &dao.Content{}
	err := service.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (service ContentService) ChangeContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Content{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	m := service.GetContentByNameAndOID(item.Name, company.ID)
	if m.ID != 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("名字重复，修改失败"), "", nil)}
	}
	err = service.ChangeModel(Orm, ID, &dao.Content{Name: item.Name, Sort: item.Sort})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (service ContentService) ChangeContentIndexAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Content{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = service.ChangeMap(Orm, ID, &dao.Content{}, map[string]interface{}{"Sort": item.Sort})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "index成功", nil)}
}
func (service ContentService) ChangeHideContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Content{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = service.ChangeMap(Orm, ID, &dao.Content{}, map[string]interface{}{"Hide": item.Hide})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "index成功", nil)}
}
