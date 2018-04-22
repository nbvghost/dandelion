package dao

import "github.com/jinzhu/gorm"

type ClassifyDao = BaseDao

func (c ClassifyDao) FindByShopID(Orm *gorm.DB, CompanyID uint64) []Classify {
	classifys := []Classify{}
	Orm.Where("CompanyID=?", CompanyID).Find(&classifys)
	return classifys
}
