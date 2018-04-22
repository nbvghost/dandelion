package service

import "dandelion/app/service/dao"

type CategoryService struct {
	dao dao.CategoryDao
}

func (self CategoryService) DelCategory(ID uint64) {

	self.dao.DelCategory(Orm, ID)
}
func (self CategoryService) AddCategory(Label string) (*dao.Category, bool) {
	return self.dao.AddCategory(Orm, Label)
}
func (self CategoryService) FindCategory() []dao.Category {
	return self.dao.FindCategory(Orm)
}
