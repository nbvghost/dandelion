package model

import "github.com/nbvghost/dandelion/library/dao"

type PushData struct {
	dao.Entity
	Content string `gorm:"column:Content;type:json"`
}

func (PushData) TableName() string {
	return "PushData"
}
