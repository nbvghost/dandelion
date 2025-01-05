package goods

import (
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Specification struct {
	Organization *model.Organization `mapping:""`
	Get          struct {
		GoodsID uint `form:"goods-id"`
	} `method:"get"`
	Delete struct {
		ID dao.PrimaryKey `form:"id"`
	} `method:"Delete"`
	Post struct {
		GoodsID dao.PrimaryKey
		List    []struct {
			Label       string
			Name        string
			LabelIndex  sqltype.Array[dao.PrimaryKey]
			Num         uint
			Weight      uint
			Stock       uint
			CostPrice   uint
			MarketPrice uint
			Brokerage   uint
		}
	} `method:"Post"`
	Put model.Specification `method:"Put"`
}

func (g *Specification) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	return result.NewData(map[string]any{"Specifications": dao.Find(db.Orm(), &model.Specification{}).Where(`"GoodsID"=?`, g.Get.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
func (g *Specification) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	if g.Put.ID == 0 {
		return nil, errors.New("数据错误")
	}

	s := make(map[string]any)
	if len(g.Put.Label) > 0 {
		s["Label"] = g.Put.Label
	}
	if len(g.Put.Name) > 0 {
		s["Name"] = g.Put.Name
	}
	if len(g.Put.LabelIndex) > 0 {
		s["LabelIndex"] = g.Put.LabelIndex
	}
	if g.Put.Num > 0 {
		s["Num"] = g.Put.Num
	}
	if g.Put.Weight > 0 {
		s["Weight"] = g.Put.Weight
	}
	if g.Put.Stock > 0 {
		s["Stock"] = g.Put.Stock
	}
	if g.Put.CostPrice > 0 {
		s["CostPrice"] = g.Put.CostPrice
	}
	if g.Put.MarketPrice > 0 {
		s["MarketPrice"] = g.Put.MarketPrice
	}
	if g.Put.Brokerage > 0 {
		s["Brokerage"] = g.Put.Brokerage
	}

	err = dao.UpdateByPrimaryKey(db.Orm(), &model.Specification{}, g.Put.ID, s)
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{"Specifications": dao.Find(db.Orm(), &model.Specification{}).Where(`"GoodsID"=?`, g.Put.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
func (g *Specification) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	if g.Post.GoodsID == 0 {
		return nil, errors.New("数据错误")
	}

	tx := db.Orm().Begin()

	specificationList := dao.Find(tx, &model.Specification{}).Where(map[string]any{"GoodsID": g.Post.GoodsID}).List()

	for i := range specificationList {
		item := specificationList[i].(*model.Specification)
		var hasIndex = false
		for ii := range g.Post.List {
			items := g.Post.List[ii]
			if item.Name == items.Name {
				hasIndex = true
				break
			}
		}
		if !hasIndex {
			err = dao.DeleteByPrimaryKey(tx, &model.Specification{}, item.ID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	for i := range g.Post.List {
		item := g.Post.List[i]
		var hasIndex = -1
		for ii := range specificationList {
			items := specificationList[ii].(*model.Specification)
			if item.Name == items.Name {
				hasIndex = ii
				break
			}
		}
		if hasIndex == -1 {
			err = dao.Create(tx, &model.Specification{
				OID:         g.Organization.ID,
				GoodsID:     g.Post.GoodsID,
				Name:        item.Name,
				Label:       item.Label,
				LabelIndex:  item.LabelIndex,
				Num:         item.Num,
				Weight:      item.Weight,
				Stock:       item.Stock,
				CostPrice:   item.CostPrice,
				MarketPrice: item.MarketPrice,
				Brokerage:   item.Brokerage,
			})
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	tx.Commit()
	return result.NewData(map[string]any{"Specifications": dao.Find(db.Orm(), &model.Specification{}).Where(`"GoodsID"=?`, g.Post.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
func (g *Specification) HandleDelete(context constrain.IContext) (r constrain.IResult, err error) {
	if g.Delete.ID == 0 {
		return nil, errors.New("没找到记录")
	}
	has := dao.GetByPrimaryKey(db.Orm(), &model.Specification{}, g.Delete.ID).(*model.Specification)
	err = service.Goods.Specification.DeleteSpecification(has.ID)
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{"Specifications": dao.Find(db.Orm(), &model.Specification{}).Where(`"GoodsID"=?`, has.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
