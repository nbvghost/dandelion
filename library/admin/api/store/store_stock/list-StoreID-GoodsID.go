package store_stock

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"strings"

	"github.com/nbvghost/dandelion/entity/model"

	"github.com/nbvghost/dandelion/library/db"
)

type ListStoreIDGoodsID struct {
	POST struct {
		Datatables *model.Datatables `body:""`
		GoodsID    uint              `uri:"GoodsID"`
		StoreID    uint              `uri:"StoreID"`
	} `method:"POST"`
}

func (m *ListStoreIDGoodsID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *ListStoreIDGoodsID) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)

	//GoodsID
	//GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	//StoreID, _ := strconv.ParseUint(context.PathParams["StoreID"], 10, 64)

	Orm := db.Orm()
	//SELECT g.ID as ID,g.Title as Title,COUNT(ss.ID) as Total,SUM(ss.Stock) as Stock FROM Goods as g,StoreStock as ss where ss.StoreID=2009 and g.ID=ss.GoodsID  group by ss.GoodsID;
	type Result struct {
		*model.StoreStock    `json:"StoreStock"`
		*model.Store         `json:"Store"`
		*model.Goods         `json:"Goods"`
		*model.Specification `json:"Specification"`
	}

	var results []Result

	var recordsTotal int64
	db := Orm.Table("StoreStock").Select("*").Joins("JOIN Goods on Goods.ID = StoreStock.GoodsID").Joins("JOIN Specification on Specification.ID = StoreStock.SpecificationID").Joins("JOIN Store on Store.ID = StoreStock.StoreID").Where("StoreStock.StoreID=?", m.POST.StoreID).Where("StoreStock.GoodsID=?", m.POST.GoodsID)

	for _, value := range m.POST.Datatables.Order {
		if !strings.EqualFold(m.POST.Datatables.Columns[value.Column].Data, "") {
			db = db.Order(m.POST.Datatables.Columns[value.Column].Data + " " + value.Dir)
		}
	}

	db.Limit(10).Offset(0).Find(&results)
	db.Offset(0).Count(&recordsTotal)
	//fmt.Println(result)
	//fmt.Println(recordsTotal)

	return &result.JsonResult{Data: map[string]interface{}{"data": results, "draw": m.POST.Datatables.Draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsTotal}}, nil
}
