package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// 账目明细
type UserJournal struct {
	dao.Entity
	UserID     dao.PrimaryKey `gorm:"column:UserID"`           //受益者
	Name       string         `gorm:"column:Name;not null"`    //
	Detail     string         `gorm:"column:Detail;not null"`  //
	Type       int            `gorm:"column:Type"`             //ddddd
	Amount     int64          `gorm:"column:Amount"`           //
	Balance    uint           `gorm:"column:Balance"`          //
	FromUserID dao.PrimaryKey `gorm:"column:FromUserID"`       //来源
	DataKV     string         `gorm:"column:DataKV;type:text"` //{Key:"",Value:""}
}

func (UserJournal) TableName() string {
	return "UserJournal"
}
