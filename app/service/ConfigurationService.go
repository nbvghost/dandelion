package service

import (
	"dandelion/app/service/dao"

	"github.com/nbvghost/gweb/tool"
)

type ConfigurationService struct {
	dao.BaseDao
}

func (b ConfigurationService) GetConfiguration(OID uint64, Key uint64) dao.Configuration {
	Orm := dao.Orm()
	var item dao.Configuration
	err := Orm.Where("K=? and OID=?", Key, OID).First(&item).Error
	//db.Where([]int64{20, 21, 22}).Find(&users
	tool.CheckError(err)
	return item
}
func (b ConfigurationService) GetConfigurations(OID uint64, Keys []uint64) map[uint64]string {
	Orm := dao.Orm()
	var items []dao.Configuration
	err := Orm.Where("K in (?) and OID=?", Keys, OID).Find(&items).Error
	//db.Where([]int64{20, 21, 22}).Find(&users
	tool.CheckError(err)

	list := make(map[uint64]string)
	for _, value := range items {
		list[value.K] = value.V
	}
	return list
	/*for key, value := range items {
		if value.ID == 0 {
			value.K = key
			err = Orm.Create(value).Error
			tool.CheckError(err)

		}
	}*/

}
func (b ConfigurationService) ChangeConfiguration(OID uint64, Key uint64, Value string) error {
	Orm := dao.Orm()
	item := b.GetConfiguration(OID, Key)
	item.V = Value
	if item.ID == 0 {
		item.K = Key
		item.V = Value
		item.OID = OID
		return Orm.Create(&item).Error
	} else {
		return Orm.Model(&dao.Configuration{}).Where("K=? and OID=?", Key, OID).Updates(map[string]interface{}{"V": Value}).Error
	}

}
