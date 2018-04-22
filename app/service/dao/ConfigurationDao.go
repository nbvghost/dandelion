package dao

import (
	"dandelion/app/util"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb/tool"
)

type ConfigurationDao struct {
}

func (b ConfigurationDao) GetConfiguration(DB *gorm.DB, Key uint64) *Configuration {
	item := &Configuration{}
	err := DB.Where("K=?", Key).First(item).Error
	tool.CheckError(err)

	if item.ID == 0 {
		item.K = Key
		err = DB.Create(item).Error
		tool.CheckError(err)
		return item
	} else {
		return item
	}
}
func (b ConfigurationDao) ChangeConfiguration(DB *gorm.DB, Key uint64, Value string) error {

	item := b.GetConfiguration(DB, Key)
	item.V = Value

	return DB.Model(&Configuration{}).Where("K=?", Key).Updates(util.StructToMap(item)).Error
}
