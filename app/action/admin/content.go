package admin

import (
	"errors"
	"fmt"
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
}

func (controller *ContentController) Init() {

	//------------------ArticleService.go-datatables------------------------
	controller.AddHandler(gweb.POSMethod("datatables/list", controller.DataTablesAction))
	controller.AddHandler(gweb.POSMethod("save", controller.SaveArticleAction))
	controller.AddHandler(gweb.GETMethod("multi/get/{ID}", controller.GetMultiArticleAction))
	controller.AddHandler(gweb.GETMethod("single/get/{ContentItemID}", controller.GetSingleArticleAction))
	controller.AddHandler(gweb.POSMethod("delete", controller.DeleteArticleAction))
	//------------------ArticleService.go-datatables------------------------

	controller.AddHandler(gweb.POSMethod("item/add", controller.AddContentAction))
	controller.AddHandler(gweb.GETMethod("item/{ContentItemID}", controller.GetContentAction))
	controller.AddHandler(gweb.GETMethod("item/list", controller.ListContentsAction))
	controller.AddHandler(gweb.DELMethod("item/{ContentItemID}", controller.DeleteContentAction))
	controller.AddHandler(gweb.PUTMethod("item/{ContentItemID}", controller.ChangeContentAction))
	controller.AddHandler(gweb.PUTMethod("item/index/{ContentItemID}", controller.ChangeContentIndexAction))
	controller.AddHandler(gweb.PUTMethod("item/hide/{ContentItemID}", controller.ChangeHideContentAction))

	controller.AddHandler(gweb.GETMethod("type/list", controller.ListContentTypeAction))
	controller.AddHandler(gweb.POSMethod("sub_type", controller.AddClassify))

	controller.AddHandler(gweb.GETMethod("sub_type/list/all/{ContentItemID}", controller.ListAllSubType))
	controller.AddHandler(gweb.GETMethod("sub_type/list/{ContentItemID}", controller.ListSubType))
	controller.AddHandler(gweb.GETMethod("sub_type/child/list/{ContentItemID}/{ParentContentSubTypeID}", controller.ListChildClassify))

	controller.AddHandler(gweb.DELMethod("sub_type/{ID}", controller.DeleteClassify))
	controller.AddHandler(gweb.PUTMethod("sub_type/{ID}", controller.ChangeClassify))
	controller.AddHandler(gweb.GETMethod("sub_type/{ID}", controller.GetContentSubTypeAction))
}
func (controller *ContentController) DeleteArticleAction(context *gweb.Context) gweb.Result {

	context.Request.ParseForm()
	fmt.Println(context.Request.FormValue("ID"))
	ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
	err := controller.Content.Delete(dao.Orm(), &dao.Content{}, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (controller *ContentController) GetSingleArticleAction(context *gweb.Context) gweb.Result {
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	article := controller.Content.GetArticleByContentItemID(ContentItemID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", article)}
}
func (controller *ContentController) GetMultiArticleAction(context *gweb.Context) gweb.Result {
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	var article dao.Content
	err := controller.Content.Get(dao.Orm(), ID, &article)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", article)}
}
func (controller *ContentController) SaveArticleAction(context *gweb.Context) gweb.Result {

	dts := &dao.Content{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	fmt.Println(dts)

	as := controller.Content.AddArticle(dts)

	return &gweb.JsonResult{Data: as}
}
func (controller *ContentController) DataTablesAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.Content.DatatablesListOrder(Orm, dts, &[]dao.Content{}, 0, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (controller *ContentController) ChangeClassify(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.ContentSubType{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	m := controller.Content.GetClassifyByName(item.Name, item.ContentItemID, item.ParentContentSubTypeID)
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
	articles := controller.Content.FindContentByContentSubTypeID(ID)
	if len(articles) > 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("包含文章，无法删除"), "删除成功", nil)}
	}

	item := &dao.ContentSubType{}
	err := controller.Content.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}

func (controller *ContentController) ListChildClassify(context *gweb.Context) gweb.Result {
	ParentContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ParentContentSubTypeID"], 10, 64)
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	list := controller.Content.FindContentSubTypesByContentItemIDAndParentContentSubTypeID(ContentItemID, ParentContentSubTypeID)

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
func (controller *ContentController) ListSubType(context *gweb.Context) gweb.Result {
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	list := controller.Content.FindContentSubTypesByContentItemID(ContentItemID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}

}
func (controller *ContentController) ListAllSubType(context *gweb.Context) gweb.Result {
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	list := controller.Content.FindContentSubTypesByContentItemID(ContentItemID)

	resultMap := make(map[uint64]interface{})

	for index := range list {
		item := list[index]
		subTypes := controller.Content.FindContentSubTypesByContentItemIDAndParentContentSubTypeID(ContentItemID, item.ID)

		childrenMap := make(map[uint64]interface{})

		for sindex := range subTypes {

			childrenMap[subTypes[sindex].ID] = subTypes[sindex]

		}

		resultMap[item.ID] = map[string]interface{}{
			"SubType":         item,
			"SubTypeChildren": childrenMap,
		}

	}

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", resultMap)}

}
func (controller *ContentController) AddClassify(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	item := &dao.ContentSubType{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	have := controller.Content.GetClassifyByName(item.Name, item.ContentItemID, item.ParentContentSubTypeID)
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
	ID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	item := &dao.ContentItem{}
	err := controller.Content.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (controller *ContentController) ListContentsAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var dts []dao.ContentItem
	Orm.Model(dao.ContentItem{}).Where("OID=?", company.ID).Order("Sort").Find(&dts)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", dts)}
}

func (controller *ContentController) DeleteContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)

	css := controller.Content.FindContentSubTypesByContentItemID(ContentItemID)
	if len(css) > 0 {

		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("包含子项内容无法删除"), "删除成功", nil)}
	}
	item := &dao.ContentItem{}
	err := controller.Content.Delete(Orm, item, ContentItemID)
	if !glog.Error(err) {
		err = controller.Content.DeleteWhere(Orm, &dao.Content{}, "ContentItemID=? and ContentSubTypeID=?", ContentItemID, 0)
	}
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (controller *ContentController) ChangeContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	ID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	item := &dao.ContentItem{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	m := controller.Content.GetContentItemByNameAndOID(item.Name, company.ID)
	if m.ID != 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("名字重复，修改失败"), "", nil)}
	}
	err = controller.Content.ChangeModel(Orm, ID, &dao.ContentItem{Name: item.Name, Sort: item.Sort})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (controller *ContentController) ChangeContentIndexAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	item := &dao.ContentItem{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = controller.Content.ChangeMap(Orm, ID, &dao.ContentItem{}, map[string]interface{}{"Sort": item.Sort})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "index成功", nil)}
}
func (controller *ContentController) ChangeHideContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	item := &dao.ContentItem{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = controller.Content.ChangeMap(Orm, ID, &dao.ContentItem{}, map[string]interface{}{"Hide": item.Hide})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "index成功", nil)}
}
func (controller *ContentController) AddContentAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	item := &dao.ContentItem{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	have := controller.Content.GetContentItemByNameAndOID(item.Name, company.ID)
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
