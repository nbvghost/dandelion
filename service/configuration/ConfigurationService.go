package configuration

import (
	"log"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/singleton"

	"github.com/nbvghost/gpa/types"
)

type ConfigurationService struct {
	model.BaseDao
}

func (b ConfigurationService) GetConfiguration(OID types.PrimaryKey, Key sqltype.ConfigurationKey) model.Configuration {
	Orm := singleton.Orm()
	var item model.Configuration
	err := Orm.Where(`"K"=? and "OID"=?`, Key, OID).First(&item).Error
	//db.Where([]int64{20, 21, 22}).Find(&users
	log.Println(err)
	return item
}
func (b ConfigurationService) GetConfigurations(OID types.PrimaryKey, Keys []sqltype.ConfigurationKey) map[sqltype.ConfigurationKey]string {
	Orm := singleton.Orm()
	var items []model.Configuration
	err := Orm.Where(`"K" in (?) and "OID"=?`, Keys, OID).Find(&items).Error
	//db.Where([]int64{20, 21, 22}).Find(&users
	log.Println(err)

	list := make(map[sqltype.ConfigurationKey]string)
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
func (b ConfigurationService) ChangeConfiguration(OID types.PrimaryKey, Key sqltype.ConfigurationKey, Value string) error {
	Orm := singleton.Orm()
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
