package activity

import (
	"context"
	"log"

	"github.com/nbvghost/dandelion/library/db"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type TimeSellService struct {
	model.BaseDao
}

func (service TimeSellService) GetTimeSellByHash(ctx context.Context, Hash string, OID dao.PrimaryKey) *model.TimeSell {
	var timesell model.TimeSell
	db.GetDB(ctx).Model(&model.TimeSell{}).Where("Hash=? and OID=?", Hash, OID).First(&timesell)
	return &timesell
}

func (service TimeSellService) GetTimeSellByGoodsID(ctx context.Context, GoodsID dao.PrimaryKey, OID dao.PrimaryKey) *model.TimeSell {
	//todo:考虑合并成一条sql语句
	//timesellGoods := service.GetTimeSellGoodsByGoodsID(GoodsID, OID)
	var timesellGoods model.TimeSellGoods
	db.GetDB(ctx).Model(&model.TimeSellGoods{}).Where(`"GoodsID"=? and "OID"=?`, GoodsID, OID).First(&timesellGoods)

	var timesell model.TimeSell
	db.GetDB(ctx).Model(&model.TimeSell{}).Where(`"Hash"=? and "OID"=?`, timesellGoods.TimeSellHash, timesellGoods.OID).First(&timesell)
	return &timesell
}
func (service TimeSellService) GetTimeSellGoodsByGoodsID(ctx context.Context, GoodsID dao.PrimaryKey, OID dao.PrimaryKey) model.TimeSellGoods {
	var timesellGoods model.TimeSellGoods
	db.GetDB(ctx).Model(&model.TimeSellGoods{}).Where(`"GoodsID"=? and "OID"=?`, GoodsID, OID).First(&timesellGoods)
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

func (service TimeSellService) DeleteTimeSell(ctx context.Context, TimeSellID dao.PrimaryKey) error {
	//timesell := TimeSellService{}.GetTimeSellByGoodsID(GoodsID)
	//var ts model.TimeSell
	ts := dao.GetByPrimaryKey(db.GetDB(ctx), &model.TimeSell{}, TimeSellID).(*model.TimeSell)
	//err := service.Delete(singleton.Orm(), &model.TimeSell{}, ts.ID)
	err := dao.DeleteBy(db.GetDB(ctx), &model.TimeSell{}, map[string]interface{}{
		"Hash": ts.Hash,
	})
	log.Println(err)
	return err
}
