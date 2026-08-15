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

type Attributes struct {
	SessionMappingData entity.SessionMappingData `mapping:""`
	POST               struct {
		GoodsID dao.PrimaryKey `form:"GoodsID"`
		GroupID dao.PrimaryKey `form:"GroupID"`
		Name    string         `form:"Name"`
		Value   string         `form:"Value"`
	} `method:"POST"`
	Put struct {
		ID    dao.PrimaryKey `form:"ID"`
		Name  string         `form:"Name"`
		Value string         `form:"Value"`
	} `method:"Put"`
	Get struct {
		GoodsID dao.PrimaryKey `form:"GoodsID"`
		GroupID dao.PrimaryKey `form:"GroupID"`
	} `method:"Get"`
	DELETE struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"DELETE"`
}

func (g *Attributes) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	at, err := service.Goods.Attributes.AddGoodsAttributes(ctx, g.SessionMappingData.OID, g.POST.GoodsID, g.POST.GroupID, g.POST.Name, g.POST.Value)
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: result.NewDataMessage(at, "属性添加成功")}, err
}

func (g *Attributes) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	list := make([]model.GoodsAttributes, 0)
	orm := db.GetDB(ctx).Model(&model.GoodsAttributes{})
	orm.Where(`"OID"=?`, g.SessionMappingData.OID)
	orm.Where(`"GoodsID"=?`, g.Get.GoodsID)
	if g.Get.GroupID > 0 {
		orm.Where(`"GroupID"=?`, g.Get.GroupID)
	}
	orm.Find(&list)
	return result.NewData(result.NewPagination(0, 0, 0, list)), nil
}

func (g *Attributes) HandleDelete(ctx constrain.IContext) (r constrain.IResult, err error) {
	err = service.Goods.Attributes.DeleteGoodsAttributes(ctx, g.DELETE.ID)
	return &result.JsonResult{Data: result.New(err, "属性删除成功")}, err
}
func (g *Attributes) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	goodsAttributes := dao.GetBy(db.GetDB(ctx), &model.GoodsAttributes{}, map[string]any{"ID": g.Put.ID, "OID": g.SessionMappingData.OID}).(*model.GoodsAttributes)
	if goodsAttributes.IsZero() {
		return nil, errors.New("找不到数据")
	}

	err := dao.UpdateByPrimaryKey(db.GetDB(ctx), &model.GoodsAttributes{}, goodsAttributes.ID, map[string]any{"Value": g.Put.Value})
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: result.New(err, "属性更新成功")}, err
}
