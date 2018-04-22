package service

import (
	"dandelion/app/service/dao"
)

//http://www.stats.gov.cn/tjsj/tjbz/xzqhdm/
type DistrictService struct {
	ProvinceDao dao.ProvinceDao
	CityDao     dao.CityDao
	AreaDao     dao.AreaDao
}

func (b DistrictService) Add(target interface{}) error {

	return Orm.Create(target).Error
}
