package model

import (
	"time"

	"github.com/nbvghost/dandelion/library/dao"
)

type QuestionTag struct {
	dao.Entity
	Name      string
	Introduce string
}

func (QuestionTag) TableName() string {
	return "QuestionTag"
}

type Question struct {
	dao.Entity
	UserID        uint      `gorm:"column:UserID"`               //
	Title         string    `gorm:"column:Title"`                //
	Describe      string    `gorm:"column:Describe;type:text"`   //描述
	Link          string    `gorm:"column:Link"`                 //相关链接
	Long          float64   `gorm:"column:Long"`                 //
	Lat           float64   `gorm:"column:Lat"`                  //
	AtLocal       bool      `gorm:"column:AtLocal"`              //是否针对本地
	Attachment    string    `gorm:"column:Attachment;type:text"` //附件json,图片，文件等
	QuestionTagID uint      `gorm:"column:QuestionTagID"`        //多对1，话题/类别
	View          uint      `gorm:"column:View"`                 //
	Follow        uint      `gorm:"column:Follow"`               //
	Expiry        time.Time `gorm:"column:Expiry"`               //问题超时时间
	Status        string    `gorm:"column:Status"`               //问题状态
	AnswerReward  uint      `gorm:"column:AnswerReward"`         //回答奖励金额，分
	ShareReward   uint      `gorm:"column:ShareReward"`          //分享奖励金额，分
	Anonymous     bool      `gorm:"column:Anonymous"`            //匿名提问
}

func (Question) TableName() string {
	return "Question"
}

type AnswerQuestion struct {
	dao.Entity
	UserID     uint
	QuestionID uint
	Content    string
	Praise     uint
}

func (AnswerQuestion) TableName() string {
	return "AnswerQuestion"
}
