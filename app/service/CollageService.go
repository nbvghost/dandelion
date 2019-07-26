package service

import (
	"strconv"
	"strings"

	"github.com/nbvghost/glog"

	"dandelion/app/play"
	"dandelion/app/service/dao"
	"dandelion/app/util"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
	"github.com/pkg/errors"
)

type CollageService struct {
	dao.BaseDao
}
func (service CollageService) AddCollageGoodsAction(context *gweb.Context) gweb.Result {
	organization := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	context.Request.ParseForm()

	GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	CollageHash := context.Request.FormValue("CollageHash")

	goods := GlobalService.Goods.FindGoodsByOrganizationIDAndGoodsID(organization.ID, GoodsID)
	if goods.ID == 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("没有找到商品"), "", nil)}
	}
	collage := service.GetCollageByHash(CollageHash, organization.ID)
	if collage.ID == 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("没有找到限时抢购"), "", nil)}
	}

	have := service.GetCollageGoodsByGoodsID(goods.ID, organization.ID)
	if have.ID > 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("这个商品已被添加为限时抢购"), "", nil)}
	}

	//service.ChangeMap(dao.Orm(), timeSell.ID, &dao.TimeSell{}, map[string]interface{}{})
	err := service.Add(dao.Orm(), &dao.CollageGoods{
		CollageHash: collage.Hash,
		GoodsID:      goods.ID,
		Disable:      false,
		OID:          organization.ID,
	})

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "Success", goods)}
}
func (service CollageService) GetCollageGoodsByGoodsID(GoodsID uint64, OID uint64) dao.CollageGoods {
	var timesellGoods dao.CollageGoods
	dao.Orm().Model(&dao.CollageGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)
	return timesellGoods
}

func (service CollageService) DeleteCollage(TimeSellID uint64) error {
	//timesell := TimeSellService{}.GetTimeSellByGoodsID(GoodsID)
	var ts dao.Collage
	service.Get(dao.Orm(), TimeSellID, &ts)
	//err := service.Delete(dao.Orm(), &dao.TimeSell{}, ts.ID)
	err := service.DeleteWhere(dao.Orm(), &dao.Collage{}, "Hash=?", ts.Hash)
	glog.Error(err)

	return err
}

func (service CollageService) DeleteItem(context *gweb.Context) gweb.Result {

	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	err := service.DeleteCollage(ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (service CollageService) DeleteGoods(context *gweb.Context) gweb.Result {

	/*Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	list := GlobalService.Goods.DeleteTimeSellGoods(Orm, ID, company.ID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "删除成功", list)}
	*/

	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	list := GlobalService.Goods.DeleteCollageGoods(Orm, ID,company.ID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "删除成功", list)}

}
func (service CollageService) ListGoods(context *gweb.Context) gweb.Result {

	//Hash := context.PathParams["Hash"]

	//list := GlobalService.Goods.FindGoodsByCollageHash(Hash)
	//var item dao.ExpressTemplate
	//err := controller.ExpressTemplate.Get(service.Orm, ID, &item)
	//return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "", list)}
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
	draw, recordsTotal, recordsFiltered, list := GlobalService.Goods.DatatablesListOrder(dao.Orm(), dts, &[]dao.Goods{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (service CollageService) DataTablesItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	dts.Groupbys = make([]string, 0)
	dts.Groupbys = append(dts.Groupbys, "Hash")

	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.Collage{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (service CollageService) GetItemByHash(Hash string) dao.Collage {
	var timesell dao.Collage
	err := dao.Orm().Model(&dao.Collage{}).Where("Hash=?", Hash).First(&timesell).Error
	glog.Error(err)
	return timesell
}
func (service CollageService) GetItem(context *gweb.Context) gweb.Result {
	//Orm := dao.Orm()
	Hash := context.PathParams["Hash"]
	item := service.GetItemByHash(Hash)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", item)}
}
func (service CollageService) GetCollageByGoodsID(GoodsID uint64, OID uint64) dao.Collage {
	//timesellGoods := service.GetTimeSellGoodsByGoodsID(GoodsID, OID)
	var timesellGoods dao.CollageGoods
	dao.Orm().Model(&dao.CollageGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)

	var timesell dao.Collage
	dao.Orm().Model(&dao.Collage{}).Where("Hash=? and OID=?", timesellGoods.CollageHash, timesellGoods.OID).First(&timesell)
	return timesell

	/*var timesell dao.Collage
	err := dao.Orm().Model(&dao.Collage{}).Where("GoodsID=?", GoodsID).First(&timesell).Error
	glog.Error(err)
	return timesell*/
}
func (service CollageService) GetCollageByHash(Hash string,OID uint64) dao.Collage {
	var timesell dao.Collage
	err := dao.Orm().Model(&dao.Collage{}).Where("Hash=? and OID=?", Hash,OID).First(&timesell).Error
	glog.Error(err)
	return timesell
}
func (service CollageService) AddCollageRecord(OrderNo, OrdersGoodsNo, No string, UserID uint64) error {
	cr := &dao.CollageRecord{}
	cr.No = No
	cr.OrderNo = OrderNo
	cr.UserID = UserID
	cr.OrdersGoodsNo = OrdersGoodsNo
	if strings.EqualFold(No, "") {
		cr.No = tool.UUID()
		cr.Collager = UserID
	} else {
		cr.No = No
		cr.Collager = 0
		_cr := service.FindCollageRecordByUserIDAndNo(UserID, No)
		if _cr.ID != 0 {
			return errors.New("您已经参加了这个活动，看看其它活动吧！")
		}
	}
	return service.Add(dao.Orm(), cr)
}
func (service CollageService) FindCollageRecordByUserIDAndNo(UserID uint64, No string) dao.CollageRecord {
	Orm := dao.Orm()
	var cr dao.CollageRecord
	Orm.Model(&dao.CollageRecord{}).Where("UserID=? and No=?").First(&cr)
	return cr

}
func (service CollageService) SaveItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	context.Request.ParseForm()
	var err error

	CollageJson := context.Request.FormValue("Collage")
	//GoodsListJson := context.Request.FormValue("GoodsListIDs")

	//GoodsListIDs := make([]uint64, 0)
	//err = util.JSONToStruct(GoodsListJson, &GoodsListIDs)
	//if err != nil {
	//	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
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
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
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
		err = service.Save(tx, item)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "提交成功", item)}

	} else {
		//修改
		_item:=service.GetCollageByHash(item.Hash, company.ID)
		if _item.ID == 0 {
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("无法修改"), "", nil)}
		}
		_item.Num=item.Num
		_item.Discount=item.Discount
		_item.TotalNum=item.TotalNum
		err = service.Save(tx, _item)

		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "提交成功", _item)}
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

	//return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "提交成功", nil)}

}
