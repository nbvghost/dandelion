package configuration

import (
	"encoding/json"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/tool/object"
	"log"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type ConfigurationService struct {
	//model.BaseDao
}

func (m ConfigurationService) GetConfiguration(OID dao.PrimaryKey, Key model.ConfigurationKey) model.Configuration {
	Orm := db.Orm()
	var item model.Configuration
	err := Orm.Where(`"K"=? and "OID"=?`, Key, OID).First(&item).Error
	//db.Where([]int64{20, 21, 22}).Find(&users
	log.Println(err)
	return item
}
func (m ConfigurationService) GetConfigurations(OID dao.PrimaryKey, keys ...model.ConfigurationKey) map[model.ConfigurationKey]string {
	Orm := db.Orm()
	var items []model.Configuration
	err := Orm.Where(`"K" in (?) and "OID"=?`, keys, OID).Find(&items).Error
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
func (m ConfigurationService) ChangeConfiguration(OID dao.PrimaryKey, Key model.ConfigurationKey, Value string) error {
	Orm := db.Orm()
	item := m.GetConfiguration(OID, Key)
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

func (m ConfigurationService) GetAdvertConfiguration(oid dao.PrimaryKey) []Advert {
	c := m.GetConfiguration(oid, model.ConfigurationKeyAdvert)
	var h []Advert
	_ = json.Unmarshal([]byte(c.V), &h)
	return h
}

func (m ConfigurationService) GetPopConfiguration(oid dao.PrimaryKey) []Pop {
	c := m.GetConfiguration(oid, model.ConfigurationKeyPop)
	var h []Pop
	_ = json.Unmarshal([]byte(c.V), &h)
	return h
}
func (m ConfigurationService) GetQuickLinkConfiguration(oid dao.PrimaryKey) []QuickLink {
	c := m.GetConfiguration(oid, model.ConfigurationKeyQuickLink)
	var h []QuickLink
	_ = json.Unmarshal([]byte(c.V), &h)
	return h
}
func (m ConfigurationService) GetHeaderConfiguration(oid dao.PrimaryKey) *Header {
	c := m.GetConfiguration(oid, model.ConfigurationKeyHeader)
	h := Header{}
	_ = json.Unmarshal([]byte(c.V), &h)
	return &h
}

func (m ConfigurationService) GetBrokerageConfiguration(oid dao.PrimaryKey) *Brokerage {
	c := m.GetConfigurations(oid,
		model.ConfigurationKeyBrokerageType,
		model.ConfigurationKeyBrokerageLeve1,
		model.ConfigurationKeyBrokerageLeve2,
		model.ConfigurationKeyBrokerageLeve3,
		model.ConfigurationKeyBrokerageLeve4,
		model.ConfigurationKeyBrokerageLeve5,
		model.ConfigurationKeyBrokerageLeve6)

	brokerageType := c[model.ConfigurationKeyBrokerageType]

	leve1 := object.ParseUint(c[model.ConfigurationKeyBrokerageLeve1])
	leve2 := object.ParseUint(c[model.ConfigurationKeyBrokerageLeve2])
	leve3 := object.ParseUint(c[model.ConfigurationKeyBrokerageLeve3])
	leve4 := object.ParseUint(c[model.ConfigurationKeyBrokerageLeve4])
	leve5 := object.ParseUint(c[model.ConfigurationKeyBrokerageLeve5])
	leve6 := object.ParseUint(c[model.ConfigurationKeyBrokerageLeve6])

	return &Brokerage{
		Type:  BrokerageType(brokerageType),
		Leve1: leve1,
		Leve2: leve2,
		Leve3: leve3,
		Leve4: leve4,
		Leve5: leve5,
		Leve6: leve6,
	}
}
