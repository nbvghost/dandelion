package service

import (
	"errors"
	"strconv"

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


func (service TimeSellService) AddTimeSellGoodsAction(context *gweb.Context) gweb.Result {
	organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	context.Request.ParseForm()

	GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	TimeSellHash := context.Request.FormValue("TimeSellHash")

	goods := GlobalService.Goods.FindGoodsByOrganizationIDAndGoodsID(organization.ID, GoodsID)
	if goods.ID == 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("没有找到商品"), "", nil)}
	}
	timeSell := service.GetTimeSellByHash(TimeSellHash, organization.ID)
	if timeSell.ID == 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("没有找到限时抢购"), "", nil)}
	}

	have := service.GetTimeSellGoodsByGoodsID(goods.ID, organization.ID)
	if have.ID > 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("这个商品已被添加为限时抢购"), "", nil)}
	}

	//service.ChangeMap(dao.Orm(), timeSell.ID, &dao.TimeSell{}, map[string]interface{}{})
	err := service.Add(dao.Orm(), &dao.TimeSellGoods{
		TimeSellHash: timeSell.Hash,
		GoodsID:      goods.ID,
		Disable:      false,
		OID:          organization.ID,
	})

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "Success", goods)}
}
func (service TimeSellService) DeleteTimeSellGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	list := GlobalService.Goods.DeleteTimeSellGoods(Orm, ID, company.ID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "删除成功", list)}

}
func (service TimeSellService) ListTimeSellGoods(context *gweb.Context) gweb.Result {
	Hash := context.PathParams["Hash"]
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	dts := &dao.Datatables{}
	//dts.Draw = 10
	//dts.Length = play.Paging
	util.RequestBodyToJSON(context.Request.Body, dts)
	GoodsIDs := []uint64{}
	dao.Orm().Model(&dao.TimeSellGoods{}).Where("TimeSellHash=? and OID=?", Hash, company.ID).Pluck("GoodsID", &GoodsIDs)
	if len(GoodsIDs) == 0 {
		GoodsIDs = []uint64{0}
	}
	dts.InIDs = GoodsIDs
	draw, recordsTotal, recordsFiltered, list := GlobalService.Goods.DatatablesListOrder(dao.Orm(), dts, &[]dao.Goods{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}

	/*Hash := context.PathParams["Hash"]
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	//dts.Groupbys = make([]string, 0)
	//dts.Groupbys = append(dts.Groupbys, "Hash")

	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.TimeSellGoods{}, company.ID, "TimeSellHash=?", Hash)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}*/

	/*Hash := context.PathParams["Hash"]

	list := GlobalService.Goods.FindGoodsByTimeSellHash(Hash)
	//var item dao.ExpressTemplate
	//err := controller.ExpressTemplate.Get(service.Orm, ID, &item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "", list)}*/
	//2002
}
func (service TimeSellService) GetItem(context *gweb.Context) gweb.Result {
	//Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Hash := context.PathParams["Hash"]
	item := service.GetTimeSellByHash(Hash, company.ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", item)}
}
func (service TimeSellService) GetTimeSellByHash(Hash string, OID uint64) dao.TimeSell {
	var timesell dao.TimeSell
	dao.Orm().Model(&dao.TimeSell{}).Where("Hash=? and OID=?", Hash, OID).First(&timesell)
	return timesell
}

//todo:
func (service TimeSellService) GetTimeSellByGoodsID(GoodsID uint64, OID uint64) *dao.TimeSell {
	//timesellGoods := service.GetTimeSellGoodsByGoodsID(GoodsID, OID)
	var timesellGoods dao.TimeSellGoods
	dao.Orm().Model(&dao.TimeSellGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)

	var timesell dao.TimeSell
	dao.Orm().Model(&dao.TimeSell{}).Where("Hash=? and OID=?", timesellGoods.TimeSellHash, timesellGoods.OID).First(&timesell)
	return &timesell
}
func (service TimeSellService) GetTimeSellGoodsByGoodsID(GoodsID uint64, OID uint64) dao.TimeSellGoods {
	var timesellGoods dao.TimeSellGoods
	dao.Orm().Model(&dao.TimeSellGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)
	return timesellGoods
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
	//GoodsListJson := context.Request.FormValue("GoodsListIDs")

	//GoodsListIDs := make([]uint64, 0)
	//err = util.JSONToStruct(GoodsListJson, &GoodsListIDs)
	//if err != nil {
	//return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	//}

	tx := dao.Orm().Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	item := dao.TimeSell{}
	err = util.JSONToStruct(TimeSellJson, &item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	Hash := tool.UUID()
	if item.ID == 0 {
		//新添加
		/*for _, value := range GoodsListIDs {
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

		}*/

		item := &dao.TimeSell{}
		err = util.JSONToStruct(TimeSellJson, item)
		//item.GoodsID = value
		item.Hash = Hash
		item.OID = company.ID
		err = service.Save(tx, item)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "提交成功", item)}
	} else {
		_item := service.GetTimeSellByHash(item.Hash, company.ID)
		if _item.ID == 0 {
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("无法修改"), "", nil)}
		}
		//_item.Hash = item.Hash
		_item.BuyNum = item.BuyNum
		_item.Enable = item.Enable
		_item.DayNum = item.DayNum
		_item.Discount = item.Discount
		_item.TotalNum = item.TotalNum
		_item.StartTime = item.StartTime
		_item.StartH = item.StartH
		_item.StartM = item.StartM
		_item.EndH = item.EndH
		_item.EndM = item.EndM
		err = service.Save(tx, _item)

		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "提交成功", _item)}

		//修改
		/*for _, value := range GoodsListIDs {
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

		}*/

	}

}

func (service TimeSellService) DataTablesItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	dts.Groupbys = make([]string, 0)
	dts.Groupbys = append(dts.Groupbys, "Hash")

	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.TimeSell{}, company.ID, "")
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
