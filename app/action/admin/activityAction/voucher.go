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

type VoucherController struct {
	gweb.BaseController
	Voucher activity.VoucherService
}

func (controller *VoucherController) Init() {

	controller.AddHandler(gweb.POSMethod("default", controller.AddItem))
	controller.AddHandler(gweb.GETMethod("{ID}", controller.GetItem))
	controller.AddHandler(gweb.POSMethod("list", controller.ListItem))
	controller.AddHandler(gweb.DELMethod("{ID}", controller.DeleteItem))
	controller.AddHandler(gweb.PUTMethod("{ID}", controller.ChangeItem))

}

func (controller *VoucherController) AddItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	item := &dao.Voucher{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	}
	item.OID = company.ID
	err = controller.Voucher.Add(Orm, item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}
}
func (controller *VoucherController) GetItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Voucher{}
	err := controller.Voucher.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", item)}
}
func (controller *VoucherController) ListItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.Voucher.DatatablesListOrder(Orm, dts, &[]dao.Voucher{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (controller *VoucherController) DeleteItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Voucher{}
	err := controller.Voucher.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}
}
func (controller *VoucherController) ChangeItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Voucher{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	}
	err = controller.Voucher.ChangeModel(Orm, ID, item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}
}
