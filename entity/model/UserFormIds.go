package model

import (
	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

type UserFormIds struct {
	base.BaseModel
	UserID types.PrimaryKey `gorm:"column:UserID"` //
	FormId string           `gorm:"column:FormId"` //formId 用于发送
}

func (UserFormIds) TableName() string {
	return "UserFormIds"
}
