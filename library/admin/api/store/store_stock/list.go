package store_stock

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"strings"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
)

type List struct {
	POST struct {
		Datatables *model.Datatables `body:""`
	} `method:"POST"`
}

func (m *List) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)

	Orm := db.Orm()
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		GoodsID  uint   `gorm:"column:GoodsID"`
		StoreID  uint   `gorm:"column:StoreID"`
		Title    string `gorm:"column:Title"`
		Total    uint   `gorm:"column:Total"`
		Stock    uint   `gorm:"column:Stock"`
		UseStock uint   `gorm:"column:UseStock"`
	}

	var results []Result

	var recordsTotal int64
	db := Orm.Table("StoreStock").Select("Goods.ID as GoodsID,StoreStock.StoreID as StoreID,Goods.Title as Title,COUNT(StoreStock.ID) as Total,SUM(StoreStock.Stock) as Stock,SUM(StoreStock.UseStock) as UseStock").Joins("JOIN Goods on Goods.ID = StoreStock.GoodsID").Where("StoreID=?", m.POST.Datatables.Columns[1].Search.Value).Group("StoreStock.GoodsID")

	for _, value := range m.POST.Datatables.Order {
		if !strings.EqualFold(m.POST.Datatables.Columns[value.Column].Data, "") {
			db = db.Order(m.POST.Datatables.Columns[value.Column].Data + " " + value.Dir)
		}
	}

	db.Limit(m.POST.Datatables.Length).Offset(m.POST.Datatables.Start).Find(&results)
	db.Offset(0).Count(&recordsTotal)
	//fmt.Println(result)
	//fmt.Println(recordsTotal)
	return &result.JsonResult{Data: map[string]interface{}{"data": results, "draw": m.POST.Datatables.Draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsTotal}}, nil

	/*dts := &model.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]model.StoreStock{})
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}*/
}

func (m *List) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}
