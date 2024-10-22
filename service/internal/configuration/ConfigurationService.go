package configuration

import (
	"encoding/json"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"github.com/nbvghost/tool/object"
	"gorm.io/gorm"
	"log"
)

type ConfigurationService struct {
	//model.BaseDao
}

func (m ConfigurationService) GetConfiguration(tx *gorm.DB, OID dao.PrimaryKey, Key model.ConfigurationKey) model.Configuration {
	var item model.Configuration
	err := tx.Where(`"K"=? and "OID"=?`, Key, OID).First(&item).Error
	if err != nil {
		log.Println(err)
	}
	return item
}
func (m ConfigurationService) GetConfigurations(OID dao.PrimaryKey, keys ...model.ConfigurationKey) map[model.ConfigurationKey]string {
	Orm := db.Orm()
	var items []model.Configuration
	err := Orm.Where(`"K" in (?) and "OID"=?`, keys, OID).Find(&items).Error
	//db.Where([]int64{20, 21, 22}).Find(&users
	if err != nil {
		log.Println(err)
	}

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
func (m ConfigurationService) ChangeConfiguration(db *gorm.DB, OID dao.PrimaryKey, Key model.ConfigurationKey, Value string) error {

	item := m.GetConfiguration(db, OID, Key)
	item.V = Value
	if item.ID == 0 {
		item.K = Key
		item.V = Value
		item.OID = OID
		return db.Create(&item).Error
	} else {
		return db.Model(&model.Configuration{}).Where(`"K"=? and "OID"=?`, Key, OID).Updates(map[string]interface{}{"V": Value}).Error
	}
}

func (m ConfigurationService) GetAdvertConfiguration(oid dao.PrimaryKey) []Advert {
	c := m.GetConfiguration(db.Orm(), oid, model.ConfigurationKeyAdvert)
	var h []Advert
	_ = json.Unmarshal([]byte(c.V), &h)
	return h
}

type EmailSTMP struct {
	EmailSTMPFrom     string
	EmailSTMPHost     string
	EmailSTMPPort     string
	EmailSTMPPassword string
}

func (m ConfigurationService) GetEmailSTMP(oid dao.PrimaryKey) *EmailSTMP {
	cs := m.GetConfigurations(oid, model.ConfigurationKeyEmailSTMPFrom, model.ConfigurationKeyEmailSTMPHost, model.ConfigurationKeyEmailSTMPPort, model.ConfigurationKeyEmailSTMPPassword)
	return &EmailSTMP{
		EmailSTMPFrom:     cs[model.ConfigurationKeyEmailSTMPFrom],
		EmailSTMPHost:     cs[model.ConfigurationKeyEmailSTMPHost],
		EmailSTMPPort:     cs[model.ConfigurationKeyEmailSTMPPort],
		EmailSTMPPassword: cs[model.ConfigurationKeyEmailSTMPPassword],
	}
}

func (m ConfigurationService) GetPopConfiguration(oid dao.PrimaryKey) []Pop {
	c := m.GetConfiguration(db.Orm(), oid, model.ConfigurationKeyPop)
	var h []Pop
	_ = json.Unmarshal([]byte(c.V), &h)
	return h
}
func (m ConfigurationService) GetQuickLinkConfiguration(oid dao.PrimaryKey) []QuickLink {
	c := m.GetConfiguration(db.Orm(), oid, model.ConfigurationKeyQuickLink)
	var h []QuickLink
	_ = json.Unmarshal([]byte(c.V), &h)
	return h
}
func (m ConfigurationService) GetBaiduTranslateConfiguration(oid dao.PrimaryKey) *serviceargument.BaiduTranslateConfig {
	c := m.GetConfigurations(oid, model.ConfigurationKeyBaiduTranslateAppID, model.ConfigurationKeyBaiduTranslateAppKey)
	h := serviceargument.BaiduTranslateConfig{
		AppID:  c[model.ConfigurationKeyBaiduTranslateAppID],
		AppKey: c[model.ConfigurationKeyBaiduTranslateAppKey],
	}
	return &h
}
func (m ConfigurationService) GetTranslate(db *gorm.DB, oid dao.PrimaryKey) []model.Configuration {
	var list []model.Configuration
	db.Model(&model.Configuration{}).Where(`"OID"=?`, oid).Where(`"K" like 'Translate%'`).Order(`"V"::int desc`).Find(&list)
	return list
}
func (m ConfigurationService) GetVolcengineConfiguration(oid dao.PrimaryKey) *serviceargument.Volcengine {
	c := m.GetConfigurations(oid, model.ConfigurationKeyVolcengineAccessKeyID, model.ConfigurationKeyVolcengineAccessKeySecret)
	h := serviceargument.Volcengine{
		AccessKeyID:     c[model.ConfigurationKeyVolcengineAccessKeyID],
		AccessKeySecret: c[model.ConfigurationKeyVolcengineAccessKeySecret],
	}
	return &h
}
func (m ConfigurationService) GetAliyunConfiguration(oid dao.PrimaryKey) *serviceargument.AliyunConfig {
	c := m.GetConfigurations(oid, model.ConfigurationKeyAliyunAccessKeyID, model.ConfigurationKeyAliyunAccessKeySecret)
	h := serviceargument.AliyunConfig{
		AccessKeyID:     c[model.ConfigurationKeyAliyunAccessKeyID],
		AccessKeySecret: c[model.ConfigurationKeyAliyunAccessKeySecret],
	}
	return &h
}
func (m ConfigurationService) GetHeaderConfiguration(oid dao.PrimaryKey) *Header {
	c := m.GetConfiguration(db.Orm(), oid, model.ConfigurationKeyHeader)
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

	leve1 := object.ParseFloat(c[model.ConfigurationKeyBrokerageLeve1])
	leve2 := object.ParseFloat(c[model.ConfigurationKeyBrokerageLeve2])
	leve3 := object.ParseFloat(c[model.ConfigurationKeyBrokerageLeve3])
	leve4 := object.ParseFloat(c[model.ConfigurationKeyBrokerageLeve4])
	leve5 := object.ParseFloat(c[model.ConfigurationKeyBrokerageLeve5])
	leve6 := object.ParseFloat(c[model.ConfigurationKeyBrokerageLeve6])

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
