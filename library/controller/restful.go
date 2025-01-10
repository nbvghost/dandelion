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
	"gorm.io/gorm"
	"net/url"
	"reflect"
	"strings"

	"github.com/nbvghost/tool/object"
)

type Restful struct {
	Admin *entity.SessionMappingData `mapping:""`
	Get   struct {
		Model string         `uri:"model"`
		ID    dao.PrimaryKey `uri:"id"`
		//PageNo      int    `form:"Page-No"`
		//PageSize    int    `form:"Page-Size"`
		//OrderField  string `form:"Order-Field"`
		//OrderMethod string `form:"Order-Method"`
	} `method:"get"`
	Post struct {
		Model string `uri:"model"`
		Body  []byte `body:""`
	} `method:"post"`
	Put struct {
		Model string         `uri:"model"`
		ID    dao.PrimaryKey `uri:"id"`
		Body  []byte         `body:""`
	} `method:"put"`
	Del struct {
		Model string         `uri:"model"`
		ID    dao.PrimaryKey `uri:"id"`
	} `method:"del"`
}

func (m *Restful) bindQuery(ctx constrain.IContext, tableType reflect.Type) url.Values {
	contextValue := contexext.FromContext(ctx)
	query := contextValue.Request.URL.Query()
	queryValues := url.Values{}
	for key := range query {
		if strings.Contains(strings.ToLower(key), "page-") {
			continue
		}
		if strings.Contains(strings.ToLower(key), "order-") {
			continue
		}
		_, ok := tableType.FieldByName(key)
		if !ok {
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
			queryValues[key] = hasList
		}
	}
	return queryValues
}
func (m *Restful) bindQueryScopes(ctx constrain.IContext, tableType reflect.Type) func(*gorm.DB) *gorm.DB {
	contextValue := contexext.FromContext(ctx)
	query := contextValue.Request.URL.Query()
	queryValues := url.Values{}
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
			queryValues[key] = hasList
		}
	}
	return func(d *gorm.DB) *gorm.DB {

		//分析url参数作为where条件
		for k := range queryValues {
			if len(queryValues[k]) > 1 {
				d.Where(fmt.Sprintf(`"%s" in ?`, k), queryValues[k])
			} else {
				field, ok := tableType.FieldByName(k)
				if !ok {
					continue
				}
				val := queryValues.Get(k)

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

		return d
	}
}

func (m *Restful) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	table, err := entity.GetModel(m.Get.Model)
	if err != nil {
		return nil, err
	}

	tableType := reflect.TypeOf(table).Elem()

	d := db.Orm().Model(table)
	_, exist := tableType.FieldByName("OID")
	if exist {
		d.Where(`"OID"=?`, m.Admin.OID)
	}
	newData := reflect.New(tableType)
	d.First(newData.Interface(), m.Get.ID)
	return result.NewData(newData.Interface()), nil
}
func (m *Restful) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	table, err := entity.GetModel(m.Post.Model)
	if err != nil {
		return nil, err
	}
	/*body, err := json.Marshal(m.Post.Body)
	if err != nil {
		return nil, err
	}*/

	modelType := reflect.TypeOf(table).Elem()

	query := m.bindQuery(ctx, modelType)
	where := map[string]any{}
	for k := range query {
		where[k] = query.Get(k)
	}

	_, existOID := modelType.FieldByName("OID")

	newData := reflect.New(modelType).Interface().(dao.IEntity)
	postItemData := reflect.New(modelType).Interface().(dao.IEntity)
	err = json.Unmarshal(m.Post.Body, postItemData)
	if err != nil {
		return nil, err
	}
	if postItemData != nil {
		v := reflect.ValueOf(newData).Elem()
		postItemDataV := reflect.ValueOf(postItemData).Elem()
		for key := range v.NumField() {
			fv := v.Field(key)
			if fv.CanSet() == false {
				continue
			}
			setv := postItemDataV.Field(key)
			if setv.CanConvert(fv.Type()) && setv.IsZero() == false {
				fv.Set(setv.Convert(fv.Type()))
			}
		}
	}
	/*if body, ok := m.Post.Body.(map[string]any); ok {
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
	}*/
	/*err = json.Unmarshal(body, newData)
	if err != nil {
		return nil, err
	}*/

	var primaryId dao.PrimaryKey

	if newData.Primary() > 0 {
		oldData := dao.GetByPrimaryKey(db.Orm(), table, newData.Primary())
		if oldData.IsZero() {
			reflect.ValueOf(newData).Elem().FieldByName("OID").Set(reflect.ValueOf(m.Admin.OID))
			err = dao.Create(db.Orm(), newData)
			primaryId = newData.Primary()
		} else {
			changeMap := com.Diff(oldData, newData)
			where["ID"] = newData.Primary()
			where["OID"] = m.Admin.OID
			err = dao.UpdateBy(db.Orm(), table, changeMap, where)
			primaryId = newData.Primary()
		}
	} else {
		if existOID && len(where) > 0 {
			where["OID"] = m.Admin.OID
			oldData := dao.GetBy(db.Orm(), table, where)
			if !oldData.IsZero() {
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

		if existOID {
			field := reflect.ValueOf(newData).Elem().FieldByName("OID")
			if field.CanSet() {
				field.Set(reflect.ValueOf(m.Admin.OID))
			}
		}
		err = dao.Create(db.Orm(), newData)
		primaryId = newData.Primary()

	}
	if err != nil {
		return nil, err
	}
	return result.NewDataMessage(dao.GetByPrimaryKey(db.Orm(), table, primaryId), "添加成功"), err
}
func (m *Restful) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	//m.bindQuery(ctx)
	table, err := entity.GetModel(m.Put.Model)
	if err != nil {
		return nil, err
	}

	modelType := reflect.TypeOf(table).Elem()

	newData := reflect.New(modelType).Interface().(dao.IEntity)
	err = json.Unmarshal(m.Put.Body, &newData)
	if err != nil {
		return nil, err
	}

	query := m.bindQuery(ctx, modelType)
	where := map[string]any{}
	for k := range query {
		where[k] = query.Get(k)
	}

	oldData := dao.GetByPrimaryKey(db.Orm(), table, m.Put.ID)
	if oldData.IsZero() {
		return nil, errors.New("数据不存在")
	}

	/*bytes, err := json.Marshal(oldData)
	if err != nil {
		return nil, err
	}
	oldMap := make(map[string]any)
	err = json.Unmarshal(bytes, &oldMap)
	if err != nil {
		return nil, err
	}*/

	changeMap := make(map[string]any) //com.Diff(&oldMap, &newMap)

	//newData := reflect.New(reflect.TypeOf(table).Elem()).Interface().(dao.IEntity)
	/*err = json.Unmarshal(body, newData)
	if err != nil {
		return nil, err
	}*/
	nDv := reflect.ValueOf(newData).Elem()
	oDv := reflect.ValueOf(oldData).Elem()
	for key := range oDv.NumField() {
		fv := nDv.Field(key)
		ofv := oDv.Field(key)

		if fv.IsZero() && fv.Kind() != reflect.Bool {
			continue
		}

		if reflect.DeepEqual(fv.Interface(), ofv.Interface()) == false {
			name := nDv.Type().Field(key).Name
			changeMap[name] = fv.Interface()
		}

		/*setv := reflect.ValueOf(changeMap[key])
		if setv.CanConvert(fv.Type()) {
			fv.Set(setv.Convert(fv.Type()))
		}
		changeMap[key] = nDv.FieldByName(key).Interface()*/
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
	return result.NewDataMessage(dao.GetByPrimaryKey(db.Orm(), table, m.Put.ID), "更新成功"), err
}
func (m *Restful) HandleDelete(ctx constrain.IContext) (constrain.IResult, error) {
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
		modelType := reflect.TypeOf(table).Elem()
		query := m.bindQuery(ctx, modelType)
		where := map[string]any{}
		for k := range query {
			where[k] = query.Get(k)
		}
		if len(where) > 0 {
			err = dao.DeleteBy(db.Orm(), table, where)
			if err != nil {
				return nil, err
			}
		}
	}
	return result.NewSuccess("删除成功"), nil
}
