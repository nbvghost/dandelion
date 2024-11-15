package model

import (
	"github.com/nbvghost/dandelion/library/dao"
	"time"
)

type UserInfoKey string

type UserInfoKeyType interface {
	string | bool | time.Time | int
}

const (
	UserInfoKeyDaySignTime     UserInfoKey = "DaySignTime"
	UserInfoKeyDaySignCount    UserInfoKey = "DaySignCount"
	UserInfoKeyLastIP          UserInfoKey = "LastIP"
	UserInfoKeyAllowAssistance UserInfoKey = "AllowAssistance"
	UserInfoKeySubscribe       UserInfoKey = "Subscribe"
	UserInfoKeyState           UserInfoKey = "State"
	UserInfoKeyAgent           UserInfoKey = "Agent"
)

type UserInfo struct {
	dao.Entity
	UserID dao.PrimaryKey `gorm:"column:UserID"`
	Key    UserInfoKey    `gorm:"column:Key"`
	Value  string         `gorm:"column:Value"`
	//DaySignTime     time.Time      `gorm:"column:DaySignTime"`     //最后一次签到
	//DaySignCount    int            `gorm:"column:DaySignCount"`    //签到次数
	//LastIP          string         `gorm:"column:LastIP"`          //登陆ip
	//AllowAssistance bool           `gorm:"column:AllowAssistance"` //AllowAssistance
	//Subscribe       bool           `gorm:"column:Subscribe"`
}

func (UserInfo) TableName() string {
	return "UserInfo"
}
