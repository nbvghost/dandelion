package service

import (
	"errors"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"
	"strconv"

	"strings"

	"github.com/nbvghost/dandelion/app/play"

	"fmt"

	"github.com/nbvghost/gweb"
)

type StoreStockService struct {
	dao.BaseDao
	Journal JournalService
}

func (service StoreStockService) VerificationSelf(StoreID, StoreStockID, Quantity uint64) *dao.ActionStatus {
	Orm := dao.Orm()
	var ss dao.StoreStock
	err := service.Get(Orm, StoreStockID, &ss)
	if err != nil {
		return (&dao.ActionStatus{}).SmartError(err, "", 0)
	}

	if ss.StoreID != StoreID {
		return (&dao.ActionStatus{}).SmartError(errors.New("找不到相关库存"), "", 0)
	}

	if Quantity == 0 {
		return (&dao.ActionStatus{}).SmartError(errors.New("无效的数量"), "", 0)
	}

	if Quantity > ss.Stock-ss.UseStock {
		//库存不足
		return (&dao.ActionStatus{}).SmartError(errors.New("库存不足"), "", 0)
	} else {

		var specification dao.Specification
		service.Get(Orm, ss.SpecificationID, &specification)

		var store dao.Store
		service.Get(Orm, StoreID, &store)

		var goods dao.Goods
		service.Get(Orm, ss.GoodsID, &goods)

		if ss.GoodsID != specification.GoodsID {
			return (&dao.ActionStatus{}).SmartError(errors.New("找不到相关库存"), "", 0)
		} else {
			//判断金额,

			if specification.CostPrice*Quantity > store.Amount {
				//金额不足
				return (&dao.ActionStatus{}).SmartError(errors.New("门店金额不足"), "", specification.CostPrice*Quantity-store.Amount)
			} else {

				tx := Orm.Begin()

				err := service.ChangeMap(tx, StoreStockID, &dao.StoreStock{}, map[string]interface{}{"UseStock": ss.UseStock + Quantity})
				if err != nil {
					tx.Rollback()
					return (&dao.ActionStatus{}).SmartError(err, "", 0)
				} else {
					detail := fmt.Sprintf("%v,规格：%v(%v)kg成本价：%v，数量：%v", goods.Title, specification.Label, float64(specification.Num)*float64(specification.Weight)/1000, specification.CostPrice, Quantity)
					err = service.Journal.AddStoreJournal(tx, StoreID, "自主核销商品库存", detail, play.StoreJournal_Type_ZZHX, -int64(specification.CostPrice*Quantity), ss.ID)
					if err != nil {
						tx.Rollback()
						return (&dao.ActionStatus{}).SmartError(err, "", 0)
					} else {
						tx.Commit()
						return (&dao.ActionStatus{}).SmartError(err, "自主核销成功", 0)
					}

				}

			}
		}

	}

}

func (service StoreStockService) GetByGoodsIDAndSpecificationIDAndStoreID(GoodsID, SpecificationID, StoreID uint64) *dao.StoreStock {
	Orm := dao.Orm()
	var ss dao.StoreStock
	Orm.Where(&dao.StoreStock{GoodsID: GoodsID, SpecificationID: SpecificationID, StoreID: StoreID}).First(&ss)
	return &ss
}

func (service StoreStockService) GetItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.StoreStock{}
	err := service.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (service StoreStockService) ListExistGoodsIDS(context *gweb.Context) gweb.Result {
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
func (service StoreStockService) ListByGoods(context *gweb.Context) gweb.Result {
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

func (service StoreStockService) ListStoreSpecifications(StoreID, GoodsID uint64) interface{} {

	Orm := dao.Orm()
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		*dao.StoreStock    `json:"StoreStock"`
		*dao.Specification `json:"Specification"`
	}

	var result []Result

	db := Orm.Table("StoreStock").Select("*").Joins("JOIN Specification on Specification.ID = StoreStock.SpecificationID").Where("StoreStock.StoreID=?", StoreID).Where("StoreStock.GoodsID=?", GoodsID).Group("StoreStock.SpecificationID")

	db.Find(&result)
	//db.Limit(dts.Length).Offset(dts.Start).Find(&result)
	//db.Offset(0).Count(&recordsTotal)

	return result
}
func (service StoreStockService) ListStoreStock(StoreID uint64) interface{} {

	Orm := dao.Orm()
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		*dao.StoreStock `json:"StoreStock"`
		*dao.Goods      `json:"Goods"`
		Total           uint64 `gorm:"column:Total"`
		TotalStock      uint64 `gorm:"column:TotalStock"`
		TotalUseStock   uint64 `gorm:"column:TotalUseStock"`
	}

	var result []Result

	db := Orm.Table("StoreStock").Select("*,COUNT(StoreStock.ID) as Total,SUM(StoreStock.Stock) as TotalStock,SUM(StoreStock.UseStock) as TotalUseStock").Joins("JOIN Goods on Goods.ID = StoreStock.GoodsID").Where("StoreStock.StoreID=?", StoreID).Group("StoreStock.GoodsID")

	db.Find(&result)
	//db.Limit(dts.Length).Offset(dts.Start).Find(&result)
	//db.Offset(0).Count(&recordsTotal)

	return result
}
func (service StoreStockService) ListItem(context *gweb.Context) gweb.Result {
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
func (service StoreStockService) DeleteItem(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.StoreStock{}
	err := service.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}

//库存的修改与添加
func (service StoreStockService) SaveItem(context *gweb.Context) gweb.Result {
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
	service.Get(Orm, ID, &item)
	if item.ID == 0 {

		if AddStoreStockStock < 0 {
			//return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("增加的库存不能小于0"), "", nil)}
			AddStoreStockStock = 0
		}

		item.StoreID = StoreID
		item.GoodsID = GoodsID
		item.SpecificationID = SpecificationID
		item.Stock = uint64(AddStoreStockStock)
		err := service.Add(Orm, item)
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

		err := service.ChangeMap(Orm, ID, item, map[string]interface{}{"SpecificationID": item.SpecificationID, "Stock": item.Stock})
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
	}

}
