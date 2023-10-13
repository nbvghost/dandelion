package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type UserJournalType int

const (
	UserJournal_Type_HX        UserJournalType = 1 //核销
	UserJournal_Type_LEVE      UserJournalType = 2 //下单，上下级结算佣金
	UserJournal_Type_TX        UserJournalType = 3 //提现
	UserJournal_Type_USER_LEVE UserJournalType = 4 //成为上下级，结算佣金
)

// UserJournal 账目明细
type UserJournal struct {
	dao.Entity
	UserID       dao.PrimaryKey  `gorm:"column:UserID"`           //受益者
	Name         string          `gorm:"column:Name;not null"`    //
	Detail       string          `gorm:"column:Detail;not null"`  //
	Type         UserJournalType `gorm:"column:Type"`             //ddddd
	Amount       int64           `gorm:"column:Amount"`           //
	Balance      uint            `gorm:"column:Balance"`          //余额是本次的账单金额
	FromUserID   dao.PrimaryKey  `gorm:"column:FromUserID"`       //来源
	FromUserName string          `gorm:"column:FromUserName"`     //来源
	DataKV       string          `gorm:"column:DataKV;type:JSON"` //{Key:"",Value:""}
}

func (UserJournal) TableName() string {
	return "UserJournal"
}

type FreezeType int

const (
	FreezeTypeFreeze   FreezeType = 0 //冻结
	FreezeTypeUnFreeze FreezeType = 1 //解冻并结算到用户的资金账户里
	FreezeTypeDisable  FreezeType = 2 //无效冻结，表示，资金来源已经无效，此冻结也无效
)

// UserFreezeJournal 被交结的账目明细
type UserFreezeJournal struct {
	dao.Entity
	UserID       dao.PrimaryKey  `gorm:"column:UserID"`           //受益者
	Name         string          `gorm:"column:Name;not null"`    //
	Detail       string          `gorm:"column:Detail;not null"`  //
	Type         UserJournalType `gorm:"column:Type"`             //ddddd
	Amount       int64           `gorm:"column:Amount"`           //
	FromUserID   dao.PrimaryKey  `gorm:"column:FromUserID"`       //来源
	FromUserName string          `gorm:"column:FromUserName"`     //来源
	DataKV       string          `gorm:"column:DataKV;type:JSON"` //{Key:"",Value:""}
	FreezeType   FreezeType      `gorm:"column:FreezeType;index"` //
}

func (UserFreezeJournal) TableName() string {
	return "UserFreezeJournal"
}

type ScoreJournalType int

const (
	ScoreJournal_Type_GM           ScoreJournalType = 1 //购买商品
	ScoreJournal_Type_DH           ScoreJournalType = 2 //积分兑换商品
	ScoreJournal_Type_LEVE         ScoreJournalType = 3 //上下级结算佣金,获取的积分
	ScoreJournal_Type_DaySign      ScoreJournalType = 4 //签到送积分
	ScoreJournal_Type_Look_Article ScoreJournalType = 5 //看文章送积分
	ScoreJournal_Type_Share        ScoreJournalType = 6 //转发获历
	ScoreJournal_Type_InviteUser   ScoreJournalType = 7 //邀请好友
	ScoreJournal_Type_Look_Video   ScoreJournalType = 8 //看视频送积分
)

// ScoreJournal Score明细
type ScoreJournal struct {
	dao.Entity
	Name    string           `gorm:"column:Name;not null"`   //
	Detail  string           `gorm:"column:Detail;not null"` //
	UserID  dao.PrimaryKey   `gorm:"column:UserID"`          //
	Score   int64            `gorm:"column:Score"`           //变动金额
	Type    ScoreJournalType `gorm:"column:Type"`            //
	Balance uint             `gorm:"column:Balance"`         //变动后的余额
	//DataKV  string         `gorm:"column:DataKV;type:text"` //{Key:"",Value:""}积分不记录，获取途径
}

func (ScoreJournal) TableName() string {
	return "ScoreJournal"
}
