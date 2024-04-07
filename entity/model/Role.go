package model

import (
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
)

type RoleStatus string

const (
	RoleStatusEnable  RoleStatus = "enable"
	RoleStatusDisable RoleStatus = "disable"
)

// Role 角色
type Role struct {
	dao.Entity
	OID    dao.PrimaryKey                `gorm:"column:OID"`
	Name   string                        `gorm:"column:Name"`
	Status RoleStatus                    `gorm:"column:Status"`
	Remark string                        `gorm:"column:Remark"`
	Routes sqltype.Array[*sqltype.Route] `gorm:"column:Routes;type:JSON"`
}

func (Role) TableName() string {
	return "Role"
}
