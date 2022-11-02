package dao

import (
	"reflect"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
)

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
func Find(tx *gorm.DB, model types.IEntity) []types.IEntity {
	return FindBy(tx, model, map[string]any{})
}

func DeleteByPrimaryKey(tx *gorm.DB, model types.IEntity, id types.PrimaryKey) error {
	return tx.Delete(reflect.New(reflect.TypeOf(model).Elem()).Interface(), id).Error
}
func DeleteBy(tx *gorm.DB, model types.IEntity, where map[string]any) error {
	return tx.Where(where).Delete(reflect.New(reflect.TypeOf(model).Elem()).Interface()).Error
}
func FindBy(tx *gorm.DB, model types.IEntity, where map[string]any) []types.IEntity {
	var list = reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
	tx.Model(model).Where(where).Find(list.Interface())
	arr := list.Elem()
	l := arr.Len()
	resultList := make([]types.IEntity, l)
	for i := 0; i < l; i++ {
		resultList[i] = arr.Index(i).Interface().(types.IEntity)
	}
	return resultList
}
