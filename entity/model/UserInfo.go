package model

import (
	"time"

	"github.com/nbvghost/gpa/types"
)

type UserInfo struct {
	types.Entity
	UserID       types.PrimaryKey `gorm:"column:UserID"`
	DaySignTime  time.Time        `gorm:"column:DaySignTime"`  //最后一次签到
	DaySignCount int              `gorm:"column:DaySignCount"` //签到次数
	LastIP       string           `gorm:"column:LastIP"`       //登陆ip
}

func (UserInfo) TableName() string {
	return "UserInfo"
}
