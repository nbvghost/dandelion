package admin

import (
	"errors"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
	"strconv"
	"strings"
)

type ContentController struct {
	gweb.BaseController
	Content service.ContentService
	Article service.ArticleService
}

func (controller *ContentController) Init() {
	controller.AddHandler(gweb.POSMethod("add", controller.AddContentAction))
	controller.AddHandler(gweb.GETMethod("{ID}", controller.GetContentAction))
	controller.AddHandler(gweb.GETMethod("list", controller.ListContentsAction))
	controller.AddHandler(gweb.DELMethod("{ID}", controller.DeleteContentAction))
	controller.AddHandler(gweb.PUTMethod("{ID}", controller.ChangeContentAction))
	controller.AddHandler(gweb.PUTMethod("index/{ID}", controller.ChangeContentIndexAction))
	controller.AddHandler(gweb.PUTMethod("hide/{ID}", controller.ChangeHideContentAction))
	controller.AddHandler(gweb.GETMethod("type/list", controller.ListContentTypeAction))
	controller.AddHandler(gweb.POSMethod("sub_type", controller.AddClassify))
	controller.AddHandler(gweb.GETMethod("sub_type/list/{ContentID}", controller.ListClassify))
	controller.AddHandler(gweb.GETMethod("sub_type/child/list/{ContentID}/{ParentContentSubTypeID}", controller.ListChildClassify))
	controller.AddHandler(gweb.DELMethod("sub_type/{ID}", controller.DeleteClassify))
	controller.AddHandler(gweb.PUTMethod("sub_type/{ID}", controller.ChangeClassify))
	controller.AddHandler(gweb.GETMethod("sub_type/{ID}", controller.GetContentSubTypeAction))
}
func (controller *ContentController) ChangeClassify(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.ContentSubType{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	m := controller.Content.GetClassifyByName(item.Name, item.ContentID, item.ParentContentSubTypeID)
	if m.ID != 0 && m.ID != item.ID {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("名字重复，修改失败"), "", nil)}
	}
	err = controller.Content.ChangeModel(Orm, ID, &dao.ContentSubType{Name: item.Name})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (controller *ContentController) DeleteClassify(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	css := controller.Content.FindContentSubTypesByParentContentSubTypeID(ID)
	if len(css) > 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("包含子项内容，无法删除"), "删除成功", nil)}
	}
	articles := controller.Article.FindArticleByContentSubTypeID(ID)
	if len(articles) > 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("包含文章，无法删除"), "删除成功", nil)}
	}

	item := &dao.ContentSubType{}
	err := controller.Content.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}

func (controller *ContentController) ListChildClassify(context *gweb.Context) gweb.Result {
	ParentContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ParentContentSubTypeID"], 10, 64)
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	list := controller.Content.FindContentSubTypesByContentIDAndParentContentSubTypeID(ContentID, ParentContentSubTypeID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}

}

func (controller *ContentController) GetContentSubTypeAction(context *gweb.Context) gweb.Result {
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

func (controller *ContentController) ListClassify(context *gweb.Context) gweb.Result {
	ContentID, _ := strconv.ParseUint(context.PathParams["ContentID"], 10, 64)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	list := controller.Content.FindContentSubTypesByContentID(ContentID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}

}
func (controller *ContentController) AddClassify(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	item := &dao.ContentSubType{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	have := controller.Content.GetClassifyByName(item.Name, item.ContentID, item.ParentContentSubTypeID)
	if have.ID != 0 && have.ID != item.ID {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("这个名字已经被使用了"), "", nil)}
	}

	//item.OID = company.ID
	err = controller.Content.Add(Orm, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (controller *ContentController) ListContentTypeAction(context *gweb.Context) gweb.Result {
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: controller.Content.ListContentType()}}
}
func (controller *ContentController) GetContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Content{}
	err := controller.Content.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (controller *ContentController) ListContentsAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var dts []dao.Content
	Orm.Model(dao.Content{}).Where("OID=?", company.ID).Order("Sort").Find(&dts)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", dts)}
}

func (controller *ContentController) DeleteContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)

	css := controller.Content.FindContentSubTypesByContentID(ID)
	if len(css) > 0 {

		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("包含子项内容无法删除"), "删除成功", nil)}
	}
	item := &dao.Content{}
	err := controller.Content.Delete(Orm, item, ID)
	if !glog.Error(err) {
		err = controller.Content.DeleteWhere(Orm, &dao.Article{}, "ContentID=? and ContentSubTypeID=? and ContentSubTypeChildID=?", ID, 0, 0)
	}
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (controller *ContentController) ChangeContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Content{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	m := controller.Content.GetContentByNameAndOID(item.Name, company.ID)
	if m.ID != 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("名字重复，修改失败"), "", nil)}
	}
	err = controller.Content.ChangeModel(Orm, ID, &dao.Content{Name: item.Name, Sort: item.Sort})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (controller *ContentController) ChangeContentIndexAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Content{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = controller.Content.ChangeMap(Orm, ID, &dao.Content{}, map[string]interface{}{"Sort": item.Sort})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "index成功", nil)}
}
func (controller *ContentController) ChangeHideContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Content{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = controller.Content.ChangeMap(Orm, ID, &dao.Content{}, map[string]interface{}{"Hide": item.Hide})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "index成功", nil)}
}
func (controller *ContentController) AddContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	item := &dao.Content{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	have := controller.Content.GetContentByNameAndOID(item.Name, company.ID)
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
	err = controller.Content.Add(Orm, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
