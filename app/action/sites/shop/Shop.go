package shop

import (
	"strconv"

	"dandelion/app/service"
	"dandelion/app/service/dao"

	"github.com/nbvghost/gweb"
)



type Controller struct {
	gweb.BaseController
	Goods service.GoodsService
	Organization service.OrganizationService
}

func (controller *Controller) Apply() {
	controller.AddHandler(gweb.ALLMethod("/resources/", controller.AddProjectdsfdsfsdAction))
	controller.AddHandler(gweb.GETMethod("article/{ArticleID}", controller.articleAction))
	controller.AddHandler(gweb.GETMethod("categories", controller.categoriesPage))
	controller.AddHandler(gweb.GETMethod("products", controller.productsPage))
	controller.AddHandler(gweb.GETMethod("product", controller.productPage))
	controller.AddHandler(gweb.GETMethod("index", controller.indexPage))
}
func (controller *Controller) indexPage(context *gweb.Context) gweb.Result {


	gs:=controller.Goods.HotList(8)
	gtcs:=controller.Goods.GetTopGoodsTypeChild(dao.Orm(),3)

	return &gweb.HTMLResult{Params: map[string]interface{}{
		"Goods":gs,
		"TopGoodsTypeChildList":gtcs,
	}}

}
func (controller *Controller) productPage(context *gweb.Context) gweb.Result {
	ID,_:=strconv.ParseUint(context.Request.URL.Query().Get("ID"),10,64)

	GoodsInfo:=controller.Goods.GetGoods(dao.Orm(),ID)

	return &gweb.HTMLResult{Params: map[string]interface{}{
		"GoodsInfo":GoodsInfo,
	}}
}
func (controller *Controller) productsPage(context *gweb.Context) gweb.Result {

	//gtid={{$v.GoodsType.ID}}&gtcid
	gtid,_:=strconv.ParseUint(context.PathParams["gtid"],10,64)
	gtcid,_:=strconv.ParseUint(context.PathParams["gtcid"],10,64)

	gs:=controller.Goods.ListGoodsChildByGoodsTypeID(gtid,gtcid)

	return &gweb.HTMLResult{Params: map[string]interface{}{
		"Goods":gs,
	}}
}
func (controller *Controller) categoriesPage(context *gweb.Context) gweb.Result {
	org:=controller.Organization.FindByDomain(dao.Orm(),context.PathParams["siteName"])
	gts := controller.Goods.ListGoodsType(org.ID)
	gtcs:=controller.Goods.ListGoodsTypeChildByOID(org.ID)

	GoodsTypeMap:=make([]interface{},0)

	for index:=range gts{
		item:=gts[index]

		childList:=make([]dao.GoodsTypeChild,0)
		for cindex:=range gtcs{
			citem:=gtcs[cindex]
			if citem.GoodsTypeID==item.ID{
				childList = append(childList,citem)
			}
		}

		GoodsTypeMap=append(GoodsTypeMap,map[string]interface{}{
			"GoodsType":item,
			"GoodsTypeChild":childList,
		})
	}



	return &gweb.HTMLResult{Params: map[string]interface{}{
		"GoodsType":GoodsTypeMap,
	}}
}
func (controller *Controller) AddProjectdsfdsfsdAction(context *gweb.Context) gweb.Result {



	return &gweb.FileServerResult{}
}
func (controller *Controller) articleAction(context *gweb.Context) gweb.Result {

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK"}}
}