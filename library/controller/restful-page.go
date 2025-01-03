package controller

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"log"
	"reflect"
)

var conditions = map[string]any{"like": true, "in": true, "=": true, "!=": true, ">": true, "<": true}

type Query struct {
	Field     string //验证
	Condition string //验证,=,!=,>,<,or,
	Value     any
}

type RestfulPage[T IOIDMapping] struct {
	Admin T `mapping:""`
	Post  struct {
		Model    string
		PageNo   int
		PageSize int
		Sort     []dao.Sort
		Query    map[string]any //and,or
	} `method:"post"`
}

type Where struct {
	Query any
	Args  []any
}

func (m *RestfulPage[T]) condition(ctx constrain.IContext, key string, tableType reflect.Type) ([]Where, error) {
	if m.Post.Query == nil {
		return nil, nil
	}
	queryValue := m.Post.Query[key] //
	if queryValue == nil {
		return nil, nil
	}
	valueType := reflect.TypeOf(queryValue)

	whereList := make([]Where, 0)

	var queryList []Query
	if valueType.Kind() == reflect.Map {
		qMap, ok := queryValue.(map[string]any)
		if !ok {
			return nil, nil
		}
		for s, a := range qMap {
			queryList = append(queryList, Query{
				Field:     s,
				Condition: "=",
				Value:     a,
			})
		}
	}
	if valueType.Kind() == reflect.Slice {
		var ok bool
		var arr []any
		arr, ok = queryValue.([]any)
		if !ok {
			return nil, errors.New("无效的查询参数")
		}
		for i := range arr {
			qMap := arr[i].(map[string]any)
			queryList = append(queryList, Query{
				Field:     qMap["Field"].(string),
				Condition: qMap["Condition"].(string),
				Value:     qMap["Value"],
			})
		}
	}

	for index := range queryList {
		query := queryList[index] //[]Query,map[string]any
		{
			if query.Value == nil {
				//空值不处理
				continue
			}
			valueValue := reflect.ValueOf(query.Value)
			if valueValue.IsZero() {
				//0值不处理
				continue
			}

			field, has := tableType.FieldByName(query.Field)
			if !has {
				//不存在的字段不处理，子struct不处理
				continue
				//return nil, errors.New(fmt.Sprintf("not find field <%s>", query.Field))
			}
			log.Println(field.Type.Kind())
			if field.Type.Kind() == reflect.Struct {
				continue
			}
			if field.Type.Kind() == reflect.Array || field.Type.Kind() == reflect.Slice {
				continue
			}
			if _, has := conditions[query.Condition]; has == false {
				//不存在的条件不处理
				//return nil, errors.New(fmt.Sprintf("not find condition <%s>", query.Condition))
				continue
			}
		}

		if query.Condition == "like" {
			whereList = append(whereList, Where{
				Query: fmt.Sprintf(`"%s" ilike ?`, query.Field),
				Args:  []any{query.Value},
			})
		} else {
			whereList = append(whereList, Where{
				Query: fmt.Sprintf(`"%s" %s ?`, query.Field, query.Condition),
				Args:  []any{query.Value},
			})
		}
	}

	return whereList, nil
}

func (m *RestfulPage[T]) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	table, err := entity.GetModel(m.Post.Model)
	if err != nil {
		return nil, err
	}

	tableType := reflect.TypeOf(table).Elem()

	d := db.Orm().Table(m.Post.Model) //dao.Find(db.Orm(), table)

	wheres, err := m.condition(ctx, "and", tableType)
	if err != nil {
		return nil, err
	}
	for i := range wheres {
		where := wheres[i]
		d.Where(where.Query, where.Args...)
	}

	wheres, err = m.condition(ctx, "or", tableType)
	if err != nil {
		return nil, err
	}
	for i := range wheres {
		where := wheres[i]
		d.Or(where.Query, where.Args...)
	}

	_, exist := tableType.FieldByName("OID")
	if exist {
		d.Where(`"OID"=?`, m.Admin.GetOID())
	}

	index := m.Post.PageNo - 1
	if index < 0 {
		index = 0
	}
	var total int64
	if index == 0 && m.Post.PageSize == 0 {
		d.Count(&total)
	} else {
		d.Limit(m.Post.PageSize).Offset(index * m.Post.PageSize).Count(&total)
	}
	for i := range m.Post.Sort {
		sort := m.Post.Sort[i]
		d.Order(fmt.Sprintf(`"%s" %s`, sort.SortField(), sort.SortMethod()))
	}

	var list = reflect.New(reflect.SliceOf(reflect.TypeOf(table)))
	d.Find(list.Interface())
	return result.NewData(result.NewPagination(index+1, m.Post.PageSize, total, list.Interface())), nil
}
func (m *RestfulPage[T]) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
