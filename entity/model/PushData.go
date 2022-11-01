package model

import "github.com/nbvghost/gpa/types"

type PushData struct {
	types.Entity
	Content string `gorm:"column:Content;type:json"`
}

func (PushData) TableName() string {
	return "PushData"
}
