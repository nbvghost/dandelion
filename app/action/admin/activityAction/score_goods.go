package activityAction

import (
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/activity"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/journal"
	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/gweb"
	"strconv"
)

type ScoreGoodsController struct {
	gweb.BaseController
	ScoreGoods activity.ScoreGoodsService
	CardItem   activity.CardItemService
	Journal    journal.JournalService
}

func (controller *ScoreGoodsController) Init() {

	controller.AddHandler(gweb.POSMethod("score_goods", controller.AddScoreGoods))
	controller.AddHandler(gweb.GETMethod("score_goods/{ID}", controller.GetScoreGoods))
	controller.AddHandler(gweb.POSMethod("score_goods/list", controller.DatatablesScoreGoods))
	controller.AddHandler(gweb.DELMethod("score_goods/{ID}", controller.DeleteScoreGoods))
	controller.AddHandler(gweb.PUTMethod("score_goods/{ID}", controller.ChangeScoreGoods))

}
func (controller *ScoreGoodsController) DeleteScoreGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.ScoreGoods{}
	err := controller.ScoreGoods.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}
}
func (controller *ScoreGoodsController) ChangeScoreGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.ScoreGoods{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	}
	err = controller.ScoreGoods.ChangeModel(Orm, ID, item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}
}
func (controller *ScoreGoodsController) AddScoreGoods(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	item := &dao.ScoreGoods{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	}
	item.OID = company.ID
	err = controller.ScoreGoods.Add(Orm, item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}
}
func (controller *ScoreGoodsController) GetScoreGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.ScoreGoods{}
	err := controller.ScoreGoods.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", item)}
}
func (controller *ScoreGoodsController) DatatablesScoreGoods(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.ScoreGoods.DatatablesListOrder(Orm, dts, &[]dao.ScoreGoods{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
