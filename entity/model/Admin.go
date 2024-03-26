package model

import (
	"github.com/lib/pq"
	"time"

	"github.com/nbvghost/dandelion/library/dao"
)

type Admin struct {
	dao.Entity
	OID         dao.PrimaryKey `gorm:"column:OID;uniqueIndex:admin_idx_unique_id"`
	Account     string         `gorm:"column:Account;not null;uniqueIndex:admin_idx_unique_id"`
	Phone       string         `gorm:"column:Phone;default:''"`
	Name        string         `gorm:"column:Name"`
	PassWord    string         `gorm:"column:PassWord;not null;default:''" json:"-"`
	RoleList    pq.Int64Array  `gorm:"column:RoleList;type:text[]"` //角色
	LastLoginAt time.Time      `gorm:"column:LastLoginAt"`
	Initiator   bool           `gorm:"column:Initiator"`
	//Authority   string         `gorm:"column:Authority;default:''"` //json 权限
}

func (Admin) TableName() string {
	return "Admin"
}
