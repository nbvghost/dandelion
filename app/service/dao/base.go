package dao

import (
	"github.com/jinzhu/gorm"
)

type BaseDao struct {
}

func (b BaseDao) Delete(DB *gorm.DB, target interface{}, ID uint64) error {

	return DB.Delete(target, "ID=?", ID).Error
}
func (b BaseDao) Add(DB *gorm.DB, target interface{}) error {

	return DB.Create(target).Error
}
func (b BaseDao) ChangeModel(DB *gorm.DB, ID uint64, target interface{}) error {

	return DB.Model(target).Where("ID=?", ID).Updates(target).Error
}
func (b BaseDao) ChangeMap(DB *gorm.DB, ID uint64, target map[string]interface{}) error {

	return DB.Model(target).Where("ID=?", ID).Updates(target).Error
}
func (b BaseDao) Get(DB *gorm.DB, ID uint64, target interface{}) error {
	return DB.Where("ID=?", ID).First(target).Error
}
func (b BaseDao) FindAll(DB *gorm.DB, target interface{}) error {

	return DB.Find(target).Error
}
