package web

import (
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/sites"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool/number"
)

type Controller struct {
	gweb.BaseController
	Template sites.TemplateService
	Content  service.ContentService
}

func (controller *Controller) Init() {
	controller.AddHandler(gweb.ALLMethod("/index", controller.index))
	controller.AddHandler(gweb.ALLMethod("/gallery/{ID}/{SubID}/{SubChildID}", controller.gallery))
	controller.AddHandler(gweb.GETMethod("/contents/{ID}/{SubID}/{SubChildID}", controller.contents))
	controller.AddHandler(gweb.ALLMethod("/contents/detail/{ID}/{SubID}/{SubChildID}/{ContentID}", controller.contentsDetail))
	controller.AddHandler(gweb.ALLMethod("/content/{ID}/{SubID}/{SubChildID}", controller.content))
	controller.AddHandler(gweb.ALLMethod("/products/{ID}/{SubID}/{SubChildID}", controller.products))
	controller.AddHandler(gweb.ALLMethod("/js/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/css/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/img/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/font/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.ALLMethod("/lib/", controller.AddProjectdsfdsfsdAction))
	//controller.AddHandler(gweb.ALLMethod("/", controller.defaultPage))
}

func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {

	//return &gweb.HTMLResult{}
	return &gweb.RedirectToUrlResult{Url: "index"}

}
func (controller *Controller) products(context *gweb.Context) gweb.Result {
	params := make(map[string]interface{})

	//{Index}/{SubIndex}/{SubChildIndex}
	ID := number.ParseInt(context.PathParams["ID"])
	SubID := number.ParseInt(context.PathParams["SubID"])
	SubChildID := number.ParseInt(context.PathParams["SubChildID"])

	menusData, menusPath := controller.Template.MenusTemplate(context)
	menusData.SetCurrentMenus(uint64(ID), uint64(SubID), uint64(SubChildID))
	//{ID}/{SubID}/{SubChildID}
	//menus := menusData.Get(uint64(ID))
	//menusSub := menus.Get(uint64(SubID))
	//menusSubChild := menusSub.Get(uint64(SubChildID))

	contentList := controller.Content.FindContentListByTypeID(&menusData, uint64(ID), uint64(SubID), uint64(SubChildID))
	params["ContentList"] = contentList

	commonPath := controller.Template.CommonTemplate(context, params)
	params["Menus"] = menusData
	return &gweb.HTMLResult{
		Name: menusData.Top.Item.TemplateName,
		Template: []string{
			menusPath, commonPath,
		},
		Params: params,
	}
}
func (controller *Controller) content(context *gweb.Context) gweb.Result {
	params := make(map[string]interface{})

	ID := number.ParseInt(context.PathParams["ID"])
	SubID := number.ParseInt(context.PathParams["SubID"])
	SubChildID := number.ParseInt(context.PathParams["SubChildID"])

	menusData, menusPath := controller.Template.MenusTemplate(context)
	menusData.SetCurrentMenus(uint64(ID), uint64(SubID), uint64(SubChildID))

	content := controller.Content.FindContentByTypeID(&menusData, uint64(ID), uint64(SubID), uint64(SubChildID))
	params["Content"] = content

	commonPath := controller.Template.CommonTemplate(context, params)

	params["Menus"] = menusData

	return &gweb.HTMLResult{
		Name: menusData.Top.Item.TemplateName,
		Template: []string{
			menusPath, commonPath,
		},
		Params: params,
	}
}
func (controller *Controller) contentsDetail(context *gweb.Context) gweb.Result {
	params := make(map[string]interface{})

	ID := number.ParseInt(context.PathParams["ID"])
	SubID := number.ParseInt(context.PathParams["SubID"])
	SubChildID := number.ParseInt(context.PathParams["SubChildID"])

	ContentID := number.ParseInt(context.PathParams["ContentID"])

	content := controller.Content.GetContentByID(uint64(ContentID))

	params["Content"] = content

	//ContentSubType:=controller.Content.GetContentByID(uint64(ContentSubTypeID))

	/*	var ContentSubTypeIDList []uint64
		if ContentSubTypeID == 0 {
			ContentSubTypeIDList = []uint64{}
		} else {
			ContentSubTypeIDList = controller.Content.GetContentSubTypeAllIDByID(uint64(ID), uint64(ContentSubTypeID))
		}*/

	contentList := controller.Content.FindContentListForLeftRight(uint64(ID), uint64(SubID), uint64(SubChildID), content.ID, content.CreatedAt)

	params["ContentLeft"] = contentList[0]
	params["ContentRight"] = contentList[1]

	menusData, menusPath := controller.Template.MenusTemplate(context)
	menusData.SetCurrentMenus(uint64(ID), uint64(SubID), uint64(SubChildID))
	commonPath := controller.Template.CommonTemplate(context, params)
	params["Menus"] = menusData
	return &gweb.HTMLResult{
		Name: "contents_detail",
		Template: []string{
			menusPath, commonPath,
		},
		Params: params,
	}
}
func (controller *Controller) contents(context *gweb.Context) gweb.Result {
	params := make(map[string]interface{})

	//{Index}/{SubIndex}/{SubChildIndex}
	ID := number.ParseInt(context.PathParams["ID"])
	SubID := number.ParseInt(context.PathParams["SubID"])
	SubChildID := number.ParseInt(context.PathParams["SubChildID"])

	menusData, menusPath := controller.Template.MenusTemplate(context)
	menusData.SetCurrentMenus(uint64(ID), uint64(SubID), uint64(SubChildID))
	//{ID}/{SubID}/{SubChildID}
	//menus := menusData.Get(uint64(ID))
	//menusSub := menus.Get(uint64(SubID))
	//menusSubChild := menusSub.Get(uint64(SubChildID))

	contentList := controller.Content.FindContentListByTypeID(&menusData, uint64(ID), uint64(SubID), uint64(SubChildID))
	params["ContentList"] = contentList

	commonPath := controller.Template.CommonTemplate(context, params)
	params["Menus"] = menusData
	return &gweb.HTMLResult{
		Name: menusData.Top.Item.TemplateName,
		Template: []string{
			menusPath, commonPath,
		},
		Params: params,
	}
}
func (controller *Controller) gallery(context *gweb.Context) gweb.Result {
	params := make(map[string]interface{})

	//{Index}/{SubIndex}/{SubChildIndex}
	ID := number.ParseInt(context.PathParams["ID"])
	SubID := number.ParseInt(context.PathParams["SubID"])
	SubChildID := number.ParseInt(context.PathParams["SubChildID"])

	menusData, menusPath := controller.Template.MenusTemplate(context)
	menusData.SetCurrentMenus(uint64(ID), uint64(SubID), uint64(SubChildID))
	//{ID}/{SubID}/{SubChildID}
	//menus := menusData.Get(uint64(ID))
	//menusSub := menus.Get(uint64(SubID))
	//menusSubChild := menusSub.Get(uint64(SubChildID))

	contentList := controller.Content.FindContentListByTypeID(&menusData, uint64(ID), uint64(SubID), uint64(SubChildID))
	params["ContentList"] = contentList

	commonPath := controller.Template.CommonTemplate(context, params)
	params["Menus"] = menusData
	return &gweb.HTMLResult{
		Name: menusData.Top.Item.TemplateName,
		Template: []string{
			menusPath, commonPath,
		},
		Params: params,
	}
}
func (controller *Controller) index(context *gweb.Context) gweb.Result {

	//siteName := context.PathParams["siteName"]

	params := make(map[string]interface{})

	menusData, menusPath := controller.Template.MenusTemplate(context)
	commonPath := controller.Template.CommonTemplate(context, params)
	params["Menus"] = menusData
	return &gweb.HTMLResult{
		Template: []string{
			menusPath, commonPath,
		},
		Params: params,
	}
}
func (controller *Controller) AddProjectdsfdsfsdAction(context *gweb.Context) gweb.Result {

	return &gweb.FileServerResult{}
}
func (controller *Controller) AddProjectAction(context *gweb.Context) gweb.Result {

	glog.Trace(context.Request.URL)
	//var project dao.Project

	//util.RequestBodyToJSON(context.Request.Body, &project)

	//fmt.Println(project)

	//controller.Dao.Add(service.Orm, &project)

	return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionOK, Message: "信息已经提交，我们会在第一时间联系您。", Data: nil}}
}
