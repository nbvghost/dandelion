package model

import (
	"time"

	"github.com/nbvghost/dandelion/library/dao"
)

type User struct {
	dao.Entity
	OID         dao.PrimaryKey `gorm:"column:OID"`
	Name        string         `gorm:"column:Name"`                  //
	OpenID      string         `gorm:"column:OpenID"`                //
	Email       string         `gorm:"column:Email"`                 //
	Phone       string         `gorm:"column:Phone"`                 //
	Password    string         `gorm:"column:Password" json:"-"`     //
	Age         int            `gorm:"column:Age"`                   //
	Amount      uint           `gorm:"column:Amount"`                //现金
	BlockAmount uint           `gorm:"column:BlockAmount"`           //冻结现金
	Score       uint           `gorm:"column:Score"`                 //积分
	Growth      uint           `gorm:"column:Growth"`                //成长值
	Portrait    string         `gorm:"column:Portrait"`              //头像
	Gender      int            `gorm:"column:Gender"`                //性别 1男  2女
	LastLoginAt time.Time      `gorm:"column:LastLoginAt;type:time"` //
	RankID      uint           `gorm:"column:RankID"`                //
	SuperiorID  dao.PrimaryKey `gorm:"column:SuperiorID"`            //
	RoleID      dao.PrimaryKey `gorm:"column:RoleID"`                //
	//Subscribe   int              `gorm:"column:Subscribe"`             //move to UserInfo
}

func (u User) TableName() string {
	return "User"
}
