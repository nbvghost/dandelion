package goods

import (
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type SkuLabelData struct {
	//Goods        goods.GoodsService
	Organization *model.Organization `mapping:""`
	Get          struct {
		GoodsID         dao.PrimaryKey `form:"goods-id"`
		GoodsSkuLabelID dao.PrimaryKey `form:"goods-sku-label-id"`
		Name            string         `form:"name"`
	} `method:"Get"`
	Post struct {
		Label           string
		Name            string
		Image           string
		GoodsSkuLabelID dao.PrimaryKey
		GoodsID         dao.PrimaryKey
	} `method:"Post"`
	Put struct {
		ID    dao.PrimaryKey
		Label string
		Name  string
		Image string
	} `method:"Put"`
	Delete struct {
		ID dao.PrimaryKey `form:"goods-sku-label-data-id"`
	} `method:"Delete"`
}

func (m *SkuLabelData) listData(goodsID, goodsSkuLabelID dao.PrimaryKey, Name string, OID dao.PrimaryKey) map[dao.PrimaryKey][]*model.GoodsSkuLabelData {
	list := make(map[dao.PrimaryKey][]*model.GoodsSkuLabelData)
	where := make(map[string]any)
	where["OID"] = OID
	if goodsID > 0 {
		where["GoodsID"] = goodsID
	}
	if goodsSkuLabelID > 0 {
		where["GoodsSkuLabelID"] = goodsSkuLabelID
	}
	if len(Name) > 0 {
		where["Name"] = Name
	}

	LabelList := dao.Find(db.Orm(), &model.GoodsSkuLabelData{}).Where(where).List()
	for i := range LabelList {
		item := LabelList[i].(*model.GoodsSkuLabelData)
		if _, ok := list[item.GoodsSkuLabelID]; !ok {
			list[item.GoodsSkuLabelID] = make([]*model.GoodsSkuLabelData, 0)
		}
		list[item.GoodsSkuLabelID] = append(list[item.GoodsSkuLabelID], item)
	}
	return list
}
func (m *SkuLabelData) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	return result.NewData(map[string]any{"SkuLabelDataList": m.listData(m.Get.GoodsID, m.Get.GoodsSkuLabelID, m.Get.Name, m.Organization.ID)}), nil
}
func (m *SkuLabelData) HandleDelete(context constrain.IContext) (constrain.IResult, error) {
	has := dao.GetByPrimaryKey(db.Orm(), &model.GoodsSkuLabelData{}, m.Delete.ID).(*model.GoodsSkuLabelData)
	if has.IsZero() {
		return nil, errors.New("记录不存在")
	}
	tx := db.Orm().Begin()

	specificationList := dao.Find(tx, &model.Specification{}).Where(`"GoodsID"=?`, has.GoodsID).List()
	for i2 := range specificationList {
		specification := specificationList[i2].(*model.Specification)
		hasItem := false
		for i3 := range specification.LabelIndex {
			if specification.LabelIndex[i3] == has.ID {
				hasItem = true
				break
			}
		}
		if hasItem {
			err := dao.DeleteByPrimaryKey(tx, &model.Specification{}, specification.ID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	err := dao.DeleteByPrimaryKey(tx, &model.GoodsSkuLabelData{}, has.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return result.NewData(map[string]any{"SkuLabelDataList": m.listData(has.GoodsID, 0, "", m.Organization.ID)}), err
}
func (m *SkuLabelData) HandlePut(context constrain.IContext) (constrain.IResult, error) {
	tx := db.Orm().Begin()

	has := dao.GetByPrimaryKey(tx, &model.GoodsSkuLabelData{}, m.Put.ID).(*model.GoodsSkuLabelData)
	if has.IsZero() {
		tx.Rollback()
		return nil, errors.New("记录不存在")
	}

	hasLabel := dao.GetBy(tx, &model.GoodsSkuLabelData{}, map[string]any{"GoodsSkuLabelID": has.GoodsSkuLabelID, "Label": m.Put.Label}).(*model.GoodsSkuLabelData)
	if !hasLabel.IsZero() && hasLabel.ID != m.Put.ID {
		tx.Rollback()
		return nil, errors.New("存在相同的名称")
	}

	changeMap := make(map[string]any)
	if len(m.Put.Name) > 0 {
		changeMap["Name"] = m.Put.Name
	}
	if len(m.Put.Label) > 0 {
		changeMap["Label"] = m.Put.Label
	}
	if len(m.Put.Image) > 0 {
		changeMap["Image"] = m.Put.Image
	}

	if len(changeMap) > 0 {
		err := dao.UpdateByPrimaryKey(tx, &model.GoodsSkuLabelData{}, has.ID, changeMap)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return result.NewData(map[string]any{"SkuLabelDataList": m.listData(has.GoodsID, 0, "", m.Organization.ID)}), nil
}

func (m *SkuLabelData) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	if m.Post.GoodsSkuLabelID == 0 {
		return nil, errors.New("数据出错")
	}

	tx := db.Orm().Begin()
	has := dao.GetBy(tx, &model.GoodsSkuLabelData{}, map[string]any{"GoodsSkuLabelID": m.Post.GoodsSkuLabelID, "Name": m.Post.Name})
	if !has.IsZero() {
		tx.Rollback()
		return nil, errors.New("存在相同的名称")
	} else {
		goodsSkuLabel := &model.GoodsSkuLabelData{
			OID:             m.Organization.ID,
			GoodsSkuLabelID: m.Post.GoodsSkuLabelID,
			Label:           m.Post.Label,
			Name:            m.Post.Name,
			Image:           m.Post.Image,
			GoodsID:         m.Post.GoodsID,
		}
		err := dao.Create(tx, goodsSkuLabel)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	return result.NewData(map[string]any{"SkuLabelDataList": m.listData(m.Post.GoodsID, 0, "", m.Organization.ID)}), err
}
