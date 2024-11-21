package express

import (
	"github.com/nbvghost/dandelion/entity/model"
)

//http://www.stats.gov.cn/tjsj/tjbz/xzqhdm/
//https://github.com/mumuy/data_location

type DistrictService struct {
	model.BaseDao
}

func (b DistrictService) Add(target interface{}) error {

	return b.Add(target)
}
