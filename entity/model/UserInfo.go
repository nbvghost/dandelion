package model

import (
	"time"

	"github.com/nbvghost/dandelion/library/dao"
)

type UserInfo struct {
	dao.Entity
	UserID          dao.PrimaryKey `gorm:"column:UserID"`
	DaySignTime     time.Time      `gorm:"column:DaySignTime"`     //最后一次签到
	DaySignCount    int            `gorm:"column:DaySignCount"`    //签到次数
	LastIP          string         `gorm:"column:LastIP"`          //登陆ip
	AllowAssistance bool           `gorm:"column:AllowAssistance"` //AllowAssistance
	Subscribe       bool           `gorm:"column:Subscribe"`
}

func (UserInfo) TableName() string {
	return "UserInfo"
}
