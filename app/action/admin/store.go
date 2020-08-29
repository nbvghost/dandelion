package admin

import (
	"errors"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/gweb"
	"strconv"
	"strings"
)

type StoreController struct {
	gweb.BaseController
	Store service.StoreService
}

func (controller *StoreController) Init() {
	controller.AddHandler(gweb.POSMethod("add", controller.AddStoreItem))
	controller.AddHandler(gweb.GETMethod("{ID}", controller.GetStoreItem))
	controller.AddHandler(gweb.POSMethod("list", controller.ListStoreItem))
	controller.AddHandler(gweb.DELMethod("{ID}", controller.DeleteStoreItem))
	controller.AddHandler(gweb.PUTMethod("{ID}", controller.ChangeStoreItem))
	controller.AddHandler(gweb.POSMethod("stock", controller.SaveStoreStockItem))
	controller.AddHandler(gweb.PUTMethod("stock", controller.SaveStoreStockItem))
	controller.AddHandler(gweb.GETMethod("stock/{ID}", controller.GetStoreStockItem))
	controller.AddHandler(gweb.GETMethod("stock/exist/goods/{StoreID}", controller.ListExistGoodsIDS))
	controller.AddHandler(gweb.POSMethod("stock/list/{StoreID}/{GoodsID}", controller.ListByGoods))
	controller.AddHandler(gweb.POSMethod("stock/list", controller.ListStoreStockItem))
	controller.AddHandler(gweb.DELMethod("stock/{ID}", controller.DeleteStoreStockItem))
}
func (controller *StoreController) DeleteStoreStockItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.StoreStock{}
	err := controller.Store.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (controller *StoreController) ListStoreStockItem(context *gweb.Context) gweb.Result {
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	Orm := dao.Orm()
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		GoodsID  uint64 `gorm:"column:GoodsID"`
		StoreID  uint64 `gorm:"column:StoreID"`
		Title    string `gorm:"column:Title"`
		Total    uint64 `gorm:"column:Total"`
		Stock    uint64 `gorm:"column:Stock"`
		UseStock uint64 `gorm:"column:UseStock"`
	}

	var result []Result

	var recordsTotal uint64
	db := Orm.Table("StoreStock").Select("Goods.ID as GoodsID,StoreStock.StoreID as StoreID,Goods.Title as Title,COUNT(StoreStock.ID) as Total,SUM(StoreStock.Stock) as Stock,SUM(StoreStock.UseStock) as UseStock").Joins("JOIN Goods on Goods.ID = StoreStock.GoodsID").Where("StoreID=?", dts.Columns[1].Search.Value).Group("StoreStock.GoodsID")

	for _, value := range dts.Order {
		if !strings.EqualFold(dts.Columns[value.Column].Data, "") {
			db = db.Order(dts.Columns[value.Column].Data + " " + value.Dir)
		}
	}

	db.Limit(dts.Length).Offset(dts.Start).Find(&result)
	db.Offset(0).Count(&recordsTotal)
	//fmt.Println(result)
	//fmt.Println(recordsTotal)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": result, "draw": dts.Draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsTotal}}

	/*dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.StoreStock{})
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}*/
}
func (controller *StoreController) ListByGoods(context *gweb.Context) gweb.Result {
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	//GoodsID
	GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	StoreID, _ := strconv.ParseUint(context.PathParams["StoreID"], 10, 64)

	Orm := dao.Orm()
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		*dao.StoreStock    `json:"StoreStock"`
		*dao.Store         `json:"Store"`
		*dao.Goods         `json:"Goods"`
		*dao.Specification `json:"Specification"`
	}

	var result []Result

	var recordsTotal uint64
	db := Orm.Table("StoreStock").Select("*").Joins("JOIN Goods on Goods.ID = StoreStock.GoodsID").Joins("JOIN Specification on Specification.ID = StoreStock.SpecificationID").Joins("JOIN Store on Store.ID = StoreStock.StoreID").Where("StoreStock.StoreID=?", StoreID).Where("StoreStock.GoodsID=?", GoodsID)

	for _, value := range dts.Order {
		if !strings.EqualFold(dts.Columns[value.Column].Data, "") {
			db = db.Order(dts.Columns[value.Column].Data + " " + value.Dir)
		}
	}

	db.Limit(10).Offset(0).Find(&result)
	db.Offset(0).Count(&recordsTotal)
	//fmt.Println(result)
	//fmt.Println(recordsTotal)

	return &gweb.JsonResult{Data: map[string]interface{}{"data": result, "draw": dts.Draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsTotal}}

}

