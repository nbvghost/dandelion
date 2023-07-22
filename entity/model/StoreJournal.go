package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// 账目明细
type StoreJournal struct {
	dao.Entity
	Name     string         `gorm:"column:Name;not null"`
	Detail   string         `gorm:"column:Detail;not null"`
	StoreID  dao.PrimaryKey `gorm:"column:StoreID"`
	Type     int            `gorm:"column:Type"`    //1=自主核销，2=在线充值
	Amount   int64          `gorm:"column:Amount"`  //变动金额
	Balance  uint           `gorm:"column:Balance"` //变动后的余额
	TargetID dao.PrimaryKey `gorm:"column:TargetID"`
}

func (StoreJournal) TableName() string {
	return "StoreJournal"
}
