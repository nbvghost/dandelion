package service

import (
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"strconv"

	"time"

	"dandelion/app/play"

	"github.com/nbvghost/gweb"
)

type VoucherService struct {
	dao.BaseDao
}

func (service VoucherService) Situation(StartTime, EndTime int64) interface{} {

	st := time.Unix(StartTime/1000, 0)
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et := time.Unix(EndTime/1000, 0).Add(24 * time.Hour)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	Orm := dao.Orm()

	type Result struct {
		TotalMoney uint64 `gorm:"column:TotalMoney"`
		TotalCount uint64 `gorm:"column:TotalCount"`
	}

	//select COUNT(ID),SUM(JSON_EXTRACT(Data,'$.Amount')) as TotalMoney from carditem where Type='Voucher' group by Type;

	var result Result

	Orm.Table("CardItem").Select("SUM(JSON_EXTRACT(Data,'$.Amount')) as TotalMoney,COUNT(ID) as TotalCount").Where("CreatedAt>=?", st).Where("CreatedAt<?", et).Where("Type=?", "Voucher").Find(&result)
	//fmt.Println(result)
	return result
}
func (service VoucherService) AddItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	item := &dao.Voucher{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	item.OID = company.ID
	err = service.Add(Orm, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (service VoucherService) GetItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Voucher{}
	err := service.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (service VoucherService) ListItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.Voucher{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (service VoucherService) DeleteItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Voucher{}
	err := service.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (service VoucherService) ChangeItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Voucher{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = service.ChangeModel(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
