package service

import "github.com/nbvghost/dandelion/app/service/dao"

//http://www.stats.gov.cn/tjsj/tjbz/xzqhdm/
//https://github.com/mumuy/data_location

type DistrictService struct {
	dao.BaseDao
}

func (b DistrictService) Add(target interface{}) error {

	return b.Add(target)
}
