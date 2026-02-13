package goods

import (
	"errors"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type SkuLabel struct {
	//Goods        goods.GoodsService
	Organization *model.Organization `mapping:""`
	Get          struct {
		GoodsID dao.PrimaryKey `form:"goods-id"`
	} `method:"Get"`
	Put struct {
		LabelList []struct {
			ID    dao.PrimaryKey
			Abel  bool
			Image bool
			Label string
			Name  string
		}
		GoodsID dao.PrimaryKey
	} `method:"Put"`
	Delete struct {
		ID dao.PrimaryKey `form:"id"`
	} `method:"Delete"`
}

func (m *SkuLabel) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	where := make(map[string]any)
	where["OID"] = m.Organization.ID
	if m.Get.GoodsID > 0 {
		where["GoodsID"] = m.Get.GoodsID
	}
	LabelList := dao.Find(db.GetDB(ctx), &model.GoodsSkuLabel{}).Where(where).List()
	return result.NewData(map[string]any{"SkuLabelList": LabelList}), nil
}
func (m *SkuLabel) HandleDelete(ctx constrain.IContext) (constrain.IResult, error) {
	if m.Delete.ID == 0 {
		return nil, errors.New("数据错误")
	}
	has := dao.GetByPrimaryKey(db.GetDB(ctx), &model.GoodsSkuLabel{}, m.Delete.ID).(*model.GoodsSkuLabel)
	if has.IsZero() {
		return nil, errors.New("找不到数据")
	}

	tx := db.GetDB(ctx).Begin()

	err := dao.DeleteByPrimaryKey(tx, &model.GoodsSkuLabel{}, has.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	specificationList := dao.Find(tx, &model.Specification{}).Where(`"GoodsID"=?`, has.GoodsID).List()

	goodsSkuLabelDataList := dao.Find(tx, &model.GoodsSkuLabelData{}).Where(`"GoodsSkuLabelID"=?`, has.ID).List()
	for i := range goodsSkuLabelDataList {
		item := goodsSkuLabelDataList[i].(*model.GoodsSkuLabelData)

		for i2 := range specificationList {
			specification := specificationList[i2].(*model.Specification)
			hasItem := false
			for i3 := range specification.LabelIndex {
				if specification.LabelIndex[i3] == item.ID {
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
		err := dao.DeleteByPrimaryKey(tx, &model.GoodsSkuLabelData{}, item.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

	}
	tx.Commit()

	LabelList := dao.Find(db.GetDB(ctx), &model.GoodsSkuLabel{}).Where(`"GoodsID"=?`, has.GoodsID).List()
	return result.NewData(map[string]any{"SkuLabelList": LabelList}), nil
}
func (m *SkuLabel) HandlePut(ctx constrain.IContext) (r constrain.IResult, err error) {
	tx := db.GetDB(ctx).Begin()
	for _, label := range m.Put.LabelList {
		if label.ID > 0 {
			err := dao.UpdateByPrimaryKey(tx, &model.GoodsSkuLabel{}, label.ID, map[string]any{
				"Name":  label.Name,
				"Label": label.Label,
				"Abel":  label.Abel,
				"Image": label.Image,
			})
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			goodsSkuLabel := dao.GetBy(tx, &model.GoodsSkuLabel{}, map[string]any{"Name": label.Name, "GoodsID": m.Put.GoodsID}).(*model.GoodsSkuLabel)
			if goodsSkuLabel.IsZero() {
				goodsSkuLabel = &model.GoodsSkuLabel{
					OID:     m.Organization.ID,
					GoodsID: m.Put.GoodsID,
					Label:   label.Label,
					Name:    label.Name,
					Abel:    label.Abel,
					Image:   label.Image,
				}
				err := dao.Create(tx, goodsSkuLabel)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			} else {
				goodsSkuLabel.Label = label.Label
				goodsSkuLabel.Abel = label.Abel
				goodsSkuLabel.Image = label.Image
				err := dao.UpdateByPrimaryKey(tx, &model.GoodsSkuLabel{}, goodsSkuLabel.ID, map[string]any{"Label": goodsSkuLabel.Label, "Abel": goodsSkuLabel.Abel, "Image": goodsSkuLabel.Image})
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	tx.Commit()
	LabelList := dao.Find(db.GetDB(ctx), &model.GoodsSkuLabel{}).Where(`"GoodsID"=?`, m.Put.GoodsID).List()
	return result.NewData(map[string]any{"SkuLabelList": LabelList}), err
}
