package goods

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/tool/object"

	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
)

type ChangeGoods struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		GoodsJSON          string `form:"goods"`
		SpecificationsJSON string `form:"specifications"`
		ParamsJSON         string `form:"params"`
	} `method:"Post"`
}

func (m *ChangeGoods) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

var spaceRegexp = regexp.MustCompile(`\s`)

func (m *ChangeGoods) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {

	var item *model.Goods
	item, err = util.JSONToStruct[*model.Goods](m.Post.GoodsJSON)
	if err != nil {
		return nil, err
	}

	var specifications []model.Specification
	if len(m.Post.SpecificationsJSON) > 0 {
		specifications, err = util.JSONToStruct[[]model.Specification](m.Post.SpecificationsJSON)
		if err != nil {
			return nil, err
		}
	}

	tx := db.GetDB(ctx).Begin()
	{
		//生成标签
		item.Tags = make(pq.StringArray, 0)
		var splitWord string
		dd, err := tx.DB()
		if err != nil {
			return nil, err
		}
		query := dd.QueryRow(fmt.Sprintf(`select to_tsvector('%s'::text)`, strings.ReplaceAll(item.Title, "'", "''")))
		err = query.Scan(&splitWord)
		if err != nil {
			return nil, err
		}
		words := spaceRegexp.Split(splitWord, -1)

		query = dd.QueryRow(fmt.Sprintf(`select to_tsvector('%s'::text)`, strings.ReplaceAll(item.Introduce, "'", "''")))
		err = query.Scan(&splitWord)
		if err != nil {
			return nil, err
		}
		words = append(words, spaceRegexp.Split(splitWord, -1)...)

		for i := range words {
			word := strings.ReplaceAll(strings.Split(words[i], ":")[0], "'", "")
			var has = false
			for i2 := range item.Tags {
				if strings.EqualFold(word, item.Tags[i2]) {
					has = true
					break
				}
			}
			if !has && len(word) >= 3 && object.ParseFloat(word) == 0.0 {
				item.Tags = append(item.Tags, word)
			}
		}
	}
	hasGoods, err := service.Goods.Goods.SaveGoods(ctx, tx, m.Organization.ID, item, specifications)
	if err != nil {
		tx.Rollback()
		as := &result.ActionResult{}
		as.Code = -55
		as.Message = err.Error()
		as.Data = map[string]interface{}{"Goods": hasGoods}
		return &result.JsonResult{Data: as}, err
	}
	tx.Commit()

	as := &result.ActionResult{}
	as.Message = "修改成功"
	as.Data = map[string]interface{}{"Goods": item}
	return &result.JsonResult{Data: as}, err

}
