package model

import (
	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//账目明细
type UserJournal struct {
	base.BaseModel
	UserID     types.PrimaryKey `gorm:"column:UserID"`           //受益者
	Name       string           `gorm:"column:Name;not null"`    //
	Detail     string           `gorm:"column:Detail;not null"`  //
	Type       int              `gorm:"column:Type"`             //ddddd
	Amount     int64            `gorm:"column:Amount"`           //
	Balance    uint             `gorm:"column:Balance"`          //
	FromUserID types.PrimaryKey `gorm:"column:FromUserID"`       //来源
	DataKV     string           `gorm:"column:DataKV;type:text"` //{Key:"",Value:""}
}

func (UserJournal) TableName() string {
	return "UserJournal"
}
