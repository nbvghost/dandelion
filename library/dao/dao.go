package dao

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func UpdateBy(tx *gorm.DB, model IEntity, updateValue interface{}, query interface{}, args ...interface{}) error {
	return tx.Model(model).Where(query, args...).Updates(updateValue).Error
}

func UpdateByPrimaryKey(tx *gorm.DB, model IEntity, id PrimaryKey, updateValue any) error {
	var item = reflect.New(reflect.TypeOf(model).Elem()).Interface().(IEntity)
	return tx.Model(model).Where(fmt.Sprintf(`"%s"=?`, item.PrimaryName()), id).Updates(updateValue).Error
}

func GetByPrimaryKey(tx *gorm.DB, model IEntity, id PrimaryKey) IEntity {
	var item = reflect.New(reflect.TypeOf(model).Elem())
	tx.Model(model).Where(fmt.Sprintf(`"%s"=?`, item.Interface().(IEntity).PrimaryName()), id).Take(item.Interface())
	return item.Interface().(IEntity)
}
func GetBy(tx *gorm.DB, model IEntity, where map[string]any) IEntity {
	var item = reflect.New(reflect.TypeOf(model).Elem())
	tx.Model(model).Where(where).Take(item.Interface())
	return item.Interface().(IEntity)
}
func Create(tx *gorm.DB, value IEntity) error {
	return tx.Model(value).Create(value).Error
}
func Save(tx *gorm.DB, value IEntity) error {
	return tx.Save(value).Error
}

func DeleteByPrimaryKey(tx *gorm.DB, model IEntity, id PrimaryKey) error {
	return tx.Delete(reflect.New(reflect.TypeOf(model).Elem()).Interface(), id).Error
}
func DeleteBy(tx *gorm.DB, model IEntity, where map[string]any) error {
	return tx.Where(where).Delete(reflect.New(reflect.TypeOf(model).Elem()).Interface()).Error
}

type FindQuery struct {
	model IEntity
	isSetOrder bool
	db    *gorm.DB
}

func (m *FindQuery) PrimaryKey(ID PrimaryKey) *FindQuery {
	m.db.Where(ID)
	return m
}
func (m *FindQuery) Select(query interface{}, args ...interface{}) *FindQuery {
	m.db.Select(query, args...)
	return m
}
func (m *FindQuery) Scan(dest interface{}) error {
	return m.db.Scan(dest).Error
}
func (m *FindQuery) Where(query interface{}, args ...interface{}) *FindQuery {
	m.db.Where(query, args...)
	return m
}
func (m *FindQuery) Order(order ...string) *FindQuery {
	m.isSetOrder=true
	m.db.Order(strings.Join(order, ","))
	return m
}
func (m *FindQuery) OrderRaw(value interface{}) *FindQuery {
	m.isSetOrder=true
	m.db.Order(value)
	return m
}
func (m *FindQuery) Count() int64 {
	var total int64
	m.db.Count(&total)
	return total
}
func (m *FindQuery) Limit(index, pageSize int) int64 {
	var total int64
	if pageSize < 0 {
		m.db.Count(&total)
	} else {
		if index < 0 {
			index = 0
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		m.db.Count(&total).Limit(pageSize).Offset(pageSize * index)
	}
	return total
}
func (m *FindQuery) LimitOnly(pageSize int) *FindQuery {
	m.db.Limit(pageSize)
	return m
}
func (m *FindQuery) Group(column string) (any, error) {
	s, ok := reflect.TypeOf(m.model).Elem().FieldByName(column)
	if !ok {
		return nil, errors.Errorf("没有找到字段%s", column)
	}
	var list = reflect.New(reflect.SliceOf(s.Type))
	m.db.Select(column).Group(column).Find(list.Interface())
	return list.Elem().Interface(), nil
}
func (m *FindQuery) Pluck(column string, dest interface{}) {
	m.db.Pluck(column, dest)
}
func (m *FindQuery) Result(dest interface{}) {
	if !m.isSetOrder {
		m.db.Order(fmt.Sprintf(`"%s" asc`, m.model.PrimaryName()))
	}
	v := reflect.TypeOf(dest).Elem()
	log.Println(v.Kind())
	if v.Kind() == reflect.Slice {
		m.db.Find(dest)
	} else {
		m.db.First(dest)
	}
}
func (m *FindQuery) List() []IEntity {
	var list = reflect.New(reflect.SliceOf(reflect.TypeOf(m.model)))
	if !m.isSetOrder {
		m.db.Order(fmt.Sprintf(`"%s" asc`, m.model.PrimaryName()))
	}
	m.db.Find(list.Interface())
	arr := list.Elem()
	l := arr.Len()
	resultList := make([]IEntity, l)
	for i := 0; i < l; i++ {
		resultList[i] = arr.Index(i).Interface().(IEntity)
	}
	return resultList

}
func Find(tx *gorm.DB, model IEntity) *FindQuery {
	t := reflect.TypeOf(model).Elem()
	return &FindQuery{
		model: reflect.New(t).Interface().(IEntity),
		db:    tx.Model(model),
	}
}
