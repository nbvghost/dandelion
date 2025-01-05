package shipping

import (
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Shipping struct {
	Organization *model.Organization `mapping:""`

	Get struct {
		ID uint `form:"ID"`
	} `method:"get"`

	Delete struct {
		ID uint `form:"ID"`
	} `method:"Delete"`

	Put struct {
		OrdersShipping model.OrdersShipping
	} `method:"Put"`
}

func (g *Shipping) Handle(context constrain.IContext) (constrain.IResult, error) {
	v := dao.GetBy(db.Orm(), &model.OrdersShipping{}, map[string]any{"ID": g.Get.ID, "OID": g.Organization.ID})
	return result.NewData(v), nil
}
func (g *Shipping) HandleDelete(context constrain.IContext) (constrain.IResult, error) {
	err := dao.DeleteBy(db.Orm(), &model.OrdersShipping{}, map[string]any{"ID": g.Delete.ID, "OID": g.Organization.ID})
	if err != nil {
		return nil, err
	}
	return nil, err
}
func (g *Shipping) HandlePut(context constrain.IContext) (constrain.IResult, error) {
	v := dao.GetBy(db.Orm(), &model.OrdersShipping{}, map[string]any{"ID": g.Put.OrdersShipping.ID, "OID": g.Organization.ID})
	if v.IsZero() {
		return nil, result.NewError(errors.New("找不到数据"))
	}
	err := dao.UpdateByPrimaryKey(db.Orm(), &model.OrdersShipping{}, v.Primary(), &model.OrdersShipping{
		Title: g.Put.OrdersShipping.Title,
		Image: g.Put.OrdersShipping.Image,
		No:    g.Put.OrdersShipping.No,
		Name:  g.Put.OrdersShipping.Name,
		Key:   g.Put.OrdersShipping.Key,
	})
	return nil, err
}
