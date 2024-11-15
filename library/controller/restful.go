package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/com"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"net/url"
	"reflect"
	"strings"

	"github.com/nbvghost/tool/object"
)

type IAdminTable interface {
	GetOID() dao.PrimaryKey
}

type Restful[T IAdminTable] struct {
	Admin T `mapping:""`
	Get   struct {
		Model       string `uri:"model"`
		PageNo      int    `form:"Page-No"`
		PageSize    int    `form:"Page-Size"`
		OrderField  string `form:"Order-Field"`
		OrderMethod string `form:"Order-Method"`
	} `method:"get"`
	Post struct {
		Model string `uri:"model"`
		Body  any    `body:""`
	} `method:"post"`
	Put struct {
		Model string         `uri:"model"`
		ID    dao.PrimaryKey `uri:"id"`
		Body  any            `body:""`
	} `method:"put"`
	Del struct {
		Model string         `uri:"model"`
		ID    dao.PrimaryKey `uri:"id"`
	} `method:"del"`
	Query url.Values
}

func (m *Restful[T]) bindQuery(ctx constrain.IContext) {
	contextValue := contexext.FromContext(ctx)
	query := contextValue.Request.URL.Query()
	m.Query = url.Values{}
	for key := range query {
		if strings.Contains(strings.ToLower(key), "page-") {
			continue
		}
		if strings.Contains(strings.ToLower(key), "order-") {
			continue
		}

		v := query[key]
		var hasList []string
		for i := range v {
			if len(v[i]) > 0 {
				hasList = append(hasList, v[i])
			}
		}

		if len(hasList) > 0 {
			m.Query[key] = hasList
		}
	}
}
func (m *Restful[T]) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	m.bindQuery(ctx)
	table, err := entity.GetModel(m.Get.Model)
	if err != nil {
		return nil, err
	}

	tableType := reflect.TypeOf(table).Elem()

	d := dao.Find(db.Orm(), table)
	//分析url参数作为where条件
	for k := range m.Query {
		if len(m.Query[k]) > 1 {
			d.Where(fmt.Sprintf(`"%s" in ?`, k), m.Query[k])
		} else {
			field, ok := tableType.FieldByName(k)
			if !ok {
				continue
			}
			val := m.Query.Get(k)

			if strings.EqualFold(field.Type.Name(), "PrimaryKey") {
				if object.ParseUint(val) > 0 {
					d.Where(fmt.Sprintf(`"%s"='%s'`, k, val))
				}
				continue
			}

			if strings.Contains(val, "%") {
				d.Where(fmt.Sprintf(`"%s" like '%s'`, k, val))
			} else {
				d.Where(fmt.Sprintf(`"%s"='%s'`, k, val))
			}
		}
	}

	_, exist := tableType.FieldByName("OID")
	if exist {
		d.Where(`"OID"=?`, m.Admin.GetOID())
	}

	index := m.Get.PageNo - 1
	if index < 0 {
		index = 0
	}
	total := d.Limit(index, m.Get.PageSize)
	if len(m.Get.OrderField) > 0 && len(m.Get.OrderMethod) > 0 {
		d.Order(fmt.Sprintf(`"%s" %s`, m.Get.OrderField, m.Get.OrderMethod))
	}
	list := d.List()
	return result.NewData(result.NewPagination(m.Get.PageNo, m.Get.PageSize, total, list)), nil
}
func (m *Restful[T]) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	m.bindQuery(ctx)
	table, err := entity.GetModel(m.Post.Model)
	if err != nil {
		return nil, err
	}
	/*body, err := json.Marshal(m.Post.Body)
	if err != nil {
		return nil, err
	}*/

	where := map[string]any{}
	for k := range m.Query {
		where[k] = m.Query.Get(k)
	}

	modelType := reflect.TypeOf(table).Elem()

	_, exist := modelType.FieldByName("OID")
	if exist {

	}

	newData := reflect.New(modelType).Interface().(dao.IEntity)

	if body, ok := m.Post.Body.(map[string]any); ok {
		v := reflect.ValueOf(newData).Elem()
		for key := range body {
			fv := v.FieldByName(key)
			if fv.CanSet() == false {
				continue
			}
			setv := reflect.ValueOf(body[key])
			if setv.CanConvert(fv.Type()) {
				fv.Set(setv.Convert(fv.Type()))
			}
		}
	}
	/*err = json.Unmarshal(body, newData)
	if err != nil {
		return nil, err
	}*/

	var primaryId dao.PrimaryKey

	if newData.Primary() > 0 {
		oldData := dao.GetByPrimaryKey(db.Orm(), table, newData.Primary())
		if oldData.IsZero() {
			reflect.ValueOf(newData).Elem().FieldByName("OID").Set(reflect.ValueOf(m.Admin.GetOID()))
			err = dao.Create(db.Orm(), newData)
			primaryId = newData.Primary()
		} else {
			changeMap := com.Diff(oldData, newData)
			where["ID"] = newData.Primary()
			where["OID"] = m.Admin.GetOID()
			err = dao.UpdateBy(db.Orm(), table, changeMap, where)
			primaryId = newData.Primary()
		}
	} else {
		where["OID"] = m.Admin.GetOID()
		oldData := dao.GetBy(db.Orm(), table, where)
		if oldData.IsZero() {
			field := reflect.ValueOf(newData).Elem().FieldByName("OID")
			if field.CanSet() {
				field.Set(reflect.ValueOf(m.Admin.GetOID()))
			}
			err = dao.Create(db.Orm(), newData)
			primaryId = newData.Primary()
		} else {
			var fieldValue string
			if len(where) > 0 {
				for k := range where {
					if strings.EqualFold("OID", k) {
						continue
					}
					fieldValue = object.ParseString(where[k])
					break
				}
			}
			return nil, fmt.Errorf("存在相同的记录:%s", fieldValue)
			//err = dao.UpdateByPrimaryKey(singleton.Orm(), table, oldData.Primary(), newData)
			//primaryId = oldData.Primary()
		}
	}
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{"Item": dao.GetByPrimaryKey(db.Orm(), table, primaryId)}), err
}
func (m *Restful[T]) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	m.bindQuery(ctx)
	table, err := entity.GetModel(m.Put.Model)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(m.Put.Body)
	if err != nil {
		return nil, err
	}

	newMap := make(map[string]any)
	err = json.Unmarshal(body, &newMap)
	if err != nil {
		return nil, err
	}

	where := map[string]any{}
	for k := range m.Query {
		where[k] = m.Query.Get(k)
	}

	oldData := dao.GetByPrimaryKey(db.Orm(), table, m.Put.ID)
	if oldData.IsZero() {
		return nil, errors.New("数据不存在")
	}

	bytes, err := json.Marshal(oldData)
	if err != nil {
		return nil, err
	}
	oldMap := make(map[string]any)
	err = json.Unmarshal(bytes, &oldMap)
	if err != nil {
		return nil, err
	}

	changeMap := com.Diff(&oldMap, &newMap)

	newData := reflect.New(reflect.TypeOf(table).Elem()).Interface().(dao.IEntity)
	/*err = json.Unmarshal(body, newData)
	if err != nil {
		return nil, err
	}*/
	v := reflect.ValueOf(newData).Elem()
	for key := range changeMap {
		fv := v.FieldByName(key)
		setv := reflect.ValueOf(changeMap[key])
		if setv.CanConvert(fv.Type()) {
			fv.Set(setv.Convert(fv.Type()))
		}
		changeMap[key] = v.FieldByName(key).Interface()
	}

	/*newData := reflect.New(reflect.TypeOf(model).Elem()).Interface().(dao.IEntity)
	err = json.Unmarshal(marshal, newData)
	if err != nil {
		return nil, err
	}
	changeMap := com.Diff(oldData, newData)*/

	//changeMap := com.Diff(oldData, m.Put.Body)

	if len(where) > 0 {
		hasData := dao.GetBy(db.Orm(), table, where)
		if hasData.IsZero() == false && hasData.Primary() != m.Put.ID {
			var fieldValue string
			if len(where) > 0 {
				for k := range where {
					if strings.EqualFold("OID", k) {
						continue
					}
					fieldValue = object.ParseString(where[k])
					break
				}
			}
			return nil, fmt.Errorf("存在相同的记录:%s", fieldValue)
		}
	}

	if len(changeMap) > 0 {
		tx := db.Orm().Begin()
		err = dao.UpdateByPrimaryKey(tx, table, m.Put.ID, changeMap)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
	}
	return result.NewData(map[string]any{"Item": dao.GetByPrimaryKey(db.Orm(), table, m.Put.ID)}), err
}
func (m *Restful[T]) HandleDelete(ctx constrain.IContext) (constrain.IResult, error) {
	m.bindQuery(ctx)

	table, err := entity.GetModel(m.Put.Model)
	if err != nil {
		return nil, err
	}
	if m.Del.ID > 0 {
		err = dao.DeleteByPrimaryKey(db.Orm(), table, m.Del.ID)
		if err != nil {
			return nil, err
		}
	} else {
		where := map[string]any{}
		for k := range m.Query {
			where[k] = m.Query.Get(k)
		}
		if len(where) > 0 {
			err = dao.DeleteBy(db.Orm(), table, where)
			if err != nil {
				return nil, err
			}
		}
	}

	return result.NewData(map[string]any{}), nil
}
