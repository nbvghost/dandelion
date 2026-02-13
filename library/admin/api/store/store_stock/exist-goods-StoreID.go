package store_stock

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ExistGoodsStoreID struct {
	GET struct {
		StoreID uint `uri:"StoreID"`
	} `method:"GET"`
}

func (m *ExistGoodsStoreID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *ExistGoodsStoreID) HandleGet(ctx constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.GetDB(ctx)
	//StoreID, _ := strconv.ParseUint(context.PathParams["StoreID"], 10, 64)
	//StoreID := object.ParseUint(context.PathParams["StoreID"])
	type Result struct {
		GoodsID uint `gorm:"column:GoodsID"`
	}

	var ids []Result
	Orm.Table("StoreStock").Select("GoodsID as GoodsID").Where(&model.StoreStock{StoreID: dao.PrimaryKey(m.GET.StoreID)}).Find(&ids)
	//fmt.Println(ids)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", ids)}, nil
}
