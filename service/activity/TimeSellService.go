package activity

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gpa/types"
)

type TimeSellService struct {
	model.BaseDao
}

func (service TimeSellService) GetTimeSellByHash(Hash string, OID types.PrimaryKey) model.TimeSell {
	var timesell model.TimeSell
	singleton.Orm().Model(&model.TimeSell{}).Where("Hash=? and OID=?", Hash, OID).First(&timesell)
	return timesell
}

//todo:
func (service TimeSellService) GetTimeSellByGoodsID(GoodsID types.PrimaryKey, OID types.PrimaryKey) *model.TimeSell {
	//timesellGoods := service.GetTimeSellGoodsByGoodsID(GoodsID, OID)
	var timesellGoods model.TimeSellGoods
	singleton.Orm().Model(&model.TimeSellGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)

	var timesell model.TimeSell
	singleton.Orm().Model(&model.TimeSell{}).Where("Hash=? and OID=?", timesellGoods.TimeSellHash, timesellGoods.OID).First(&timesell)
	return &timesell
}
func (service TimeSellService) GetTimeSellGoodsByGoodsID(GoodsID types.PrimaryKey, OID types.PrimaryKey) model.TimeSellGoods {
	var timesellGoods model.TimeSellGoods
	singleton.Orm().Model(&model.TimeSellGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)
	return timesellGoods
}

/*
func (service TimeSellService) AddTimeSellAction(context *gweb.Context) (r gweb.Result,err error) {
	//:Hash/:GoodsID
	context.Request.ParseForm()
	//Hash := context.Request.FormValue("Hash")
	//GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "", nil)}
}*/

func (service TimeSellService) DeleteTimeSell(TimeSellID types.PrimaryKey) error {
	//timesell := TimeSellService{}.GetTimeSellByGoodsID(GoodsID)
	var ts model.TimeSell
	service.Get(singleton.Orm(), TimeSellID, &ts)
	//err := service.Delete(singleton.Orm(), &model.TimeSell{}, ts.ID)
	err := service.DeleteWhere(singleton.Orm(), &model.TimeSell{}, map[string]interface{}{
		"Hash": ts.Hash,
	})
	glog.Error(err)
	return err
}
