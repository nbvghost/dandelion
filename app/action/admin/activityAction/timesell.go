package activityAction

import (
	"errors"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/activity"
	"github.com/nbvghost/dandelion/app/service/content"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/goods"
	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
	"strconv"
)

type TimesellController struct {
	gweb.BaseController
	Content  content.ContentService
	Goods    goods.GoodsService
	TimeSell activity.TimeSellService
}

func (controller *TimesellController) Init() {

	controller.AddHandler(gweb.POSMethod("save", controller.SaveItem))
	controller.AddHandler(gweb.POSMethod("change", controller.SaveItem))
	controller.AddHandler(gweb.GETMethod("{Hash}", controller.GetItem))
	controller.AddHandler(gweb.POSMethod("datatables/list", controller.DataTablesItem))
	controller.AddHandler(gweb.POSMethod("goods/{Hash}/list", controller.ListTimeSellGoods))
	controller.AddHandler(gweb.DELMethod("goods/{GoodsID}", controller.DeleteTimeSellGoods))
	controller.AddHandler(gweb.POSMethod("goods/add", controller.AddTimeSellGoodsAction))
	controller.AddHandler(gweb.DELMethod("{ID}", controller.DeleteItem))

}

func (controller *TimesellController) DeleteItem(context *gweb.Context) gweb.Result {

	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	err := controller.TimeSell.DeleteTimeSell(ID)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}
}

func (controller *TimesellController) AddTimeSellGoodsAction(context *gweb.Context) gweb.Result {
	organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	context.Request.ParseForm()

	GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	TimeSellHash := context.Request.FormValue("TimeSellHash")

	goods := controller.Goods.FindGoodsByOrganizationIDAndGoodsID(organization.ID, GoodsID)
	if goods.ID == 0 {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("没有找到商品"), "", nil)}
	}
	timeSell := controller.TimeSell.GetTimeSellByHash(TimeSellHash, organization.ID)
	if timeSell.ID == 0 {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("没有找到限时抢购"), "", nil)}
	}

	have := controller.TimeSell.GetTimeSellGoodsByGoodsID(goods.ID, organization.ID)
	if have.ID > 0 {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("这个商品已被添加为限时抢购"), "", nil)}
	}

	//service.ChangeMap(dao.Orm(), timeSell.ID, &dao.TimeSell{}, map[string]interface{}{})
	err := controller.TimeSell.Add(dao.Orm(), &dao.TimeSellGoods{
		TimeSellHash: timeSell.Hash,
		GoodsID:      goods.ID,
		Disable:      false,
		OID:          organization.ID,
	})

	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "Success", goods)}
}

func (controller *TimesellController) ListTimeSellGoods(context *gweb.Context) gweb.Result {
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
	draw, recordsTotal, recordsFiltered, list := controller.TimeSell.DatatablesListOrder(dao.Orm(), dts, &[]dao.Goods{}, company.ID, "")
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
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "", list)}*/
	//2002
}
func (controller *TimesellController) DataTablesItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	dts.Groupbys = make([]string, 0)
	dts.Groupbys = append(dts.Groupbys, "Hash")

	draw, recordsTotal, recordsFiltered, list := controller.TimeSell.DatatablesListOrder(Orm, dts, &[]dao.TimeSell{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (controller *TimesellController) GetItem(context *gweb.Context) gweb.Result {
	//Orm := dao.Orm()
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Hash := context.PathParams["Hash"]
	item := controller.TimeSell.GetTimeSellByHash(Hash, company.ID)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", item)}
}
func (controller *TimesellController) DeleteTimeSellGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	list := controller.Goods.DeleteTimeSellGoods(Orm, ID, company.ID)

	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "删除成功", list)}

}
func (controller *TimesellController) SaveItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	context.Request.ParseForm()
	var err error

	TimeSellJson := context.Request.FormValue("TimeSell")
	//GoodsListJson := context.Request.FormValue("GoodsListIDs")

	//GoodsListIDs := make([]uint64, 0)
	//err = util.JSONToStruct(GoodsListJson, &GoodsListIDs)
	//if err != nil {
	//return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
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
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
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
		err = controller.TimeSell.Save(tx, item)
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", item)}
	} else {
		_item := controller.TimeSell.GetTimeSellByHash(item.Hash, company.ID)
		if _item.ID == 0 {
			return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无法修改"), "", nil)}
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
		err = controller.TimeSell.Save(tx, _item)

		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", _item)}

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
