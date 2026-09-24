package goods

import (
	"fmt"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
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
	Keyword                string
	Code                   string
	GoodsTypeID            dao.PrimaryKey
	LocationID             dao.PrimaryKey
	CreatedAts             []string
	OnlyMatchSpecification bool
}

func (m *List) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
func (m *List) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	orm := db.GetDB(ctx).Model(&model.Goods{})
	orm.Joins(`left join "GoodsType" on ("GoodsType"."ID" = "Goods"."GoodsTypeID")`)
	//Joins(`left join "Specification" on ("Specification"."GoodsID" = "Goods"."ID")`).

	orm.Joins(`left join (SELECT "GoodsID",(to_json(array_agg("GoodsAttributes"))) AS "GoodsAttributesList" FROM "GoodsAttributes" GROUP BY "GoodsID") "GoodsAttributes" ON ("GoodsAttributes".
"GoodsID" = "Goods"."ID")`)
	orm.Joins(`left join (SELECT "GoodsID",(to_json(array_agg("AssetCode"))) AS "AssetCodeList" FROM "AssetCode" GROUP BY "GoodsID") "AssetCode" ON ("AssetCode"."GoodsID" = "Goods"."ID")`)
	orm.Where(`"Goods"."OID"=?`, m.Admin.OID).Select(`"Goods".*,"GoodsType".*,COALESCE("Specification"."SpecificationList", '[]'::json) as "SpecificationList",COALESCE("AssetCode"."AssetCodeList", 
'[]'::json) as "AssetCodeList",COALESCE("GoodsAttributes"."GoodsAttributesList", '[]'::json) as "GoodsAttributesList"`)
	//Group(`"Goods"."ID"`).
	//Group(`"GoodsType"."ID"`).
	orm.Order(m.Post.Sort.OrderByColumn(`"Goods"."CreatedAt"`, true))

	if m.Post.Query.GoodsTypeID > 0 {
		orm.Where(`"Goods"."GoodsTypeID" = ?`, m.Post.Query.GoodsTypeID)
	}
	if m.Post.Query.LocationID > 0 {
		orm.Where(`EXISTS (SELECT 1 FROM "AssetCode" ac WHERE ac."GoodsID" = "Goods"."ID" AND ac."LocationID" = ?)`, m.Post.Query.LocationID)
	}

	keyword := strings.TrimSpace(m.Post.Query.Keyword)
	if len(keyword) > 0 {

		if m.Post.Query.OnlyMatchSpecification {
			orm.Joins(`left join (SELECT "GoodsID",(to_json(array_agg("Specification"))) AS "SpecificationList" FROM "Specification" where "Specification"."Label" ilike ? or "Specification".
"CodeNo" ilike ? or "Specification"."Language" ->> 'Label' ilike ? GROUP BY "GoodsID") "Specification" ON ("Specification"."GoodsID" = "Goods"."ID")`, fmt.Sprintf("%%%s%%", keyword),
				fmt.Sprintf("%%%s%%", keyword), fmt.Sprintf("%%%s%%", keyword))
		} else {
			orm.Joins(`left join (SELECT "GoodsID",(to_json(array_agg("Specification"))) AS "SpecificationList" FROM "Specification" GROUP BY "GoodsID") "Specification" ON ("Specification"."GoodsID" = "Goods"."ID")`)
		}
		orm.Where(`
"Goods"."Title" ilike ? or
"Goods"."Summary" ilike ? or
"Goods"."Introduce" ilike ? or
"Goods"."Language"->>'Title' ilike ? or
"Goods"."Language"->>'Summary' ilike ? or
"Goods"."Language"->>'Introduce' ilike ? or
EXISTS (SELECT 1 FROM "Specification" ac WHERE ac."GoodsID" = "Goods"."ID" AND ac."Label" ilike ?) OR
EXISTS (SELECT 1 FROM "Specification" ac WHERE ac."GoodsID" = "Goods"."ID" AND ac."CodeNo" ilike ?) OR
EXISTS (SELECT 1 FROM "Specification" ac WHERE ac."GoodsID" = "Goods"."ID" AND ac."Language"->>'Label' ilike ?) OR
EXISTS (SELECT 1 FROM "GoodsAttributes" ac WHERE ac."GoodsID" = "Goods"."ID" AND ac."Name" ilike ?) OR
EXISTS (SELECT 1 FROM "GoodsAttributes" ac WHERE ac."GoodsID" = "Goods"."ID" AND ac."Value" ilike ?)
`,
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword))
	} else {
		orm.Joins(`left join (SELECT "GoodsID",(to_json(array_agg("Specification"))) AS "SpecificationList" FROM "Specification" GROUP BY "GoodsID") "Specification" ON ("Specification"."GoodsID" = "Goods"."ID")`)
	}

	if len(m.Post.Query.CreatedAts) == 2 {
		InventoryAtStart, err := time.Parse(time.DateOnly, m.Post.Query.CreatedAts[0])
		if err != nil {
			return nil, err
		}
		InventoryAtStart = time.Date(InventoryAtStart.Year(), InventoryAtStart.Month(), InventoryAtStart.Day(), 0, 0, 0, 0, time.Local)

		InventoryAtEnd, err := time.Parse(time.DateOnly, m.Post.Query.CreatedAts[1])
		if err != nil {
			return nil, err
		}
		InventoryAtEnd = time.Date(InventoryAtEnd.Year(), InventoryAtEnd.Month(), InventoryAtEnd.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
		orm.Where(`"Goods"."CreatedAt">=? and "Goods"."CreatedAt"<?`, InventoryAtStart, InventoryAtEnd)
	}

	var total int64
	orm.Count(&total)

	var list []struct {
		model.Goods
		model.GoodsType     `json:"GoodsType"`
		SpecificationList   sqltype.Array[model.Specification]
		AssetCodeList       sqltype.Array[model.AssetCode]
		GoodsAttributesList sqltype.Array[model.GoodsAttributes]
	}
	orm.Limit(m.Post.PageSize).Offset((m.Post.PageNo - 1) * m.Post.PageSize).Find(&list)

	return result.NewData(result.NewPagination(m.Post.PageNo, m.Post.PageSize, total, list)), nil
}
