package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

// 核销记录-user，store
type Verification struct {
	dao.Entity
	VerificationNo string         `gorm:"column:VerificationNo;unique"` //订单号
	UserID         dao.PrimaryKey `gorm:"column:UserID"`
	Name           string         `gorm:"column:Name"`
	Label          string         `gorm:"column:Label"`
	CardItemID     dao.PrimaryKey `gorm:"column:CardItemID"`
	StoreID        dao.PrimaryKey `gorm:"column:StoreID"`
	StoreUserID    dao.PrimaryKey `gorm:"column:StoreUserID"` //门店核销管理员的用户ID
	Quantity       uint           `gorm:"column:Quantity"`    //核销数量
}

func (Verification) TableName() string {
	return "Verification"
}
