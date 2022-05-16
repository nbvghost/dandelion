package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
)

type Logger struct {
	base.BaseModel
	OID   uint   `gorm:"column:OID"`
	Key   int    `gorm:"column:Key"`
	Title string `gorm:"column:Title"`
	Data  string `gorm:"column:Data"`
}

func (u *Logger) BeforeCreate(scope *gorm.DB) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))

	}
	return nil
}
func (Logger) TableName() string {
	return "Logger"
}
