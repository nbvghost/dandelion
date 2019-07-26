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

	//"GoodsTypeID":GoodsTypeID,
	//	"GoodsTypeChildID":GoodsTypeChildID,

	//gtid={{$v.GoodsType.ID}}&gtcid
	GoodsTypeID,_:=strconv.ParseUint(context.PathParams["GoodsTypeID"],10,64)
	GoodsTypeChildID,_:=strconv.ParseUint(context.PathParams["GoodsTypeChildID"],10,64)

	gs:=controller.Goods.ListGoodsChildByGoodsTypeID(GoodsTypeID,GoodsTypeChildID)

	return &gweb.HTMLResult{Params: map[string]interface{}{
		"Goods":gs,
	}}
}
func (controller *Controller) categoriesPage(context *gweb.Context) gweb.Result {
	org:=controller.Organization.FindByDomain(dao.Orm(),context.PathParams["siteName"])
	gts := controller.Goods.ListGoodsType(org.ID)
	gtcs:=controller.Goods.ListGoodsTypeChildByOID(org.ID)

	//GoodsTypeID,GoodsTypeChildID




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
	GoodsTypeID,_:=strconv.ParseUint(context.Request.URL.Query().Get("GoodsTypeID"),10,64)
	GoodsTypeChildID,_:=strconv.ParseUint(context.Request.URL.Query().Get("GoodsTypeChildID"),10,64)

	if GoodsTypeID<=0 && len(GoodsTypeMap)>0{
		item:=GoodsTypeMap[0].(map[string]interface{})
		GoodsTypeID=item["GoodsType"].(dao.GoodsType).ID

		if item["GoodsTypeChild"]!=nil{
			itemChild:=item["GoodsTypeChild"].([]dao.GoodsTypeChild)
			if len(itemChild)>0{
				GoodsTypeChildID=itemChild[0].ID
			}

		}


	}


	hotGoods:=controller.Goods.HotListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID,GoodsTypeChildID,3)
	newGoods:=controller.Goods.NewListByGoodsTypeIDAndGoodsTypeChildID(GoodsTypeID,GoodsTypeChildID,3)

	return &gweb.HTMLResult{Params: map[string]interface{}{
		"GoodsType":GoodsTypeMap,
		"NewGoods":newGoods,
		"HotGoods":hotGoods,
		"GoodsTypeID":GoodsTypeID,
		"GoodsTypeChildID":GoodsTypeChildID,
	}}
}
func (controller *Controller) AddProjectdsfdsfsdAction(context *gweb.Context) gweb.Result {



	return &gweb.FileServerResult{}
}
func (controller *Controller) articleAction(context *gweb.Context) gweb.Result {

	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK"}}
}