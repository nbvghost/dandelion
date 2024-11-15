package dao

import "gorm.io/gorm"

func NewPagerScope(page int, pageSize int) func(*gorm.DB) *gorm.DB {
	if page < 1 {
		page = 1
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize<10{
		pageSize=10
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
}
