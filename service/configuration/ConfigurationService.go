package configuration

import (
	"github.com/nbvghost/dandelion/library/db"
	"log"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type ConfigurationService struct {
	model.BaseDao
}

func (b ConfigurationService) GetConfiguration(OID dao.PrimaryKey, Key model.ConfigurationKey) model.Configuration {
	Orm := db.Orm()
	var item model.Configuration
	err := Orm.Where(`"K"=? and "OID"=?`, Key, OID).First(&item).Error
	//db.Where([]int64{20, 21, 22}).Find(&users
	log.Println(err)
	return item
}
func (b ConfigurationService) GetConfigurations(OID dao.PrimaryKey, Keys ...model.ConfigurationKey) map[model.ConfigurationKey]string {
	Orm := db.Orm()
	var items []model.Configuration
	err := Orm.Where(`"K" in (?) and "OID"=?`, Keys, OID).Find(&items).Error
	//db.Where([]int64{20, 21, 22}).Find(&users
	log.Println(err)

	list := make(map[model.ConfigurationKey]string)
	for _, value := range items {
		list[value.K] = value.V
	}
	return list
	/*for key, value := range items {
		if value.ID == 0 {
			value.K = key
			err = Orm.Create(value).Error
			log.Println(err)

		}
	}*/

}
func (b ConfigurationService) ChangeConfiguration(OID dao.PrimaryKey, Key model.ConfigurationKey, Value string) error {
	Orm := db.Orm()
	item := b.GetConfiguration(OID, Key)
	item.V = Value
	if item.ID == 0 {
		item.K = Key
		item.V = Value
		item.OID = OID
		return Orm.Create(&item).Error
	} else {
		return Orm.Model(&model.Configuration{}).Where(`"K"=? and "OID"=?`, Key, OID).Updates(map[string]interface{}{"V": Value}).Error
	}

}
