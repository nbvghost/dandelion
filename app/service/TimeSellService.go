package service

import (
	"strconv"
	"strings"

	"dandelion/app/service/dao"
	"dandelion/app/util"

	"dandelion/app/play"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type TimeSellService struct {
	dao.BaseDao
}

func (service TimeSellService) DeleteTimeSellGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)

	list := GlobalService.Goods.DeleteTimeSellGoods(Orm, ID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "删除成功", list)}

}
func (service TimeSellService) ListTimeSellGoods(context *gweb.Context) gweb.Result {

	Hash := context.PathParams["Hash"]

	list := GlobalService.Goods.FindGoodsByTimeSellHash(Hash)
	//var item dao.ExpressTemplate
	//err := controller.ExpressTemplate.Get(service.Orm, ID, &item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "", list)}
	//2002
}
func (service TimeSellService) GetItem(context *gweb.Context) gweb.Result {
	//Orm := dao.Orm()
	Hash := context.PathParams["Hash"]
	item := service.GetTimeSellByHash(Hash)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", item)}
}
func (service TimeSellService) GetTimeSellByHash(Hash string) dao.TimeSell {
	var timesell dao.TimeSell
	err := dao.Orm().Model(&dao.TimeSell{}).Where("Hash=?", Hash).First(&timesell).Error
	glog.Error(err)
	return timesell
}
func (service TimeSellService) GetTimeSellListByHash(Hash string) []dao.TimeSell {
	var timesells []dao.TimeSell
	err := dao.Orm().Model(&dao.TimeSell{}).Where("Hash=?", Hash).Find(&timesells).Error
	glog.Error(err)
	return timesells
}
func (service TimeSellService) GetTimeSellByGoodsID(GoodsID uint64) dao.TimeSell {
	var timesell dao.TimeSell
	err := dao.Orm().Model(&dao.TimeSell{}).Where("GoodsID=?", GoodsID).First(&timesell).Error
	glog.Error(err)
	return timesell
}

/*
func (service TimeSellService) AddTimeSellAction(context *gweb.Context) gweb.Result {
	//:Hash/:GoodsID
	context.Request.ParseForm()
	//Hash := context.Request.FormValue("Hash")
	//GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "", nil)}
}*/
func (service TimeSellService) SaveItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	context.Request.ParseForm()
	var err error

	TimeSellJson := context.Request.FormValue("TimeSell")
	GoodsListJson := context.Request.FormValue("GoodsListIDs")

	GoodsListIDs := make([]uint64, 0)
	err = util.JSONToStruct(GoodsListJson, &GoodsListIDs)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	tx := dao.Orm().Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	item := &dao.TimeSell{}
	err = util.JSONToStruct(TimeSellJson, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	Hash := tool.UUID()
	if strings.EqualFold(item.Hash, "") {
		//新添加
		for _, value := range GoodsListIDs {
			isHaveTS := service.GetTimeSellByGoodsID(value)
			if isHaveTS.ID != 0 && isHaveTS.OID != company.ID {
				continue
			}

			item := &dao.TimeSell{}
			err = util.JSONToStruct(TimeSellJson, item)
			item.GoodsID = value
			item.Hash = Hash
			item.OID = company.ID
			err = service.Save(tx, item)

		}

	} else {
		//修改
		for _, value := range GoodsListIDs {
			isHaveTS := service.GetTimeSellByGoodsID(value)
			if isHaveTS.ID != 0 {
				if strings.EqualFold(item.Hash, isHaveTS.Hash) && isHaveTS.OID == company.ID {
					_item := &dao.TimeSell{}
					err = util.JSONToStruct(TimeSellJson, _item)
					_item.GoodsID = value
					_item.Hash = item.Hash
					_item.OID = company.ID
					_item.ID = isHaveTS.ID
					err = service.Save(tx, _item)
				}
				continue
			}

			_item := &dao.TimeSell{}
			err = util.JSONToStruct(TimeSellJson, _item)
			_item.GoodsID = value
			_item.Hash = item.Hash
			_item.OID = company.ID
			_item.ID = 0
			err = service.Save(tx, _item)

		}

	}

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "提交成功", nil)}

}

func (service TimeSellService) DataTablesItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	dts.Groupbys = make([]string, 0)
	dts.Groupbys = append(dts.Groupbys, "Hash")

	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.TimeSell{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (service TimeSellService) DeleteItem(context *gweb.Context) gweb.Result {

	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	err := service.DeleteTimeSell(ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}

func (service TimeSellService) DeleteTimeSell(TimeSellID uint64) error {
	//timesell := TimeSellService{}.GetTimeSellByGoodsID(GoodsID)
	var ts dao.TimeSell
	service.Get(dao.Orm(), TimeSellID, &ts)
	//err := service.Delete(dao.Orm(), &dao.TimeSell{}, ts.ID)
	err := service.DeleteWhere(dao.Orm(), &dao.TimeSell{}, "Hash=?", ts.Hash)
	glog.Error(err)
	return err
}
