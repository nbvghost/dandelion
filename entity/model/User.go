package model

import (
	"time"

	"github.com/nbvghost/gpa/types"
)

type User struct {
	types.Entity
	Name        string           `gorm:"column:Name"`                  //
	OpenID      string           `gorm:"column:OpenID"`                //
	Email       string           `gorm:"column:Email"`                 //
	Tel         string           `gorm:"column:Tel"`                   //
	Password    string           `gorm:"column:Password" json:"-"`     //
	Age         int              `gorm:"column:Age"`                   //
	Region      string           `gorm:"column:Region"`                //
	Amount      uint             `gorm:"column:Amount"`                //现金
	BlockAmount uint             `gorm:"column:BlockAmount"`           //冻结现金
	Score       uint             `gorm:"column:Score"`                 //积分
	Growth      uint             `gorm:"column:Growth"`                //成长值
	Portrait    string           `gorm:"column:Portrait"`              //头像
	Gender      int              `gorm:"column:Gender"`                //性别 1男  2女
	Subscribe   int              `gorm:"column:Subscribe"`             //
	LastLoginAt time.Time        `gorm:"column:LastLoginAt;type:time"` //
	RankID      uint             `gorm:"column:RankID"`                //
	SuperiorID  types.PrimaryKey `gorm:"column:SuperiorID"`            //
}

func (u User) TableName() string {
	return "User"
}
