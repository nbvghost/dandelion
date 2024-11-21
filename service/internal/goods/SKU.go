package goods

import (
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"gorm.io/gorm"
)

type SKUService struct {
}

/*
	func (service SKUService) SkuLabelByGoodsID(orm *gorm.DB, goodsID dao.PrimaryKey,specifications []*model.Specification) []extends.SkuLabel {
		skuLabelMap := make(map[dao.PrimaryKey]*model.GoodsSkuLabel)//make([]extends.SkuLabel, len(goodsSkuLabel))
		{
			skuLabel := dao.Find(orm, &model.GoodsSkuLabel{}).Where(`"GoodsID"=? and "Abel"=?`, goodsID, true).Order(`"Image" desc`).List()
			for i := range skuLabel {
				item := skuLabel[i].(*model.GoodsSkuLabel)
				//skuLabelList[i].Label = item
				//skuLabelList[i].Data = goodsSkuLabelDataMap[item.ID]
				skuLabelMap[item.ID]=item
			}
		}


		skuLabelDataMap := make(map[dao.PrimaryKey]*model.GoodsSkuLabelData)
		{
			goodsSkuLabelData := dao.Find(orm, &model.GoodsSkuLabelData{}).Where(`"GoodsID"=?`, goodsID).List()
			for i := range goodsSkuLabelData {
				item := goodsSkuLabelData[i].(*model.GoodsSkuLabelData)
				skuLabelDataMap[item.GoodsSkuLabelID] = item
			}
		}

		var skuNodes []extends.SkuNode


		var maxLabelIndex =0

		for i := 0; i < len(specifications); i++ {
			specification:=specifications[i]
			skuLabelDataID:=specification.LabelIndex[0]

			//计算，规格中最大的属性数量
			if len(specification.LabelIndex)>maxLabelIndex{
				maxLabelIndex=len(specification.LabelIndex)
			}

			skuLabelData:=skuLabelDataMap[skuLabelDataID]
			skuLabel:=skuLabelMap[skuLabelData.GoodsSkuLabelID]

			var has bool

			for ii := 0; ii < len(skuNodes); ii++ {
				if skuNodes[ii].Data.ID==skuLabelDataID {
					has=true
					break
				}
			}

			if !has{
				skuNodes=append(skuNodes,extends.SkuNode{
					Label:  skuLabel,
					Data:   skuLabelData,
					Childs: make([]extends.SkuNode, 0),
				})
			}
		}

		{

			for i := 0; i < maxLabelIndex; i++ {
				if i==0{
					continue
				}

				specification:=specifications[i]
				skuLabelDataID:=specification.LabelIndex[i]

				skuLabelData:=skuLabelDataMap[skuLabelDataID]
				skuLabel:=skuLabelMap[skuLabelData.GoodsSkuLabelID]



			}

		}




		return skuLabelList
	}
*/
func (service SKUService) SkuLabelByGoodsID(orm *gorm.DB, goodsID dao.PrimaryKey) []extends.SkuLabel {
	goodsSkuLabelDataMap := make(map[dao.PrimaryKey][]*model.GoodsSkuLabelData)
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

/*
	func (service SKUService) SkuLabelByOID(orm *gorm.DB, oid dao.PrimaryKey) []extends.SkuLabel {
		goodsSkuLabelData := dao.Find(orm, &model.GoodsSkuLabelData{}).Where(`"OID"=?`, oid).List()
		goodsSkuLabel := dao.Find(orm, &model.GoodsSkuLabel{}).Where(`"OID"=? and "Abel"=?`, oid, true).Order(`"Image" desc`).List()

		goodsSkuLabelDataMap := make(map[dao.PrimaryKey][]*model.GoodsSkuLabelData)
		{

			for i := range goodsSkuLabelData {
				item := goodsSkuLabelData[i].(*model.GoodsSkuLabelData)
				if _, ok := goodsSkuLabelDataMap[item.GoodsSkuLabelID]; !ok {
					goodsSkuLabelDataMap[item.GoodsSkuLabelID] = make([]*model.GoodsSkuLabelData, 0)
				}
				goodsSkuLabelDataMap[item.GoodsSkuLabelID] = append(goodsSkuLabelDataMap[item.GoodsSkuLabelID], item)
			}
		}


		skuLabelList := make([]extends.SkuLabel, len(goodsSkuLabel))
		for i := range goodsSkuLabel {
			item := goodsSkuLabel[i].(*model.GoodsSkuLabel)
			skuLabelList[i].Label = item
			skuLabelList[i].Data = goodsSkuLabelDataMap[item.ID]
		}
		return skuLabelList
	}
*/
func (service SKUService) SkuLabel(goodsSkuLabelList []model.GoodsSkuLabel, goodsSkuLabelDataList []model.GoodsSkuLabelData) []extends.SkuLabel {
	goodsSkuLabelDataMap := make(map[dao.PrimaryKey][]*model.GoodsSkuLabelData)
	for i := range goodsSkuLabelDataList {
		item := goodsSkuLabelDataList[i]
		if _, ok := goodsSkuLabelDataMap[item.GoodsSkuLabelID]; !ok {
			goodsSkuLabelDataMap[item.GoodsSkuLabelID] = make([]*model.GoodsSkuLabelData, 0)
		}
		goodsSkuLabelDataMap[item.GoodsSkuLabelID] = append(goodsSkuLabelDataMap[item.GoodsSkuLabelID], &item)
	}

	skuLabelList := make([]extends.SkuLabel, len(goodsSkuLabelList))
	for i := range goodsSkuLabelList {
		item := goodsSkuLabelList[i]
		skuLabelList[i].Label = &item
		skuLabelList[i].Data = goodsSkuLabelDataMap[item.ID]
	}
	return skuLabelList
}
