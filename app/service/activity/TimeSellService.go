package activity

import (
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/glog"
)

type TimeSellService struct {
	dao.BaseDao
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
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "", nil)}
}*/

func (service TimeSellService) DeleteTimeSell(TimeSellID uint64) error {
	//timesell := TimeSellService{}.GetTimeSellByGoodsID(GoodsID)
	var ts dao.TimeSell
	service.Get(dao.Orm(), TimeSellID, &ts)
	//err := service.Delete(dao.Orm(), &dao.TimeSell{}, ts.ID)
	err := service.DeleteWhere(dao.Orm(), &dao.TimeSell{}, "Hash=?", ts.Hash)
	glog.Error(err)
	return err
}
