package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"strings"
)

type List struct {
	Admin *entity.SessionMappingData `mapping:""`
	Post  struct {
		Query    ListQuery
		Sort     dao.Sort
		PageNo   int
		PageSize int
	} `method:"Post"`
}

type ListQuery struct {
	Keyword       string
	Title         string
	GoodsTypeID   dao.PrimaryKey
	Introduce     string
	Specification struct {
		Label string
	}
}

func (m *List) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
func (m *List) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	orm := db.Orm().Model(&model.Goods{}).
		Joins(`left join "GoodsType" on ("GoodsType"."ID" = "Goods"."GoodsTypeID")
         left join "Specification" on ("Specification"."GoodsID" = "Goods"."ID")`).
		Where(`"Goods"."OID"=?`, m.Admin.OID).Select(`"Goods".*,"GoodsType".*,to_json(array_agg("Specification")) as "SpecificationList"`).
		Group(`"Goods"."ID"`).
		Group(`"GoodsType"."ID"`).
		Order(m.Post.Sort.OrderByColumn(`"Goods"."CreatedAt"`, true))

	if len(m.Post.Query.Title) > 0 {
		orm.Where(`"Goods"."Title" ilike ?`, fmt.Sprintf("%%%s%%", m.Post.Query.Title))
	}
	if m.Post.Query.GoodsTypeID > 0 {
		orm.Where(`"Goods"."GoodsTypeID" = ?`, m.Post.Query.GoodsTypeID)
	}
	if len(m.Post.Query.Introduce) > 0 {
		orm.Where(`"Goods"."Introduce" ilike ?`, fmt.Sprintf("%%%s%%", m.Post.Query.Introduce))
	}
	if len(m.Post.Query.Specification.Label) > 0 {
		orm.Where(`"Specification"."Label" ilike ? or "Specification"."Language"->>'Label' ilike ?`, fmt.Sprintf("%%%s%%", m.Post.Query.Specification.Label), fmt.Sprintf("%%%s%%", m.Post.Query.Specification.Label))
	}

	keyword := strings.TrimSpace(m.Post.Query.Keyword)
	if len(keyword) > 0 {
		orm.Where(`
"Goods"."Title" ilike ? or
"Goods"."Summary" ilike ? or
"Goods"."Introduce" ilike ? or
"Goods"."Language"->>'Title' ilike ? or
"Goods"."Language"->>'Summary' ilike ? or
"Goods"."Language"->>'Introduce' ilike ? or
"Specification"."Label" ilike ? or 
"Specification"."CodeNo" ilike ? or 
"Specification"."Language"->>'Label' ilike ?
`,
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword))
	}

	var total int64
	orm.Count(&total)

	var list []struct {
		model.Goods
		model.GoodsType   `json:"GoodsType"`
		SpecificationList sqltype.Array[model.Specification]
	}
	orm.Limit(m.Post.PageSize).Offset((m.Post.PageNo - 1) * m.Post.PageSize).Find(&list)

	return result.NewData(result.NewPagination(m.Post.PageNo, m.Post.PageSize, total, list)), nil
}
