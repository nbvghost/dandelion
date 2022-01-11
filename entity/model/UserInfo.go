package model

import (
	"time"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

type UserInfo struct {
	base.BaseModel
	UserID       types.PrimaryKey `gorm:"column:UserID"`
	DaySignTime  time.Time        `gorm:"column:DaySignTime"`  //最后一次签到
	DaySignCount int              `gorm:"column:DaySignCount"` //签到次数
}

func (UserInfo) TableName() string {
	return "UserInfo"
}
