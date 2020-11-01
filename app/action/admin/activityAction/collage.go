package activityAction

import (
	"errors"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/activity"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/goods"
	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
	"strconv"
)

type CollageController struct {
	gweb.BaseController
	Collage activity.CollageService

	Goods goods.GoodsService
}

func (controller *CollageController) Init() {

	controller.AddHandler(gweb.POSMethod("save", controller.SaveItem))
	controller.AddHandler(gweb.POSMethod("change", controller.SaveItem))
	controller.AddHandler(gweb.GETMethod("{Hash}", controller.GetItem))
	controller.AddHandler(gweb.POSMethod("datatables/list", controller.DataTablesItem))
	controller.AddHandler(gweb.POSMethod("goods/{Hash}/list", controller.ListGoods))
	controller.AddHandler(gweb.DELMethod("goods/{GoodsID}", controller.DeleteGoods))
	controller.AddHandler(gweb.POSMethod("goods/add", controller.AddCollageGoodsAction))
	controller.AddHandler(gweb.DELMethod("{ID}", controller.DeleteItem))

}

func (controller *CollageController) DeleteItem(context *gweb.Context) gweb.Result {

	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	err := controller.Collage.DeleteCollage(ID)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}
}

func (controller *CollageController) AddCollageGoodsAction(context *gweb.Context) gweb.Result {
	organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	context.Request.ParseForm()

	GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	CollageHash := context.Request.FormValue("CollageHash")

	goods := controller.Goods.FindGoodsByOrganizationIDAndGoodsID(organization.ID, GoodsID)
	if goods.ID == 0 {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("没有找到商品"), "", nil)}
	}
	collage := controller.Collage.GetCollageByHash(CollageHash, organization.ID)
	if collage.ID == 0 {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("没有找到限时抢购"), "", nil)}
	}

	have := controller.Collage.GetCollageGoodsByGoodsID(goods.ID, organization.ID)
	if have.ID > 0 {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("这个商品已被添加为限时抢购"), "", nil)}
	}

	//service.ChangeMap(dao.Orm(), timeSell.ID, &dao.TimeSell{}, map[string]interface{}{})
	err := controller.Collage.Add(dao.Orm(), &dao.CollageGoods{
		CollageHash: collage.Hash,
		GoodsID:     goods.ID,
		Disable:     false,
		OID:         organization.ID,
	})

	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "Success", goods)}
}
func (controller *CollageController) DeleteGoods(context *gweb.Context) gweb.Result {

	/*Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	list := GlobalService.Goods.DeleteTimeSellGoods(Orm, ID, company.ID)

	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "删除成功", list)}
	*/

	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	list := controller.Goods.DeleteCollageGoods(Orm, ID, company.ID)

	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "删除成功", list)}

}
func (controller *CollageController) ListGoods(context *gweb.Context) gweb.Result {

	//Hash := context.PathParams["Hash"]

	//list := GlobalService.Goods.FindGoodsByCollageHash(Hash)
	//var item dao.ExpressTemplate
	//err := controller.ExpressTemplate.Get(service.Orm, ID, &item)
	//return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "", list)}
	//2002

	Hash := context.PathParams["Hash"]
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	dts := &dao.Datatables{}
	//dts.Draw = 10
	//dts.Length = play.Paging
	util.RequestBodyToJSON(context.Request.Body, dts)
	GoodsIDs := []uint64{}
	dao.Orm().Model(&dao.CollageGoods{}).Where("CollageHash=? and OID=?", Hash, company.ID).Pluck("GoodsID", &GoodsIDs)
	if len(GoodsIDs) == 0 {
		GoodsIDs = []uint64{0}
	}
	dts.InIDs = GoodsIDs
	draw, recordsTotal, recordsFiltered, list := controller.Goods.DatatablesListOrder(dao.Orm(), dts, &[]dao.Goods{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (controller *CollageController) DataTablesItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	dts.Groupbys = make([]string, 0)
	dts.Groupbys = append(dts.Groupbys, "Hash")

	draw, recordsTotal, recordsFiltered, list := controller.Collage.DatatablesListOrder(Orm, dts, &[]dao.Collage{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (controller *CollageController) GetItem(context *gweb.Context) gweb.Result {
	//Orm := dao.Orm()
	Hash := context.PathParams["Hash"]
	item := controller.Collage.GetItemByHash(Hash)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", item)}
}

func (controller *CollageController) SaveItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	context.Request.ParseForm()
	var err error

	CollageJson := context.Request.FormValue("Collage")
	//GoodsListJson := context.Request.FormValue("GoodsListIDs")

	//GoodsListIDs := make([]uint64, 0)
	//err = util.JSONToStruct(GoodsListJson, &GoodsListIDs)
	//if err != nil {
	//	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	//}

	tx := dao.Orm().Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	item := &dao.Collage{}
	err = util.JSONToStruct(CollageJson, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	}
	Hash := tool.UUID()
	//if strings.EqualFold(item.Hash, "") {
	if item.ID == 0 {
		//新添加
		/*for _, value := range GoodsListIDs {
			isHaveTS := service.GetCollageByGoodsID(value)
			if isHaveTS.ID != 0 && isHaveTS.OID != company.ID {
				continue
			}

			item := &dao.Collage{}
			err = util.JSONToStruct(CollageJson, item)
			item.GoodsID = value
			item.Hash = Hash
			item.OID = company.ID
			err = service.Save(tx, item)

		}*/

		item := &dao.Collage{}
		err = util.JSONToStruct(CollageJson, item)
		//item.GoodsID = value
		item.Hash = Hash
		item.OID = company.ID
		err = controller.Collage.Save(tx, item)
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", item)}

	} else {
		//修改
		_item := controller.Collage.GetCollageByHash(item.Hash, company.ID)
		if _item.ID == 0 {
			return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无法修改"), "", nil)}
		}
		_item.Num = item.Num
		_item.Discount = item.Discount
		_item.TotalNum = item.TotalNum
		err = controller.Collage.Save(tx, _item)

		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", _item)}
		/*for _, value := range GoodsListIDs {
			isHaveTS := service.GetCollageByGoodsID(value)
			if isHaveTS.ID != 0 {
				if strings.EqualFold(item.Hash, isHaveTS.Hash) && isHaveTS.OID == company.ID {
					_item := &dao.Collage{}
					err = util.JSONToStruct(CollageJson, _item)
					_item.GoodsID = value
					_item.Hash = item.Hash
					_item.OID = company.ID
					_item.ID = isHaveTS.ID
					err = service.Save(tx, _item)
				}
				continue
			}

			_item := &dao.Collage{}
			err = util.JSONToStruct(CollageJson, _item)
			_item.GoodsID = value
			_item.Hash = item.Hash
			_item.OID = company.ID
			_item.ID = 0
			err = service.Save(tx, _item)

		}*/

	}

	//return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", nil)}

}
