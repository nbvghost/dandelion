package dao

import "github.com/jinzhu/gorm"

type CategoryDao struct{}

func (CategoryDao) DelCategory(Orm *gorm.DB, ID uint64) {

	Orm.Delete(Category{}, "ID=?", ID)
}
func (CategoryDao) AddCategory(Orm *gorm.DB, Label string) (*Category, bool) {
	_category := Category{}
	Orm.Where("Label=?", Label).First(&_category)
	if _category.ID != 0 {

		return &_category, false
	} else {
		_category.Label = Label
		Orm.Create(&_category)
		return &_category, true
	}
}
func (CategoryDao) FindCategory(Orm *gorm.DB) []Category {
	list := []Category{}

	Orm.Find(&list)

	return list
}
