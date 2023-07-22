package goods

import (
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"gorm.io/gorm"
)

type SKUService struct {
}

func (service SKUService) SkuLabelByGoodsID(orm *gorm.DB, goodsID dao.PrimaryKey) []extends.SkuLabel {
	goodsSkuLabelDataMap := make(map[dao.PrimaryKey][]*model.GoodsSkuLabelData, 0)
	{
		goodsSkuLabelData := dao.Find(orm, &model.GoodsSkuLabelData{}).Where(`"GoodsID"=?`, goodsID).List()
		for i := range goodsSkuLabelData {
			item := goodsSkuLabelData[i].(*model.GoodsSkuLabelData)
			if _, ok := goodsSkuLabelDataMap[item.GoodsSkuLabelID]; !ok {
				goodsSkuLabelDataMap[item.GoodsSkuLabelID] = make([]*model.GoodsSkuLabelData, 0)
			}
			goodsSkuLabelDataMap[item.GoodsSkuLabelID] = append(goodsSkuLabelDataMap[item.GoodsSkuLabelID], item)
		}
	}

	goodsSkuLabel := dao.Find(orm, &model.GoodsSkuLabel{}).Where(`"GoodsID"=? and "Abel"=?`, goodsID, true).Order(`"Image" desc`).List()
	skuLabelList := make([]extends.SkuLabel, len(goodsSkuLabel))
	for i := range goodsSkuLabel {
		item := goodsSkuLabel[i].(*model.GoodsSkuLabel)
		skuLabelList[i].Label = item
		skuLabelList[i].Data = goodsSkuLabelDataMap[item.ID]
	}
	return skuLabelList
}
func (service SKUService) SkuLabelByOID(orm *gorm.DB, oid dao.PrimaryKey) []extends.SkuLabel {
	goodsSkuLabelDataMap := make(map[dao.PrimaryKey][]*model.GoodsSkuLabelData, 0)
	{
		goodsSkuLabelData := dao.Find(orm, &model.GoodsSkuLabelData{}).Where(`"OID"=?`, oid).List()
		for i := range goodsSkuLabelData {
			item := goodsSkuLabelData[i].(*model.GoodsSkuLabelData)
			if _, ok := goodsSkuLabelDataMap[item.GoodsSkuLabelID]; !ok {
				goodsSkuLabelDataMap[item.GoodsSkuLabelID] = make([]*model.GoodsSkuLabelData, 0)
			}
			goodsSkuLabelDataMap[item.GoodsSkuLabelID] = append(goodsSkuLabelDataMap[item.GoodsSkuLabelID], item)
		}
	}

	goodsSkuLabel := dao.Find(orm, &model.GoodsSkuLabel{}).Where(`"OID"=? and "Abel"=?`, oid, true).Order(`"Image" desc`).List()
	skuLabelList := make([]extends.SkuLabel, len(goodsSkuLabel))
	for i := range goodsSkuLabel {
		item := goodsSkuLabel[i].(*model.GoodsSkuLabel)
		skuLabelList[i].Label = item
		skuLabelList[i].Data = goodsSkuLabelDataMap[item.ID]
	}
	return skuLabelList
}