func (controller *StoreController) ListExistGoodsIDS(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	StoreID, _ := strconv.ParseUint(context.PathParams["StoreID"], 10, 64)

	type Result struct {
		GoodsID uint64 `gorm:"column:GoodsID"`
	}

	var ids []Result
	Orm.Table("StoreStock").Select("GoodsID as GoodsID").Where(&dao.StoreStock{StoreID: StoreID}).Find(&ids)
	//fmt.Println(ids)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", ids)}
}
func (controller *StoreController) GetStoreStockItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.StoreStock{}
	err := controller.Store.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}

//库存的修改与添加
func (controller *StoreController) SaveStoreStockItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	context.Request.ParseForm()

	//form.StoreID=parseInt($routeParams.ID);
	//form.GoodsID=$scope.SelectGoods.ID;
	//form.ID=$scope.StoreStock.ID;
	//form.SpecificationID=$scope.StoreStock.SpecificationID;
	//form.AddStoreStockStock=$scope.AddStoreStockStock;

	StoreID, _ := strconv.ParseUint(context.Request.FormValue("StoreID"), 10, 64)
	GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
	SpecificationID, _ := strconv.ParseUint(context.Request.FormValue("SpecificationID"), 10, 64)
	AddStoreStockStock, _ := strconv.ParseInt(context.Request.FormValue("AddStoreStockStock"), 10, 64)

	item := &dao.StoreStock{}
	controller.Store.Get(Orm, ID, &item)
	if item.ID == 0 {

		if AddStoreStockStock < 0 {
			//return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("增加的库存不能小于0"), "", nil)}
			AddStoreStockStock = 0
		}

		item.StoreID = StoreID
		item.GoodsID = GoodsID
		item.SpecificationID = SpecificationID
		item.Stock = uint64(AddStoreStockStock)
		err := controller.Store.Add(Orm, item)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
	} else {

		stock := int64(item.Stock) + AddStoreStockStock
		if stock < 0 {
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("库存不能为负的"), "", nil)}
		}
		if stock < int64(item.UseStock) {
			//return (&dao.ActionStatus{}).SmartError(errors.New("库存不能为负的"), "", 0)
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("库存不能为负的"), "", nil)}
		}

		item.Stock = uint64(stock)
		item.SpecificationID = SpecificationID

		err := controller.Store.ChangeMap(Orm, ID, item, map[string]interface{}{"SpecificationID": item.SpecificationID, "Stock": item.Stock})
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
	}

}
func (controller *StoreController) ChangeStoreItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Store{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	var _store dao.Store
	controller.Store.GetByPhone(item.Phone)
	if _store.ID > 0 && _store.ID != item.ID {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("手机号："+item.Phone+"已经被使用"), "", nil)}
	}

	err = controller.Store.ChangeModel(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (controller *StoreController) DeleteStoreItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Store{}
	err := controller.Store.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (controller *StoreController) ListStoreItem(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.Store.DatatablesListOrder(Orm, dts, &[]dao.Store{}, company.ID, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (controller *StoreController) GetStoreItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Store{}
	err := controller.Store.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (controller *StoreController) AddStoreItem(context *gweb.Context) gweb.Result {

	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	Orm := dao.Orm()
	item := &dao.Store{}
	item.OID = company.ID
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	var _store dao.Store
	_store = controller.Store.GetByPhone(item.Phone)
	if _store.ID > 0 {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("手机号："+item.Phone+"已经被使用"), "", nil)}
	}

	err = controller.Store.Add(Orm, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
