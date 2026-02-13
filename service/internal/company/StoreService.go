package company

import (
	"context"
	"log"

	"github.com/nbvghost/dandelion/library/db"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type StoreService struct {
	model.BaseDao
}

func (service StoreService) GetByPhone(ctx context.Context, Phone string) model.Store {
	Orm := db.GetDB(ctx)
	var store model.Store
	Orm.Model(&model.Store{}).Where(&model.Store{Phone: Phone}).First(&store)
	return store
}
func (service StoreService) LocationList(ctx context.Context, Latitude, Longitude float64) []map[string]interface{} {
	Orm := db.GetDB(ctx)

	rows, err := Orm.Model(&model.Store{}).Select("ID,Images,Name,Address,ServicePhone,Stars,StarsCount,ROUND(6378.138*2*ASIN(SQRT(POW(SIN((?*PI()/180-Latitude*PI()/180)/2),2)+COS(?*PI()/180)*COS(Latitude*PI()/180)*POW(SIN((?*PI()/180-Longitude*PI()/180)/2),2)))*1000) AS Distance", Latitude, Latitude, Longitude).Order("Distance asc").Rows()
	log.Println(err)
	defer rows.Close()

	list := make([]map[string]interface{}, 0)
	for rows.Next() {

		var ID uint
		var Images string
		var Name string
		var Address string
		var ServicePhone string
		var Stars uint
		var StarsCount uint
		var Distance float64

		err = rows.Scan(&ID, &Images, &Name, &Address, &ServicePhone, &Stars, &StarsCount, &Distance)
		log.Println(err)
		list = append(list, map[string]interface{}{"ID": ID, "Images": Images, "Name": Name, "Address": Address, "ServicePhone": ServicePhone, "Stars": Stars, "StarsCount": StarsCount, "Distance": Distance})
	}

	return list
}

func (service StoreService) GetByGoodsIDAndSpecificationIDAndStoreID(ctx context.Context, GoodsID, SpecificationID, StoreID dao.PrimaryKey) *model.StoreStock {
	Orm := db.GetDB(ctx)
	var ss model.StoreStock
	Orm.Where(&model.StoreStock{GoodsID: GoodsID, SpecificationID: SpecificationID, StoreID: StoreID}).First(&ss)
	return &ss
}

func (service StoreService) ListStoreSpecifications(ctx context.Context, StoreID, GoodsID dao.PrimaryKey) interface{} {

	Orm := db.GetDB(ctx)
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		*model.StoreStock    `json:"StoreStock"`
		*model.Specification `json:"Specification"`
	}

	var result []Result

	db := Orm.Table("StoreStock").Select("*").Joins("JOIN Specification on Specification.ID = StoreStock.SpecificationID").Where("StoreStock.StoreID=?", StoreID).Where("StoreStock.GoodsID=?", GoodsID).Group("StoreStock.SpecificationID")

	db.Find(&result)
	//db.Limit(dts.Length).Offset(dts.Start).Find(&result)
	//db.Offset(0).Count(&recordsTotal)

	return result
}
func (service StoreService) ListStoreStock(ctx context.Context, StoreID dao.PrimaryKey) interface{} {

	Orm := db.GetDB(ctx)
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		*model.StoreStock `json:"StoreStock"`
		*model.Goods      `json:"Goods"`
		Total             uint `gorm:"column:Total"`
		TotalStock        uint `gorm:"column:TotalStock"`
		TotalUseStock     uint `gorm:"column:TotalUseStock"`
	}

	var result []Result

	db := Orm.Table("StoreStock").Select("*,COUNT(StoreStock.ID) as Total,SUM(StoreStock.Stock) as TotalStock,SUM(StoreStock.UseStock) as TotalUseStock").Joins("JOIN Goods on Goods.ID = StoreStock.GoodsID").Where("StoreStock.StoreID=?", StoreID).Group("StoreStock.GoodsID")

	db.Find(&result)
	//db.Limit(dts.Length).Offset(dts.Start).Find(&result)
	//db.Offset(0).Count(&recordsTotal)

	return result
}
