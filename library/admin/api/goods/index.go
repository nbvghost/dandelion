package goods

import (
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Index struct {
	SessionMappingData entity.SessionMappingData `mapping:""`
	Get                struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"get"`
	Put struct {
		Goods *model.Goods
	} `method:"Put"`
	Post struct {
		Goods *model.Goods
	} `method:"Post"`
	Delete struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"Delete"`
}

func (m *Index) Handle(context constrain.IContext) (constrain.IResult, error) {
	return result.NewData(map[string]any{"Goods": dao.GetBy(db.Orm(), &model.Goods{}, map[string]any{"OID": m.SessionMappingData.GetOID(), "ID": m.Get.ID})}), nil
}
func (m *Index) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	if m.Post.Goods.ID > 0 {
		return nil, errors.New("数据错误")
	}
	tx := db.Orm().Begin()
	goods, err := service.Goods.Goods.SaveGoods(tx, m.SessionMappingData.GetOID(), m.Post.Goods, nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return result.NewDataMessage(map[string]any{"Goods": goods}, "添加成功"), nil
}
func (m *Index) HandlePut(context constrain.IContext) (constrain.IResult, error) {
	if m.Put.Goods.ID == 0 {
		return nil, errors.New("数据错误")
	}
	has := dao.GetBy(db.Orm(), &model.Goods{}, map[string]any{"OID": m.SessionMappingData.GetOID(), "ID": m.Put.Goods.ID}).(*model.Goods)
	if has.IsZero() {
		return nil, errors.New("找不到数据")
	}

	tx := db.Orm().Begin()
	goods, err := service.Goods.Goods.SaveGoods(tx, m.SessionMappingData.GetOID(), m.Put.Goods, nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return result.NewDataMessage(map[string]any{"Goods": goods}, "修改成功"), nil
}
func (m *Index) HandleDelete(context constrain.IContext) (constrain.IResult, error) {
	return &result.JsonResult{Data: service.Goods.Goods.DeleteGoods(dao.PrimaryKey(m.Delete.ID))}, nil
}
