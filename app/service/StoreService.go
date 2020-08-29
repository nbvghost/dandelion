package service

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/glog"
)

type StoreService struct {
	dao.BaseDao
	Journal JournalService
}

func (service StoreService) GetByPhone(Phone string) dao.Store {
	Orm := dao.Orm()
	var store dao.Store
	Orm.Model(&dao.Store{}).Where(&dao.Store{Phone: Phone}).First(&store)
	return store
}
func (service StoreService) LocationList(Latitude, Longitude float64) []map[string]interface{} {
	Orm := dao.Orm()

	rows, err := Orm.Model(&dao.Store{}).Select("ID,Images,Name,Address,ServicePhone,Stars,StarsCount,ROUND(6378.138*2*ASIN(SQRT(POW(SIN((?*PI()/180-Latitude*PI()/180)/2),2)+COS(?*PI()/180)*COS(Latitude*PI()/180)*POW(SIN((?*PI()/180-Longitude*PI()/180)/2),2)))*1000) AS Distance", Latitude, Latitude, Longitude).Order("Distance asc").Rows()
	glog.Error(err)
	defer rows.Close()

	list := make([]map[string]interface{}, 0)
	for rows.Next() {

		var ID uint64
		var Images string
		var Name string
		var Address string
		var ServicePhone string
		var Stars uint64
		var StarsCount uint64
		var Distance float64

		err = rows.Scan(&ID, &Images, &Name, &Address, &ServicePhone, &Stars, &StarsCount, &Distance)
		glog.Error(err)
		list = append(list, map[string]interface{}{"ID": ID, "Images": Images, "Name": Name, "Address": Address, "ServicePhone": ServicePhone, "Stars": Stars, "StarsCount": StarsCount, "Distance": Distance})
	}

	return list
}

func (service StoreService) VerificationSelf(StoreID, StoreStockID, Quantity uint64) *dao.ActionStatus {
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

func (service StoreService) GetByGoodsIDAndSpecificationIDAndStoreID(GoodsID, SpecificationID, StoreID uint64) *dao.StoreStock {
	Orm := dao.Orm()
	var ss dao.StoreStock
	Orm.Where(&dao.StoreStock{GoodsID: GoodsID, SpecificationID: SpecificationID, StoreID: StoreID}).First(&ss)
	return &ss
}

func (service StoreService) ListStoreSpecifications(StoreID, GoodsID uint64) interface{} {

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
func (service StoreService) ListStoreStock(StoreID uint64) interface{} {

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
