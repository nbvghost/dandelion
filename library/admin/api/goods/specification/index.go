package specification

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/tool/object"
	"github.com/samber/lo"
	"strings"
)

type Index struct {
	OIDMapping *entity.SessionMappingData `mapping:""`
	Get        struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"get"`
	Delete struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"Delete"`
	Post struct {
		Specification struct {
			ID          dao.PrimaryKey
			GoodsID     dao.PrimaryKey
			Label       string
			Num         uint
			Unit        string
			Weight      float64
			Stock       uint
			CostPrice   float64
			MarketPrice float64
			Brokerage   float64
			Pictures    sqltype.Array[sqltype.Image]
			Language    model.SpecificationLanguage
			Remark      string
		}
	} `method:"Post"`
	Put struct {
		Specification struct {
			ID          dao.PrimaryKey
			GoodsID     dao.PrimaryKey
			Label       string
			Num         uint
			Unit        string
			Weight      float64
			Stock       uint
			CostPrice   float64
			MarketPrice float64
			Brokerage   float64
			Pictures    sqltype.Array[sqltype.Image]
			Language    model.SpecificationLanguage
			Remark      string
		}
	} `method:"Put"`
}

func (g *Index) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	return result.NewData(dao.GetByPrimaryKey(db.Orm(), &model.Specification{}, g.Get.ID)), err
}
func (g *Index) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	if g.Put.Specification.ID == 0 {
		return nil, errors.New("数据错误")
	}

	s := make(map[string]any)
	if len(g.Put.Specification.Label) > 0 {
		s["Label"] = g.Put.Specification.Label
	}
	/*if len(g.Put.Name) > 0 {
		s["Name"] = g.Put.Name
	}*/
	/*if len(g.Put.Specification.LabelIndex) > 0 {
		s["LabelIndex"] = g.Put.Specification.LabelIndex
	}*/
	if g.Put.Specification.Num > 0 {
		s["Num"] = g.Put.Specification.Num
	}
	if len(g.Put.Specification.Unit) > 0 {
		s["Unit"] = g.Put.Specification.Unit
	}
	if g.Put.Specification.Weight > 0 {
		s["Weight"] = object.ParseUint(g.Put.Specification.Weight * 1000.0)
	}
	if g.Put.Specification.Stock > 0 {
		s["Stock"] = g.Put.Specification.Stock
	}
	if g.Put.Specification.CostPrice > 0 {
		s["CostPrice"] = object.ParseUint(g.Put.Specification.CostPrice * 100.0)
	}
	if g.Put.Specification.MarketPrice > 0 {
		s["MarketPrice"] = object.ParseUint(g.Put.Specification.MarketPrice * 100.0)
	}
	if g.Put.Specification.Brokerage > 0 {
		s["Brokerage"] = object.ParseUint(g.Put.Specification.Brokerage * 100.0)
	}
	if len(g.Put.Specification.Remark) > 0 {
		s["Remark"] = g.Put.Specification.Remark
	}
	if len(g.Put.Specification.Language.Label) > 0 {
		s["Language"] = g.Put.Specification.Language
	}

	if len(g.Put.Specification.Pictures) > 0 {
		for i := 0; i < len(g.Put.Specification.Pictures); i++ {
			item := g.Put.Specification.Pictures[i]
			if strings.HasPrefix(item.Src, oss.TempFilePrefix) {
				fileBytes, err := oss.GetTempFile(item.Src)
				if err != nil {
					//tx.Rollback()
					return nil, err
				}
				file, err := oss.UploadFile(context, fileBytes, fmt.Sprintf("goods/specification/%d", g.Put.Specification.ID), "", true, "")
				if err != nil {
					//tx.Rollback()
					return nil, err
				}
				g.Put.Specification.Pictures[i].Src = file.Data.Path
			}
		}
	}
	s["Pictures"] = g.Put.Specification.Pictures

	err = dao.UpdateByPrimaryKey(db.Orm(), &model.Specification{}, g.Put.Specification.ID, s)
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{"Specifications": dao.Find(db.Orm(), &model.Specification{}).Where(`"GoodsID"=?`, g.Put.Specification.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
func (g *Index) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	if g.Post.Specification.GoodsID == 0 {
		return nil, errors.New("数据错误")
	}

	tx := db.Orm().Begin()

	specificationList := dao.Find(tx, &model.Specification{}).Where(map[string]any{"GoodsID": g.Post.Specification.GoodsID}).List()

	/*for i := range specificationList {
		item := specificationList[i].(*model.Specification)
		var hasIndex = false
		for ii := range g.Post.List {
			items := g.Post.List[ii]
			if item.Label == items.Label {
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
	}*/

	specification := g.Post.Specification

	hasItem, has := lo.Find(specificationList, func(item dao.IEntity) bool {
		s := item.(*model.Specification)
		if strings.EqualFold(s.Label, specification.Label) {
			return true
		} else {
			return false
		}
	})
	if has {
		return nil, errors.New(fmt.Sprintf("存在相同的规格:%s", hasItem.(*model.Specification).Label))
	}

	newSpecification := &model.Specification{
		OID:     g.OIDMapping.OID,
		GoodsID: g.Post.Specification.GoodsID,
		//Name:        item.Name,
		Label: specification.Label,
		//LabelIndex:  specification.LabelIndex,
		Num:         specification.Num,
		Unit:        specification.Unit,
		Weight:      object.ParseUint(specification.Weight * 1000.0),
		Stock:       specification.Stock,
		CostPrice:   object.ParseUint(specification.CostPrice * 100.0),
		MarketPrice: object.ParseUint(specification.MarketPrice * 100.0),
		Brokerage:   object.ParseUint(specification.Brokerage * 100.0),
		Language:    specification.Language,
		Remark:      specification.Remark,
	}
	err = dao.Create(tx, newSpecification)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for i := 0; i < len(specification.Pictures); i++ {
		item := specification.Pictures[i]
		if strings.HasPrefix(item.Src, oss.TempFilePrefix) {
			fileBytes, err := oss.GetTempFile(item.Src)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			file, err := oss.UploadFile(context, fileBytes, fmt.Sprintf("goods/specification/%d", newSpecification.ID), "", true, "")
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			specification.Pictures[i].Src = file.Data.Path
		}
	}
	err = dao.UpdateByPrimaryKey(tx, &model.Specification{}, newSpecification.ID, map[string]any{"Pictures": specification.Pictures})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return result.NewData(map[string]any{"Specifications": dao.Find(db.Orm(), &model.Specification{}).Where(`"GoodsID"=?`, g.Post.Specification.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
func (g *Index) HandleDelete(context constrain.IContext) (r constrain.IResult, err error) {
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
