package service

import (
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"strconv"

	"dandelion/app/play"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb"
)

type FullCutService struct {
	dao.BaseDao
}

func (service FullCutService) FindOrderByAmountDesc(DB *gorm.DB) []dao.FullCut {
	var list []dao.FullCut
	DB.Model(&dao.FullCut{}).Order("Amount desc").Find(&list)
	return list
}
func (service FullCutService) FindOrderByAmountASC(DB *gorm.DB) []dao.FullCut {
	var list []dao.FullCut
	DB.Model(&dao.FullCut{}).Order("Amount asc").Find(&list)
	return list
}
func (service FullCutService) SaveItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	item := &dao.FullCut{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	if Orm.NewRecord(item) {
		err = service.Add(Orm, item)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
	} else {
		err = service.Save(Orm, item)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
	}

}
func (service FullCutService) GetItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.FullCut{}
	err := service.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (service FullCutService) DataTablesItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.FullCut{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (service FullCutService) DeleteItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.FullCut{}
	err := service.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
