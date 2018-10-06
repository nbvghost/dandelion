package service

import (
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"strconv"

	"dandelion/app/play"

	"github.com/nbvghost/gweb"
)

type TimeSellService struct {
	dao.BaseDao
	GoodsService GoodsService
}

func (service TimeSellService) SaveItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	context.Request.ParseForm()

	TimeSell := context.Request.FormValue("TimeSell")
	GoodsList := context.Request.FormValue("GoodsList")

	//fmt.Println(TimeSell)
	//fmt.Println(GoodsList)

	//form.TimeSell=JSON.stringify($scope.Item);
	//form.GoodsList=JSON.stringify($scope.GoodsList);

	item := &dao.TimeSell{}
	err := util.JSONToStruct(TimeSell, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	var list []dao.Goods
	err = util.JSONToStruct(GoodsList, &list)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	tx := Orm.Begin()

	successTxt := ""
	if Orm.NewRecord(item) {
		err = service.Add(tx, item)
		successTxt = "添加成功"
	} else {
		err = service.Save(tx, item)
		successTxt = "修改成功"
		//return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
	}

	for _, value := range list {

		err = service.ChangeModel(tx, value.ID, &dao.Goods{TimeSellID: item.ID})

	}

	defer func() {

		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, successTxt, nil)}

}
func (service TimeSellService) GetItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.TimeSell{}
	err := service.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (service TimeSellService) DataTablesItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.TimeSell{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (service TimeSellService) DeleteItem(context *gweb.Context) gweb.Result {

	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.TimeSell{}
	err := service.Delete(Orm, item, ID)

	list := service.GoodsService.FindByTimeSellID(ID)
	for _, value := range list {
		err = service.GoodsService.DeleteTimeSellGoods(Orm, value.ID)
	}

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
