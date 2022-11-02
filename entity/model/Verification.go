package model

import (
	"github.com/nbvghost/gpa/types"
)

//核销记录-user，store
type Verification struct {
	types.Entity
	VerificationNo string           `gorm:"column:VerificationNo;unique"` //订单号
	UserID         types.PrimaryKey `gorm:"column:UserID"`
	Name           string           `gorm:"column:Name"`
	Label          string           `gorm:"column:Label"`
	CardItemID     types.PrimaryKey `gorm:"column:CardItemID"`
	StoreID        types.PrimaryKey `gorm:"column:StoreID"`
	StoreUserID    types.PrimaryKey `gorm:"column:StoreUserID"` //门店核销管理员的用户ID
	Quantity       uint             `gorm:"column:Quantity"`    //核销数量
}

func (Verification) TableName() string {
	return "Verification"
}
