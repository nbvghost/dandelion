package dao

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
)

func UpdateBy(tx *gorm.DB, model types.IEntity, value interface{}, query interface{}, args ...interface{}) error {
	return tx.Model(model).Where(query, args...).Updates(value).Error
}

func UpdateByPrimaryKey(tx *gorm.DB, model types.IEntity, id types.PrimaryKey, value any) error {
	return tx.Model(model).Where(fmt.Sprintf(`"%s"=?`, model.PrimaryName()), id).Updates(value).Error
}

func GetByPrimaryKey(tx *gorm.DB, model types.IEntity, id types.PrimaryKey) types.IEntity {
	var item = reflect.New(reflect.TypeOf(model).Elem())
	tx.Model(model).Where(fmt.Sprintf(`"%s"=?`, model.PrimaryName()), id).Take(item.Interface())
	return item.Interface().(types.IEntity)
}
func GetBy(tx *gorm.DB, model types.IEntity, where map[string]any) types.IEntity {
	var item = reflect.New(reflect.TypeOf(model).Elem())
	tx.Model(model).Where(where).Take(item.Interface())
	return item.Interface().(types.IEntity)
}
func Create(tx *gorm.DB, value types.IEntity) error {
	return tx.Model(value).Create(value).Error
}
func Save(tx *gorm.DB, value types.IEntity) error {
	return tx.Save(value).Error
}

func DeleteByPrimaryKey(tx *gorm.DB, model types.IEntity, id types.PrimaryKey) error {
	return tx.Delete(reflect.New(reflect.TypeOf(model).Elem()).Interface(), id).Error
}
func DeleteBy(tx *gorm.DB, model types.IEntity, where map[string]any) error {
	return tx.Where(where).Delete(reflect.New(reflect.TypeOf(model).Elem()).Interface()).Error
}

type FindQuery struct {
	model types.IEntity
	order []string
	db    *gorm.DB
}

func (m *FindQuery) PrimaryKey(ID types.PrimaryKey) *FindQuery {
	m.db.Where(ID)
	return m
}
func (m *FindQuery) Select(query interface{}, args ...interface{}) *FindQuery {
	m.db.Select(query, args...)
	return m
}
func (m *FindQuery) Where(query interface{}, args ...interface{}) *FindQuery {
	m.db.Where(query, args...)
	return m
}
func (m *FindQuery) Order(order ...string) *FindQuery {
	m.db.Order(strings.Join(order, ","))
	return m
}
func (m *FindQuery) Count() int64 {
	var total int64
	m.db.Count(&total)
	return total
}
func (m *FindQuery) Limit(index, pageSize int) int64 {
	if index < 0 {
		index = 0
	}
	var total int64
	m.db.Count(&total).Limit(pageSize).Offset(pageSize * index)
	return total
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
func (m *FindQuery) List() []types.IEntity {
	var list = reflect.New(reflect.SliceOf(reflect.TypeOf(m.model)))
	if len(m.order) == 0 {
		m.db.Order(fmt.Sprintf(`"%s" asc`, m.model.PrimaryName()))
	}
	m.db.Find(list.Interface())
	arr := list.Elem()
	l := arr.Len()
	resultList := make([]types.IEntity, l)
	for i := 0; i < l; i++ {
		resultList[i] = arr.Index(i).Interface().(types.IEntity)
	}
	return resultList

}
func Find(tx *gorm.DB, model types.IEntity) *FindQuery {
	t := reflect.TypeOf(model).Elem()
	return &FindQuery{
		model: reflect.New(t).Interface().(types.IEntity),
		db:    tx.Model(model),
	}
}
