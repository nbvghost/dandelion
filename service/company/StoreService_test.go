package company

import (
	"fmt"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"testing"
)

func TestStoreStockService_ListStoreSpecifications(t *testing.T) {

	ss := StoreService{}
	fmt.Println(ss.ListStoreSpecifications(2009, 2005))
}
func TestStoreStockService_ListStoreStock(t *testing.T) {

	ss := StoreService{}
	fmt.Println(ss.ListStoreStock(2009))

}
func TestStoreStockService_ListByGoodsIDAndStoreID(t *testing.T) {
	Orm := singleton.Orm()

	type Result struct {
		GoodsID uint `gorm:"column:GoodsID"`
	}

	var ids []Result
	Orm.Table("StoreStock").Select("GoodsID as GoodsID").Where(&model.StoreStock{StoreID: 2009}).Find(&ids)
	fmt.Println(ids)

	/*Orm :=singleton.Orm()
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		*model.StoreStock    `json:"StoreStock"`
		*model.Store         `json:"Store"`
		*model.Goods         `json:"Goods"`
		*model.Specification `json:"Specification"`
	}

	var result []Result

	var recordsTotal uint
	db := Orm.Table("StoreStock").Select("*").Joins("JOIN Goods as Goods on Goods.ID = StoreStock.GoodsID").Joins("JOIN Specification on Specification.ID = StoreStock.SpecificationID").Joins("JOIN Store on Store.ID = StoreStock.StoreID")
	db.Where("StoreID=?", 2009).Where("GoodsID=?", 2004)
	db.Limit(10).Offset(0).Find(&result)
	db.Offset(0).Count(&recordsTotal)
	fmt.Println(result)
	fmt.Println(recordsTotal)

	b, _ := json.Marshal(result)
	fmt.Println(string(b))*/

}

/*func TestStoreStockService_ListItem(t *testing.T) {
	Orm :=singleton.Orm()
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		ID      uint `gorm:"column:ID"`
		StoreID uint `gorm:"column:StoreID"`
		Title   string `gorm:"column:Title"`
		Total   uint `gorm:"column:Total"`
		Stock   uint `gorm:"column:Stock"`
	}

	var result []Result

	var recordsTotal uint
	db := Orm.Table("StoreStock").Select("Goods.ID as ID,StoreStock.StoreID as StoreID,Goods.Title as Title,COUNT(StoreStock.ID) as Total,SUM(StoreStock.Stock) as Stock").Joins("JOIN Goods on Goods.ID = StoreStock.GoodsID").Where("StoreID=?", 2009).Group("StoreStock.GoodsID")
	db.Limit(10).Offset(0).Find(&result)
	db.Offset(0).Count(&recordsTotal)
	fmt.Println(result)
	fmt.Println(recordsTotal)

}
*/
