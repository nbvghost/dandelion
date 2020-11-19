package activityAction

import (
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/activity"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/gweb"
	"strconv"
)

type FullcutController struct {
	gweb.BaseController
	FullCut activity.FullCutService
}

func (controller *FullcutController) Init() {

	controller.AddHandler(gweb.POSMethod("save", controller.SaveItem))
	controller.AddHandler(gweb.GETMethod("{ID}", controller.GetItem))
	controller.AddHandler(gweb.POSMethod("datatables/list", controller.DataTablesItem))
	controller.AddHandler(gweb.DELMethod("{ID}", controller.DeleteItem))

}
func (controller *FullcutController) DeleteItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.FullCut{}
	err := controller.FullCut.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}
}

func (controller *FullcutController) DataTablesItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.FullCut.DatatablesListOrder(Orm, dts, &[]dao.FullCut{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (controller *FullcutController) GetItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.FullCut{}
	err := controller.FullCut.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", item)}
}
func (controller *FullcutController) SaveItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	item := &dao.FullCut{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	}
	item.OID = company.ID
	if Orm.NewRecord(item) {
		err = controller.FullCut.Add(Orm, item)
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}
	} else {
		err = controller.FullCut.Save(Orm, item)
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}
	}

}
