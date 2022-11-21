package dao

import (
	"reflect"
	"strings"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
)

func UpdateBy(tx *gorm.DB, model types.IEntity, value interface{}, query interface{}, args ...interface{}) error {
	return tx.Model(model).Where(query, args...).Updates(value).Error
}

func UpdateByPrimaryKey(tx *gorm.DB, model types.IEntity, id types.PrimaryKey, value any) error {
	return tx.Model(model).Where(`"ID"=?`, id).Updates(value).Error
}

func GetByPrimaryKey(tx *gorm.DB, model types.IEntity, id types.PrimaryKey) types.IEntity {
	var item = reflect.New(reflect.TypeOf(model).Elem())
	tx.Model(model).Where(`"ID"=?`, id).Take(item.Interface())
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
func (m *FindQuery) Where(query interface{}, args ...interface{}) *FindQuery {
	m.db.Where(query, args...)
	return m
}
func (m *FindQuery) Order(order ...string) *FindQuery {
	m.db.Order(strings.Join(order, ","))
	return m
}
func (m *FindQuery) List() []types.IEntity {
	var list = reflect.New(reflect.SliceOf(reflect.TypeOf(m.model)))
	if len(m.order) == 0 {
		m.db.Order(`"ID" asc`)
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
	return &FindQuery{
		model: model,
		db:    tx.Model(model),
	}
}
