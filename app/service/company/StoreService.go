package company

import (
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/glog"
)

type StoreService struct {
	dao.BaseDao
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
