package goods

import (
	"context"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type SortService struct {
}

func (service SortService) HotSaleList(ctx context.Context, OID dao.PrimaryKey, count uint) []model.Goods {
	Orm := db.GetDB(ctx)
	var result []model.Goods
	db := Orm.Model(&model.Goods{}).Where(map[string]interface{}{"OID": OID}).Order(`"CountSale" desc`).Limit(int(count))
	db.Find(&result)
	return result
}
func (service SortService) HotViewList(ctx context.Context, OID dao.PrimaryKey, count uint) []model.Goods {
	Orm := db.GetDB(ctx)
	var result []model.Goods
	db := Orm.Model(&model.Goods{}).Where(map[string]interface{}{"OID": OID}).Order(`"CountView" desc`).Limit(int(count))
	db.Find(&result)
	return result
}
func (service SortService) RecommendGoods(ctx context.Context, OID, GoodsID, GoodsTypeID, GoodsTypeChildID dao.PrimaryKey) []model.Goods {
	var contentList []model.Goods
	db.GetDB(ctx).
		Model(&model.Goods{}).
		Where(`"OID"=? and "ID"<>? and ("GoodsTypeID"=? or "GoodsTypeChildID"=?)`, OID, GoodsID, GoodsTypeID, GoodsTypeChildID).
		Order("RAND()").
		Limit(6).
		Find(&contentList)
	return contentList
}
