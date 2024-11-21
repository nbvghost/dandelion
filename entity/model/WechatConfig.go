package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type WechatConfig struct {
	dao.Entity
	OID                        dao.PrimaryKey `gorm:"column:OID;unique"`
	Name                       string         `gorm:"column:Name"`
	AppID                      string         `gorm:"column:AppID;unique"`
	AppSecret                  string         `gorm:"column:AppSecret"`
	MchID                      string         `gorm:"column:MchID"`
	MchAPIv2Key                string         `gorm:"column:MchAPIv2Key"`
	MchAPIv3Key                string         `gorm:"column:MchAPIv3Key"`
	MchCertificateSerialNumber string         `gorm:"column:MchCertificateSerialNumber"`
	PrivateKey                 string         `gorm:"column:PrivateKey"`
	OrderNotifyUrl             string         `gorm:"column:OrderNotifyUrl"`
	RefundNotifyUrl            string         `gorm:"column:RefundNotifyUrl"`
}

func (WechatConfig) TableName() string {
	return "WechatConfig"
}
