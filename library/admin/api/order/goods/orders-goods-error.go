package goods

import (
	"errors"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type OrdersGoodsError struct {
	SessionMappingData *entity.SessionMappingData `mapping:""`
	Post               struct {
		ID    dao.PrimaryKey
		Error string
	} `method:"post"`
	Del struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"Delete"`
}

func (m *OrdersGoodsError) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *OrdersGoodsError) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	var order model.OrdersGoods
	db.GetDB(ctx).Model(&model.OrdersGoods{}).Where(`"OID"=? and "ID"=?`, m.SessionMappingData.OID, m.Post.ID).Find(&order)
	if order.IsZero() {
		return nil, errors.New("找不到数据")
	}
	err := dao.UpdateByPrimaryKey(db.GetDB(ctx), &model.OrdersGoods{}, m.Post.ID, map[string]any{"Error": m.Post.Error})
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("添加成功"), nil
}
func (m *OrdersGoodsError) HandleDelete(ctx constrain.IContext) (constrain.IResult, error) {
	var order model.OrdersGoods
	db.GetDB(ctx).Model(&model.OrdersGoods{}).Where(`"OID"=? and "ID"=?`, m.SessionMappingData.OID, m.Del.ID).Find(&order)
	if order.IsZero() {
		return nil, errors.New("找不到数据")
	}
	err := dao.UpdateByPrimaryKey(db.GetDB(ctx), &model.OrdersGoods{}, m.Post.ID, map[string]any{"Error": ""})
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("清空成功"), nil
}
