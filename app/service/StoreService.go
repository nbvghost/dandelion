package service

import (
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"strconv"

	"dandelion/app/play"
	"errors"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
)

type StoreService struct {
	dao.BaseDao
}

func (service StoreService) GetByPhone(Phone string) dao.Store {
	Orm := dao.Orm()
	var store dao.Store
	Orm.Model(&dao.Store{}).Where(&dao.Store{Phone: Phone}).First(&store)
	return store
}
func (service StoreService) LocationList(Latitude, Longitude float64) []map[string]interface{} {
	Orm := dao.Orm()

	rows, err := Orm.Model(&dao.Store{}).Select("ID,Images,Name,Address,ServicePhone,Stars,StarsCount,ROUND(6378.138*2*ASIN(SQRT(POW(SIN((?*PI()/180-Latitude*PI()/180)/2),2)+COS(?*PI()/180)*COS(Latitude*PI()/180)*POW(SIN((?*PI()/180-Longitude*PI()/180)/2),2)))*1000) AS Distance", Latitude, Latitude, Longitude).Order("Distance asc").Rows()
	glog.Error(err)
	defer rows.Close()

	list := make([]map[string]interface{}, 0)
	for rows.Next() {

		var ID uint64
		var Images string
		var Name string
		var Address string
		var ServicePhone string
		var Stars uint64
		var StarsCount uint64
		var Distance float64

		err = rows.Scan(&ID, &Images, &Name, &Address, &ServicePhone, &Stars, &StarsCount, &Distance)
		glog.Error(err)
		list = append(list, map[string]interface{}{"ID": ID, "Images": Images, "Name": Name, "Address": Address, "ServicePhone": ServicePhone, "Stars": Stars, "StarsCount": StarsCount, "Distance": Distance})
	}

	return list
}
func (service StoreService) AddItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	item := &dao.Store{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	var _store dao.Store
	_store = service.GetByPhone(item.Phone)
	if _store.ID > 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("手机号："+item.Phone+"已经被使用"), "", nil)}
	}

	err = service.Add(Orm, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (service StoreService) GetItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Store{}
	err := service.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (service StoreService) ListItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.Store{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (service StoreService) DeleteItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Store{}
	err := service.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (service StoreService) ChangeItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Store{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	var _store dao.Store
	service.GetByPhone(item.Phone)
	if _store.ID > 0 && _store.ID != item.ID {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("手机号："+item.Phone+"已经被使用"), "", nil)}
	}

	err = service.ChangeModel(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
