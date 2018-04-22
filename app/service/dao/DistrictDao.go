package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb/tool"
)

type ProvinceDao struct {
}

func (b ProvinceDao) ListP(DB *gorm.DB) []Province {
	var target []Province
	err := DB.Find(&target).Error
	tool.CheckError(err)
	return target
}

type CityDao struct {
}

func (b CityDao) ListC(DB *gorm.DB, P int) []City {
	var target []City
	err := DB.Where("P=?", P).Find(&target).Error
	tool.CheckError(err)
	return target
}

type AreaDao struct {
}

func (b AreaDao) ListA(DB *gorm.DB, P int, C int) []Area {
	var target []Area
	err := DB.Where("P=?", P).Where("C=?", C).Find(&target).Error
	tool.CheckError(err)
	return target
}
