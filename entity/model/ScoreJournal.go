package model

import (
	"github.com/nbvghost/gpa/types"
)

//Score明细
type ScoreJournal struct {
	types.Entity
	Name    string           `gorm:"column:Name;not null"`    //
	Detail  string           `gorm:"column:Detail;not null"`  //
	UserID  types.PrimaryKey `gorm:"column:UserID"`           //
	Score   int64            `gorm:"column:Score"`            //变动金额
	Type    int              `gorm:"column:Type"`             //
	Balance uint             `gorm:"column:Balance"`          //变动后的余额
	DataKV  string           `gorm:"column:DataKV;type:text"` //{Key:"",Value:""}
}

func (ScoreJournal) TableName() string {
	return "ScoreJournal"
}
