package dao

import "time"

type QuestionTag struct {
	BaseModel
	Name      string
	Introduce string
}

func (QuestionTag) TableName() string {
	return "QuestionTag"
}

type Question struct {
	BaseModel
	UserID        uint64    `gorm:"column:UserID"`               //
	Title         string    `gorm:"column:Title"`                //
	Describe      string    `gorm:"column:Describe;type:text"`   //描述
	Link          string    `gorm:"column:Link"`                 //相关链接
	Long          float64   `gorm:"column:Long"`                 //
	Lat           float64   `gorm:"column:Lat"`                  //
	AtLocal       bool      `gorm:"column:AtLocal"`              //是否针对本地
	Attachment    string    `gorm:"column:Attachment;type:text"` //附件json,图片，文件等
	QuestionTagID uint64    `gorm:"column:QuestionTagID"`        //多对1，话题/类别
	View          uint64    `gorm:"column:View"`                 //
	Follow        uint64    `gorm:"column:Follow"`               //
	Expiry        time.Time `gorm:"column:Expiry"`               //问题超时时间
	Status        string    `gorm:"column:Status"`               //问题状态
	AnswerReward  uint64    `gorm:"column:AnswerReward"`         //回答奖励金额，分
	ShareReward   uint64    `gorm:"column:ShareReward"`          //分享奖励金额，分
	Anonymous     bool      `gorm:"column:Anonymous"`            //匿名提问
}

func (Question) TableName() string {
	return "Question"
}

type AnswerQuestion struct {
	BaseModel
	UserID     uint64
	QuestionID uint64
	Content    string
	Praise     uint64
}

func (AnswerQuestion) TableName() string {
	return "AnswerQuestion"
}
