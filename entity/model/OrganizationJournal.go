package model

import (
	"github.com/nbvghost/gpa/types"
)

//Organization
//商店账目明细
type OrganizationJournal struct {
	types.Entity
	OID     types.PrimaryKey `gorm:"column:OID"`              //OID
	Name    string           `gorm:"column:Name;not null"`    //
	Detail  string           `gorm:"column:Detail;not null"`  //
	Type    int              `gorm:"column:Type"`             //ddddd
	Amount  int64            `gorm:"column:Amount"`           //
	Balance uint             `gorm:"column:Balance"`          //
	DataKV  string           `gorm:"column:DataKV;type:text"` //{Key:"",Value:""}
}

func (OrganizationJournal) TableName() string {
	return "OrganizationJournal"
}
